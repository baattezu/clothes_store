package cart

import (
	"encoding/json"
	"net/http"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) HandleCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Header.Get("User-ID") // Предполагаем, что User-ID передается в заголовке
	switch r.Method {
	case http.MethodGet:
		h.getCart(w, r, userID)
	case http.MethodPost:
		h.addToCart(w, r, userID)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *Handler) getCart(w http.ResponseWriter, r *http.Request, userID string) {
	cartItems, err := h.service.GetCartItems(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cartItems)
}

func (h *Handler) addToCart(w http.ResponseWriter, r *http.Request, userID string) {
	var dto AddToCartDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := h.service.AddToCart(r.Context(), userID, dto)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
