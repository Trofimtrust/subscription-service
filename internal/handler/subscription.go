package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"subscription-service/internal/model"
	"subscription-service/internal/repository"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

// CreateSubscription godoc
// @Summary Создать подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body model.Subscription true "Subscription"
// @Success 201 {object} model.Subscription
// @Router /subscriptions [post]

func (h *Handler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	var request model.SubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", request.StartDate)
	if err != nil {
		http.Error(w, "start_date must be in MM-YYYY format", http.StatusBadRequest)
		return
	}

	var endDate *time.Time

	if request.EndDate != nil {
		t, err := time.Parse("01-2006", *request.EndDate)
		if err != nil {
			http.Error(w, "end_date must be in MM-YYYY format", http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	subscription := model.Subscription{
		ID:          uuid.New().String(),
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = h.repo.Create(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(toResponse(subscription))
}

// GetSubscriptions godoc
// @Summary Получить все подписки
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Router /subscriptions [get]

func (h *Handler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetSubscriptions called")

	subscriptions, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	responses := make([]model.SubscriptionResponse, 0, len(subscriptions))

	for _, s := range subscriptions {
		responses = append(responses, toResponse(s))
	}

	json.NewEncoder(w).Encode(responses)
}

// GetSubscriptionByID godoc
// @Summary Получить подписку по ID
// @Tags subscriptions
// @Produce json
// @Param id path string true "Subscription ID"
// @Success 200 {object} model.Subscription
// @Router /subscriptions/{id} [get]

func (h *Handler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	subscription, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Subscription not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toResponse(*subscription))
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Tags subscriptions
// @Param id path string true "Subscription ID"
// @Success 204
// @Router /subscriptions/{id} [delete]

func (h *Handler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.repo.Delete(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "Subscription ID"
// @Param subscription body model.Subscription true "Subscription"
// @Success 200 {object} model.Subscription
// @Router /subscriptions/{id} [put]

func (h *Handler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var request model.SubscriptionRequest

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	startDate, err := time.Parse("01-2006", request.StartDate)
	if err != nil {
		http.Error(w, "start_date must be in MM-YYYY format", http.StatusBadRequest)
		return
	}

	var endDate *time.Time

	if request.EndDate != nil {
		t, err := time.Parse("01-2006", *request.EndDate)
		if err != nil {
			http.Error(w, "end_date must be in MM-YYYY format", http.StatusBadRequest)
			return
		}
		endDate = &t
	}

	subscription := model.Subscription{
		ID:          id,
		ServiceName: request.ServiceName,
		Price:       request.Price,
		UserID:      request.UserID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	err = h.repo.Update(subscription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(toResponse(subscription))
}

// CalculateCost godoc
// @Summary Рассчитать стоимость подписок
// @Tags subscriptions
// @Produce json
// @Param user_id query string false "User ID"
// @Param service_name query string false "Service Name"
// @Param from query string false "YYYY-MM"
// @Param to query string false "YYYY-MM"
// @Success 200 {object} map[string]int
// @Router /subscriptions/cost [get]

func (h *Handler) CalculateCost(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("user_id")
	serviceName := r.URL.Query().Get("service_name")
	from := r.URL.Query().Get("from")
	to := r.URL.Query().Get("to")

	total, err := h.repo.CalculateCost(userID, serviceName, from, to)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]int{
		"total_cost": total,
	})
}

func toResponse(subscription model.Subscription) model.SubscriptionResponse {
	response := model.SubscriptionResponse{
		ID:          subscription.ID,
		ServiceName: subscription.ServiceName,
		Price:       subscription.Price,
		UserID:      subscription.UserID,
		StartDate:   subscription.StartDate.Format("01-2006"),
	}

	if subscription.EndDate != nil {
		end := subscription.EndDate.Format("01-2006")
		response.EndDate = &end
	}

	return response
}
