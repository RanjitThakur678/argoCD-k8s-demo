// Package payment is write-heavy and must be idempotent - retried requests (client timeout, network
// blip) must never double-charge. We never store card data ourselves; a real gateway (Stripe/Razorpay)
// would be called from ChargeGateway, keeping this service out of PCI card-data scope.
package payment

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
)

type PaymentRequest struct {
	OrderID string  `json:"order_id"`
	Amount  float64 `json:"amount"`
}

type Payment struct {
	ID        string    `json:"id"`
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// Create handles POST /api/payments. Requires an `Idempotency-Key` header - if a payment with that
// key already exists, the existing record is returned instead of charging again.
func Create(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := r.Header.Get("Idempotency-Key")
	if idempotencyKey == "" {
		http.Error(w, "Idempotency-Key header is required", http.StatusBadRequest)
		return
	}

	var req PaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.OrderID == "" {
		http.Error(w, "invalid payment request", http.StatusBadRequest)
		return
	}

	if platform.DB == nil {
		http.Error(w, "payment store unavailable", http.StatusServiceUnavailable)
		return
	}
	ctx := r.Context()

	var existing Payment
	err := platform.DB.QueryRow(ctx,
		`SELECT id, order_id, amount, status, created_at FROM payments WHERE idempotency_key = $1`,
		idempotencyKey,
	).Scan(&existing.ID, &existing.OrderID, &existing.Amount, &existing.Status, &existing.CreatedAt)
	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(existing)
		return
	}

	status := chargeGateway(req)

	payment := Payment{
		OrderID:   req.OrderID,
		Amount:    req.Amount,
		Status:    status,
		CreatedAt: time.Now().UTC(),
	}
	err = platform.DB.QueryRow(ctx,
		`INSERT INTO payments (order_id, amount, status, idempotency_key, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		payment.OrderID, payment.Amount, payment.Status, idempotencyKey, payment.CreatedAt,
	).Scan(&payment.ID)
	if err != nil {
		http.Error(w, "failed to record payment", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payment)
}

// chargeGateway is a stub for a real gateway call (Stripe/Razorpay) - swap this out, never add card storage here.
func chargeGateway(req PaymentRequest) string {
	if req.Amount <= 0 {
		return "failed"
	}
	return "succeeded"
}
