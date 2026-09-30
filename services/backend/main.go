// Command backend is the Go core service: catalog, cart, reviews, checkout, and payment domains,
// run as one modular monolith (simpler ops at this scale) behind the FastAPI gateway.
package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/cart"
	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/catalog"
	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/checkout"
	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/payment"
	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/platform"
	"github.com/RanjitThakur678/argoCD-k8s-demo/services/backend/internal/reviews"
)

func main() {
	platform.InitDB()
	platform.InitCache()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
	mux.Handle("GET /metrics", platform.MetricsHandler())

	// Catalog - read-heavy, cache-aside.
	mux.HandleFunc("GET /api/products", catalog.List)
	mux.HandleFunc("GET /api/products/{id}", catalog.Get)

	// Reviews - read+write heavy, cached listing invalidated on write.
	mux.HandleFunc("GET /api/products/{id}/reviews", reviews.List)
	mux.HandleFunc("POST /api/products/{id}/reviews", reviews.Create)

	// Cart - write-heavy, Redis-only, ephemeral.
	mux.HandleFunc("GET /api/cart", cart.Get)
	mux.HandleFunc("POST /api/cart/items", cart.AddItem)
	mux.HandleFunc("DELETE /api/cart/items/{productID}", cart.RemoveItem)

	// Checkout / Payment - write-heavy, transactional, never cached.
	mux.HandleFunc("POST /api/checkout", checkout.Create)
	mux.HandleFunc("POST /api/payments", payment.Create)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      withMetrics(mux),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("backend listening on :%s", port)
	log.Fatal(srv.ListenAndServe())
}

// withMetrics records request duration per route/status so Grafana dashboards can show per-domain latency.
func withMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		platform.RequestDuration.WithLabelValues(domainOf(r.URL.Path), r.Method, http.StatusText(rec.status)).
			Observe(time.Since(start).Seconds())
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func domainOf(path string) string {
	switch {
	case len(path) >= len("/api/cart") && path[:len("/api/cart")] == "/api/cart":
		return "cart"
	case len(path) >= len("/api/checkout") && path[:len("/api/checkout")] == "/api/checkout":
		return "checkout"
	case len(path) >= len("/api/payments") && path[:len("/api/payments")] == "/api/payments":
		return "payment"
	case len(path) >= len("/api/products") && path[:len("/api/products")] == "/api/products":
		return "catalog"
	default:
		return "other"
	}
}
