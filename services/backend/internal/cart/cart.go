// Package cart is the write-heavy, ephemeral domain: add/remove/update line items per user session.
// Backed entirely by Redis (a hash per user), never Postgres - cart contents don't need durability
// or transactional guarantees the way orders/payments do, and TTL naturally expires abandoned carts.
package cart

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
)

type Item struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

const cartTTL = 24 * time.Hour

func key(userID string) string { return "cart:" + userID }

func userID(r *http.Request) string {
	if id := r.Header.Get("X-User-Id"); id != "" {
		return id
	}
	return "anonymous"
}

// Get handles GET /api/cart.
func Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	uid := userID(r)
	w.Header().Set("Content-Type", "application/json")

	if platform.Cache == nil {
		json.NewEncoder(w).Encode([]Item{})
		return
	}

	raw, err := platform.Cache.HGetAll(ctx, key(uid)).Result()
	if err != nil {
		http.Error(w, "cart unavailable", http.StatusServiceUnavailable)
		return
	}

	items := make([]Item, 0, len(raw))
	for productID, qtyStr := range raw {
		qty, _ := strconv.Atoi(qtyStr)
		items = append(items, Item{ProductID: productID, Quantity: qty})
	}
	json.NewEncoder(w).Encode(items)
}

// AddItem handles POST /api/cart/items - increments quantity for a product in the caller's cart.
func AddItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil || item.ProductID == "" {
		http.Error(w, "invalid item", http.StatusBadRequest)
		return
	}
	if item.Quantity <= 0 {
		item.Quantity = 1
	}
	if platform.Cache == nil {
		http.Error(w, "cart store unavailable", http.StatusServiceUnavailable)
		return
	}

	ctx := r.Context()
	k := key(userID(r))
	if err := platform.Cache.HIncrBy(ctx, k, item.ProductID, int64(item.Quantity)).Err(); err != nil {
		http.Error(w, "failed to update cart", http.StatusInternalServerError)
		return
	}
	platform.Cache.Expire(ctx, k, cartTTL)
	w.WriteHeader(http.StatusNoContent)
}

// RemoveItem handles DELETE /api/cart/items/{productID}.
func RemoveItem(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("productID")
	if platform.Cache == nil {
		http.Error(w, "cart store unavailable", http.StatusServiceUnavailable)
		return
	}
	platform.Cache.HDel(r.Context(), key(userID(r)), productID)
	w.WriteHeader(http.StatusNoContent)
}
