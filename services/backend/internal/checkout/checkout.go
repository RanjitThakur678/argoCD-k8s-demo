// Package checkout is write-heavy and consistency-critical: order creation must be transactional
// and is never cached or served stale - Postgres is the only source of truth here.
package checkout

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
)

type LineItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type OrderRequest struct {
	UserID string     `json:"user_id"`
	Items  []LineItem `json:"items"`
}

type Order struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Items     []LineItem `json:"items"`
	Total     float64    `json:"total"`
	Status    string     `json:"status"`
	CreatedAt time.Time  `json:"created_at"`
}

// Create handles POST /api/checkout - creates an order and its line items in a single transaction,
// so a partial write (order row without items, or vice versa) can never happen.
func Create(w http.ResponseWriter, r *http.Request) {
	var req OrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || len(req.Items) == 0 {
		http.Error(w, "invalid order", http.StatusBadRequest)
		return
	}

	var total float64
	for _, item := range req.Items {
		total += item.UnitPrice * float64(item.Quantity)
	}

	order := Order{
		UserID:    req.UserID,
		Items:     req.Items,
		Total:     total,
		Status:    "pending_payment",
		CreatedAt: time.Now().UTC(),
	}

	if platform.DB == nil {
		http.Error(w, "checkout store unavailable", http.StatusServiceUnavailable)
		return
	}

	ctx := r.Context()
	tx, err := platform.DB.Begin(ctx)
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO orders (user_id, total, status, created_at) VALUES ($1, $2, $3, $4) RETURNING id`,
		order.UserID, order.Total, order.Status, order.CreatedAt,
	).Scan(&order.ID)
	if err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	for _, item := range req.Items {
		_, err = tx.Exec(ctx,
			`INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES ($1, $2, $3, $4)`,
			order.ID, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			http.Error(w, "failed to save order items", http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "failed to commit order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}
