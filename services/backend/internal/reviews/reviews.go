// Package reviews is read-and-write heavy: customers post reviews often, and every product page reads them.
// Writes go straight to Postgres; the listing read is cached with a short TTL and invalidated on new writes.
package reviews

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
)

type Review struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	Author    string    `json:"author"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

const cacheTTLPrefix = "reviews:"
const cacheTTL = 30 * time.Second

// List handles GET /api/products/{id}/reviews.
func List(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	ctx := r.Context()
	cacheKey := cacheTTLPrefix + productID

	if platform.Cache != nil {
		if cached, err := platform.Cache.Get(ctx, cacheKey).Result(); err == nil {
			platform.CacheResult.WithLabelValues("reviews", "hit").Inc()
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(cached))
			return
		}
	}
	platform.CacheResult.WithLabelValues("reviews", "miss").Inc()

	list := fetch(ctx, productID)
	body, _ := json.Marshal(list)
	if platform.Cache != nil {
		platform.Cache.Set(ctx, cacheKey, body, cacheTTL)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// Create handles POST /api/products/{id}/reviews.
func Create(w http.ResponseWriter, r *http.Request) {
	productID := r.PathValue("id")
	var rev Review
	if err := json.NewDecoder(r.Body).Decode(&rev); err != nil || rev.Author == "" {
		http.Error(w, "invalid review", http.StatusBadRequest)
		return
	}
	rev.ProductID = productID
	rev.CreatedAt = time.Now().UTC()

	ctx := r.Context()
	if platform.DB == nil {
		http.Error(w, "reviews store unavailable", http.StatusServiceUnavailable)
		return
	}
	err := platform.DB.QueryRow(ctx,
		`INSERT INTO reviews (product_id, author, rating, comment, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		rev.ProductID, rev.Author, rev.Rating, rev.Comment, rev.CreatedAt,
	).Scan(&rev.ID)
	if err != nil {
		http.Error(w, "failed to save review", http.StatusInternalServerError)
		return
	}

	// Invalidate the cached listing so the new review is visible on the next read.
	if platform.Cache != nil {
		platform.Cache.Del(ctx, cacheTTLPrefix+productID)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rev)
}

func fetch(ctx context.Context, productID string) []Review {
	if platform.DB == nil {
		return []Review{}
	}
	rows, err := platform.DB.Query(ctx,
		`SELECT id, product_id, author, rating, comment, created_at
		 FROM reviews WHERE product_id = $1 ORDER BY created_at DESC`, productID)
	if err != nil {
		return []Review{}
	}
	defer rows.Close()

	list := []Review{}
	for rows.Next() {
		var rev Review
		if err := rows.Scan(&rev.ID, &rev.ProductID, &rev.Author, &rev.Rating, &rev.Comment, &rev.CreatedAt); err == nil {
			list = append(list, rev)
		}
	}
	return list
}
