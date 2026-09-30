"""Catalog + reviews - proxies to the Go backend. Also the placeholder home for a future
recommendation/classification model (Python's ML ecosystem is why this layer exists).
"""
from typing import Optional

from fastapi import APIRouter, HTTPException
from pydantic import BaseModel

from app.config import backend_client

router = APIRouter()


class ReviewIn(BaseModel):
    author: str
    rating: int
    comment: Optional[str] = None


@router.get("")
async def list_products():
    resp = await backend_client.get("/api/products")
    resp.raise_for_status()
    return resp.json()


@router.get("/{product_id}")
async def get_product(product_id: str):
    resp = await backend_client.get(f"/api/products/{product_id}")
    if resp.status_code == 404:
        raise HTTPException(status_code=404, detail="product not found")
    resp.raise_for_status()
    return resp.json()


@router.get("/{product_id}/reviews")
async def list_reviews(product_id: str):
    resp = await backend_client.get(f"/api/products/{product_id}/reviews")
    resp.raise_for_status()
    return resp.json()


@router.post("/{product_id}/reviews", status_code=201)
async def create_review(product_id: str, review: ReviewIn):
    resp = await backend_client.post(f"/api/products/{product_id}/reviews", json=review.model_dump())
    resp.raise_for_status()
    return resp.json()


@router.get("/{product_id}/recommendations")
async def recommendations(product_id: str):
    # TODO: replace with a real model (collaborative filtering / content-based) served from this layer.
    # Placeholder: return up to 3 other catalog products so the frontend has something to render.
    resp = await backend_client.get("/api/products")
    resp.raise_for_status()
    products = [p for p in resp.json() if p["id"] != product_id]
    return products[:3]
