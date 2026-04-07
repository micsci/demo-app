package handler

import (
	"context"

	pb "github.com/demo-marketplace/gen/go/marketplace/v1"
	"github.com/demo-marketplace/internal/model"
	"github.com/demo-marketplace/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Handler struct {
	pb.UnimplementedMarketplaceServiceServer
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// --- Health ---

func (h *Handler) GetHealth(ctx context.Context, _ *pb.GetHealthRequest) (*pb.GetHealthResponse, error) {
	return &pb.GetHealthResponse{Status: "ok"}, nil
}

// --- Categories ---

func (h *Handler) CreateCategory(ctx context.Context, req *pb.CreateCategoryRequest) (*pb.CreateCategoryResponse, error) {
	if req.GetCategory() == nil {
		return nil, status.Error(codes.InvalidArgument, "category is required")
	}

	cat, err := h.svc.CreateCategory(ctx, model.Category{
		Name: req.GetCategory().GetName(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateCategoryResponse{Category: categoryToProto(cat)}, nil
}

func (h *Handler) ListCategories(ctx context.Context, _ *pb.ListCategoriesRequest) (*pb.ListCategoriesResponse, error) {
	cats, err := h.svc.ListCategories(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Category, len(cats))
	for i, c := range cats {
		result[i] = categoryToProto(c)
	}
	return &pb.ListCategoriesResponse{Categories: result}, nil
}

func (h *Handler) DeleteCategory(ctx context.Context, req *pb.DeleteCategoryRequest) (*pb.DeleteCategoryResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid category id")
	}

	if err := h.svc.DeleteCategory(ctx, id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteCategoryResponse{}, nil
}

// --- Applications ---

func (h *Handler) CreateApplication(ctx context.Context, req *pb.CreateApplicationRequest) (*pb.CreateApplicationResponse, error) {
	if req.GetApplication() == nil {
		return nil, status.Error(codes.InvalidArgument, "application is required")
	}

	catID, err := uuid.Parse(req.GetApplication().GetCategoryId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid category_id")
	}

	app, err := h.svc.CreateApplication(ctx, model.Application{
		Name:        req.GetApplication().GetName(),
		DisplayName: req.GetApplication().GetDisplayName(),
		Description: req.GetApplication().GetDescription(),
		Active:      req.GetApplication().GetActive(),
		Visible:     req.GetApplication().GetVisible(),
		CategoryID:  catID,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateApplicationResponse{Application: applicationToProto(app)}, nil
}

func (h *Handler) GetApplication(ctx context.Context, req *pb.GetApplicationRequest) (*pb.GetApplicationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application id")
	}

	app, err := h.svc.GetApplication(ctx, id)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &pb.GetApplicationResponse{Application: applicationToProto(app)}, nil
}

func (h *Handler) ListApplications(ctx context.Context, req *pb.ListApplicationsRequest) (*pb.ListApplicationsResponse, error) {
	apps, err := h.svc.ListApplications(ctx, req.GetActiveOnly())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Application, len(apps))
	for i, a := range apps {
		result[i] = applicationToProto(a)
	}
	return &pb.ListApplicationsResponse{Applications: result}, nil
}

func (h *Handler) UpdateApplication(ctx context.Context, req *pb.UpdateApplicationRequest) (*pb.UpdateApplicationResponse, error) {
	if req.GetApplication() == nil {
		return nil, status.Error(codes.InvalidArgument, "application is required")
	}

	id, err := uuid.Parse(req.GetApplication().GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application id")
	}

	catID, err := uuid.Parse(req.GetApplication().GetCategoryId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid category_id")
	}

	app, err := h.svc.UpdateApplication(ctx, model.Application{
		ID:          id,
		Name:        req.GetApplication().GetName(),
		DisplayName: req.GetApplication().GetDisplayName(),
		Description: req.GetApplication().GetDescription(),
		Active:      req.GetApplication().GetActive(),
		Visible:     req.GetApplication().GetVisible(),
		CategoryID:  catID,
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdateApplicationResponse{Application: applicationToProto(app)}, nil
}

func (h *Handler) DeleteApplication(ctx context.Context, req *pb.DeleteApplicationRequest) (*pb.DeleteApplicationResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application id")
	}

	if err := h.svc.DeleteApplication(ctx, id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteApplicationResponse{}, nil
}

// --- Connections ---

func (h *Handler) CreateConnection(ctx context.Context, req *pb.CreateConnectionRequest) (*pb.CreateConnectionResponse, error) {
	if req.GetConnection() == nil {
		return nil, status.Error(codes.InvalidArgument, "connection is required")
	}

	appID, err := uuid.Parse(req.GetConnection().GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application_id")
	}

	conn, err := h.svc.CreateConnection(ctx, model.Connection{
		ApplicationID: appID,
		MerchantID:    req.GetConnection().GetMerchantId(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateConnectionResponse{Connection: connectionToProto(conn)}, nil
}

func (h *Handler) ListConnectionsByApplication(ctx context.Context, req *pb.ListConnectionsByApplicationRequest) (*pb.ListConnectionsByApplicationResponse, error) {
	appID, err := uuid.Parse(req.GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application_id")
	}

	conns, err := h.svc.ListConnectionsByApp(ctx, appID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Connection, len(conns))
	for i, c := range conns {
		result[i] = connectionToProto(c)
	}
	return &pb.ListConnectionsByApplicationResponse{Connections: result}, nil
}

func (h *Handler) UpdateConnectionStatus(ctx context.Context, req *pb.UpdateConnectionStatusRequest) (*pb.UpdateConnectionStatusResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid connection id")
	}

	conn, err := h.svc.UpdateConnectionStatus(ctx, id, req.GetStatus())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.UpdateConnectionStatusResponse{Connection: connectionToProto(conn)}, nil
}

func (h *Handler) DeleteConnection(ctx context.Context, req *pb.DeleteConnectionRequest) (*pb.DeleteConnectionResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid connection id")
	}

	if err := h.svc.DeleteConnection(ctx, id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteConnectionResponse{}, nil
}

// --- Reviews ---

func (h *Handler) CreateReview(ctx context.Context, req *pb.CreateReviewRequest) (*pb.CreateReviewResponse, error) {
	if req.GetReview() == nil {
		return nil, status.Error(codes.InvalidArgument, "review is required")
	}

	appID, err := uuid.Parse(req.GetReview().GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application_id")
	}

	review, err := h.svc.CreateReview(ctx, model.Review{
		ApplicationID: appID,
		MerchantID:    req.GetReview().GetMerchantId(),
		Rating:        int(req.GetReview().GetRating()),
		Body:          req.GetReview().GetBody(),
	})
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateReviewResponse{Review: reviewToProto(review)}, nil
}

func (h *Handler) ListReviewsByApplication(ctx context.Context, req *pb.ListReviewsByApplicationRequest) (*pb.ListReviewsByApplicationResponse, error) {
	appID, err := uuid.Parse(req.GetApplicationId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid application_id")
	}

	reviews, err := h.svc.ListReviewsByApp(ctx, appID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	result := make([]*pb.Review, len(reviews))
	for i, r := range reviews {
		result[i] = reviewToProto(r)
	}

	// Fetch average rating separately for each application.
	avg, err := h.svc.GetAverageRating(ctx, appID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListReviewsByApplicationResponse{
		Reviews:       result,
		AverageRating: avg,
	}, nil
}

func (h *Handler) DeleteReview(ctx context.Context, req *pb.DeleteReviewRequest) (*pb.DeleteReviewResponse, error) {
	id, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid review id")
	}

	if err := h.svc.DeleteReview(ctx, id); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &pb.DeleteReviewResponse{}, nil
}

// --- Proto converters ---

func categoryToProto(c model.Category) *pb.Category {
	return &pb.Category{
		Id:        c.ID.String(),
		Name:      c.Name,
		CreatedAt: timestamppb.New(c.CreatedAt),
		UpdatedAt: timestamppb.New(c.UpdatedAt),
	}
}

func applicationToProto(a model.Application) *pb.Application {
	return &pb.Application{
		Id:          a.ID.String(),
		Name:        a.Name,
		DisplayName: a.DisplayName,
		Description: a.Description,
		Active:      a.Active,
		Visible:     a.Visible,
		CategoryId:  a.CategoryID.String(),
		CreatedAt:   timestamppb.New(a.CreatedAt),
		UpdatedAt:   timestamppb.New(a.UpdatedAt),
	}
}

func connectionToProto(c model.Connection) *pb.Connection {
	return &pb.Connection{
		Id:            c.ID.String(),
		ApplicationId: c.ApplicationID.String(),
		MerchantId:    c.MerchantID,
		Status:        c.Status,
		CreatedAt:     timestamppb.New(c.CreatedAt),
		UpdatedAt:     timestamppb.New(c.UpdatedAt),
	}
}

func reviewToProto(r model.Review) *pb.Review {
	return &pb.Review{
		Id:            r.ID.String(),
		ApplicationId: r.ApplicationID.String(),
		MerchantId:    r.MerchantID,
		Rating:        int32(r.Rating),
		Body:          r.Body,
		CreatedAt:     timestamppb.New(r.CreatedAt),
		UpdatedAt:     timestamppb.New(r.UpdatedAt),
	}
}

// Ensure Handler implements the interface at compile time.
var _ pb.MarketplaceServiceServer = (*Handler)(nil)

