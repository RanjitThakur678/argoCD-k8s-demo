// Package platform wires up shared infrastructure (Postgres, Redis, metrics) for all domain packages.
package platform

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

// DB is a pgx connection pool; nil if DATABASE_URL was unset or unreachable at startup,
// domain packages must fall back to seed data in that case (keeps the demo runnable with no Postgres).
var DB *pgxpool.Pool

// Cache is a Redis client; nil if REDIS_ADDR was unset or unreachable at startup.
var Cache *redis.Client

func InitDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("platform: DATABASE_URL not set, catalog/reviews/checkout run without persistence")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Printf("platform: postgres connect failed: %v", err)
		return
	}
	if err := pool.Ping(ctx); err != nil {
		log.Printf("platform: postgres ping failed: %v", err)
		return
	}
	DB = pool
	log.Println("platform: connected to postgres")
}

func InitCache() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		log.Println("platform: REDIS_ADDR not set, catalog/cart run without cache")
		return
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("platform: redis ping failed: %v", err)
		return
	}
	Cache = client
	log.Println("platform: connected to redis")
}
