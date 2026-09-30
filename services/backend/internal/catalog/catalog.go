// Package catalog is the read-heavy domain: product tiles/dashboard/detail pages.
// Cache-aside against Redis in front of Postgres - falls back to in-memory seed data if neither is configured,
// so the service is runnable standalone before the datastores are wired up.
package catalog

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
)

type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Image string  `json:"image"`
}

const cacheKey = "catalog:products"
const cacheTTL = 60 * time.Second

var seed = []Product{
	{ID: "p1", Name: "Aurora Backpack", Price: 59.99, Image: "/img/backpack.jpg"},
	{ID: "p2", Name: "Nimbus Sneakers", Price: 89.00, Image: "/img/sneakers.jpg"},
	{ID: "p3", Name: "Halo Desk Lamp", Price: 34.50, Image: "/img/lamp.jpg"},
}

// List handles GET /api/products - the highest-traffic route in the app, hence cache-first.
func List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if platform.Cache != nil {
		if cached, err := platform.Cache.Get(ctx, cacheKey).Result(); err == nil {
			platform.CacheResult.WithLabelValues("catalog", "hit").Inc()
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cached))
			return
		}
	}
	platform.CacheResult.WithLabelValues("catalog", "miss").Inc()

	products := fetchFromDB(ctx)
	body, _ := json.Marshal(products)

	if platform.Cache != nil {
		platform.Cache.Set(ctx, cacheKey, body, cacheTTL)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Get handles GET /api/products/{id}.
func Get(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	for _, p := range fetchFromDB(r.Context()) {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "product not found", http.StatusNotFound)
}

func fetchFromDB(ctx context.Context) []Product {
	if platform.DB == nil {
		return seed
	}
	rows, err := platform.DB.Query(ctx, "SELECT id, name, price, image FROM products ORDER BY name")
	if err != nil {
		return seed
	}
	defer rows.Close()

	var products []Product
	for rows.Next() {
		var p Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Image); err == nil {
			products = append(products, p)
		}
	}
	if len(products) == 0 {
		return seed
	}
	return products
}

// Invalidate drops the cached product list - call this after any product write (admin update, restock, etc.).
func Invalidate(ctx context.Context) {
	if platform.Cache != nil {
		platform.Cache.Del(ctx, cacheKey)
	}
}
