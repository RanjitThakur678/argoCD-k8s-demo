"""Cart - write-heavy, thin proxy to the Go backend's Redis-backed cart. `X-User-Id` is forwarded
as-is; swap in real session/auth extraction here once auth is added to this layer.
"""
from fastapi import APIRouter, Header, Response
from pydantic import BaseModel

from app.config import backend_client

router = APIRouter()


class CartItemIn(BaseModel):
    product_id: str
    quantity: int = 1


def _headers(user_id: str | None) -> dict:
    return {"X-User-Id": user_id} if user_id else {}


@router.get("")
async def get_cart(x_user_id: str | None = Header(default=None)):
    resp = await backend_client.get("/api/cart", headers=_headers(x_user_id))
    resp.raise_for_status()
    return resp.json()


@router.post("/items", status_code=204)
async def add_item(item: CartItemIn, x_user_id: str | None = Header(default=None)):
    resp = await backend_client.post("/api/cart/items", json=item.model_dump(), headers=_headers(x_user_id))
    resp.raise_for_status()
    return Response(status_code=204)


@router.delete("/items/{product_id}", status_code=204)
async def remove_item(product_id: str, x_user_id: str | None = Header(default=None)):
    resp = await backend_client.delete(f"/api/cart/items/{product_id}", headers=_headers(x_user_id))
    resp.raise_for_status()
    return Response(status_code=204)
