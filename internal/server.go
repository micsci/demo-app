package internal

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	pb "github.com/demo-marketplace/gen/go/marketplace/v1"
	"github.com/demo-marketplace/internal/config"
	"github.com/demo-marketplace/internal/handler"
	"github.com/demo-marketplace/internal/service"
	"github.com/demo-marketplace/internal/storage/postgres"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg    config.Config
	logger *slog.Logger
}

func NewServer(logger *slog.Logger, cfg config.Config) *Server {
	return &Server{cfg: cfg, logger: logger}
}

func (s *Server) Run(ctx context.Context) error {
	db, err := sqlx.Connect("postgres", s.cfg.DSN())
	if err != nil {
		return fmt.Errorf("connecting to database: %w", err)
	}
	defer db.Close()

	if s.cfg.MigrationsUp {
		if err := goose.Up(db.DB, "migrations"); err != nil {
			return fmt.Errorf("running migrations: %w", err)
		}
		s.logger.Info("migrations applied")
	}

	store := postgres.New(db)
	svc := service.New(store)
	h := handler.New(svc)

	// gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterMarketplaceServiceServer(grpcServer, h)
	reflection.Register(grpcServer)

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%s", s.cfg.GRPCPort))
	if err != nil {
		return fmt.Errorf("listening on gRPC port: %w", err)
	}

	// gRPC client connection (gateway → gRPC server)
	conn, err := grpc.NewClient(
		fmt.Sprintf("localhost:%s", s.cfg.GRPCPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("creating gRPC client: %w", err)
	}
	defer conn.Close()

	// gRPC-Gateway HTTP mux
	gwMux := runtime.NewServeMux()
	if err := pb.RegisterMarketplaceServiceHandler(ctx, gwMux, conn); err != nil {
		return fmt.Errorf("registering gateway handler: %w", err)
	}

	// Review REST endpoints
	reviewH := handler.NewReviewHandler(svc, s.logger)
	mux := http.NewServeMux()
	mux.Handle("/", gwMux)
	mux.HandleFunc("POST /v1/reviews", reviewH.CreateReview)
	mux.HandleFunc("GET /v1/applications/{application_id}/reviews", reviewH.ListReviews)
	mux.HandleFunc("GET /v1/applications/{application_id}/reviews/search", reviewH.SearchReviews)
	mux.HandleFunc("DELETE /v1/reviews/{id}", reviewH.DeleteReview)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.cfg.AppPort),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		s.logger.Info("starting gRPC server", "port", s.cfg.GRPCPort)
		return grpcServer.Serve(grpcLis)
	})

	g.Go(func() error {
		s.logger.Info("starting HTTP server (gRPC-Gateway)", "port", s.cfg.AppPort)
		return httpServer.ListenAndServe()
	})

	g.Go(func() error {
		<-gCtx.Done()
		s.logger.Info("shutting down")
		grpcServer.GracefulStop()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		return httpServer.Shutdown(shutdownCtx)
	})

	return g.Wait()
}
