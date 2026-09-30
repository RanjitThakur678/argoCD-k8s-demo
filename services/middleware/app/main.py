"""FastAPI gateway/BFF: auth + request routing to the Go backend, and the home for future ML features
(recommendations/classification) that need Python's ecosystem rather than Go's.
"""
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from prometheus_fastapi_instrumentator import Instrumentator

from app.routers import cart, checkout, payments, products

app = FastAPI(title="ecom-middleware", version="1.0.0")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # tighten to the deployed frontend origin before production use
    allow_methods=["*"],
    allow_headers=["*"],
)

Instrumentator().instrument(app).expose(app, endpoint="/metrics")

app.include_router(products.router, prefix="/api/products", tags=["catalog", "reviews"])
app.include_router(cart.router, prefix="/api/cart", tags=["cart"])
app.include_router(checkout.router, prefix="/api/checkout", tags=["checkout"])
app.include_router(payments.router, prefix="/api/payments", tags=["payment"])


@app.get("/health")
def health():
    return {"status": "ok"}
