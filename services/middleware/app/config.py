"""Environment-driven configuration - kept in one place so every router shares the same backend client."""
import os

import httpx

BACKEND_URL = os.getenv("BACKEND_URL", "http://backend:8080")
REQUEST_TIMEOUT_SECONDS = float(os.getenv("BACKEND_TIMEOUT_SECONDS", "5"))

# Shared client (connection pooling) - created once at import time, reused by every router.
backend_client = httpx.AsyncClient(base_url=BACKEND_URL, timeout=REQUEST_TIMEOUT_SECONDS)
