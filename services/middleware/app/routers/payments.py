"""Payment - write-heavy, must be idempotent. `Idempotency-Key` is required and forwarded as-is
to the Go backend, which is where the actual dedupe check against Postgres happens."""
from fastapi import APIRouter, Header, HTTPException
from pydantic import BaseModel

from app.config import backend_client

router = APIRouter()


class PaymentIn(BaseModel):
    order_id: str
    amount: float


@router.post("", status_code=201)
async def create_payment(payment: PaymentIn, idempotency_key: str | None = Header(default=None)):
    if not idempotency_key:
        raise HTTPException(status_code=400, detail="Idempotency-Key header is required")
    resp = await backend_client.post(
        "/api/payments",
        json=payment.model_dump(),
        headers={"Idempotency-Key": idempotency_key},
    )
    resp.raise_for_status()
    return resp.json()
