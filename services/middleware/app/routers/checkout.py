"""Checkout - write-heavy, consistency-critical. Straight pass-through to the Go backend's
transactional order creation; no caching, no retries that could double-submit."""
from fastapi import APIRouter
from pydantic import BaseModel

from app.config import backend_client

router = APIRouter()


class LineItemIn(BaseModel):
    product_id: str
    quantity: int
    unit_price: float


class OrderIn(BaseModel):
    user_id: str
    items: list[LineItemIn]


@router.post("", status_code=201)
async def create_order(order: OrderIn):
    resp = await backend_client.post("/api/checkout", json=order.model_dump())
    resp.raise_for_status()
    return resp.json()
