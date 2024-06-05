package order

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID") // Предполагаем, что User-ID передается в заголовке
	switch r.Method {
	case http.MethodPost:
		h.createOrder(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) createOrder(w http.ResponseWriter, r *http.Request, userID string) {
	var dto CreateOrderDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	orderID, err := h.service.CreateOrder(r.Context(), userID, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"order_id": strconv.Itoa(orderID)})
}
