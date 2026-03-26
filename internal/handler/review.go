package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/demo-marketplace/internal/model"
	"github.com/demo-marketplace/internal/service"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	svc    *service.Service
	logger *slog.Logger
}

func NewReviewHandler(svc *service.Service, logger *slog.Logger) *ReviewHandler {
	return &ReviewHandler{svc: svc, logger: logger}
}

func (h *ReviewHandler) CreateReview(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ApplicationID string `json:"application_id"`
		AuthorEmail   string `json:"author_email"`
		Rating        int    `json:"rating"`
		Comment       string `json:"comment"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	h.logger.Info("creating review",
		"application_id", req.ApplicationID,
		"author", req.AuthorEmail,
		"rating", req.Rating,
	)

	appID, err := uuid.Parse(req.ApplicationID)
	if err != nil {
		http.Error(w, "invalid application_id", http.StatusBadRequest)
		return
	}

	review, err := h.svc.CreateReview(r.Context(), model.Review{
		ApplicationID: appID,
		AuthorEmail:   req.AuthorEmail,
		Rating:        req.Rating,
		Comment:       req.Comment,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(review)
}

func (h *ReviewHandler) ListReviews(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(r.PathValue("application_id"))
	if err != nil {
		http.Error(w, "invalid application_id", http.StatusBadRequest)
		return
	}

	reviews, err := h.svc.ListReviewsByApp(r.Context(), appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	avg, err := h.svc.GetAverageRating(r.Context(), appID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Reviews       []model.Review `json:"reviews"`
		AverageRating float64        `json:"average_rating"`
		Total         int            `json:"total"`
	}{
		Reviews:       reviews,
		AverageRating: avg,
		Total:         len(reviews),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ReviewHandler) SearchReviews(w http.ResponseWriter, r *http.Request) {
	appID, err := uuid.Parse(r.PathValue("application_id"))
	if err != nil {
		http.Error(w, "invalid application_id", http.StatusBadRequest)
		return
	}

	keyword := r.URL.Query().Get("q")

	reviews, err := h.svc.SearchReviews(r.Context(), appID, keyword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(reviews)
}

func (h *ReviewHandler) DeleteReview(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid review id", http.StatusBadRequest)
		return
	}

	if err := h.svc.DeleteReview(r.Context(), id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
