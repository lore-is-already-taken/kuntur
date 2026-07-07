"""Admin authentication for write/triage endpoints.

This module is the single seam where admin authentication happens. Today it
validates a static bearer token against the ``ADMIN_TOKEN`` environment
variable; when a real token service (per-user tokens, expiry, rotation) is
introduced, only the internals of :func:`require_admin` change — the routes
that depend on it stay untouched.

Usage::

    @router.post("/", dependencies=[Depends(require_admin)])
    async def create_thing(...): ...
"""

import secrets
from os import getenv

from fastapi import Depends, HTTPException
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer

# OpenAPI documentation for the responses this dependency can produce.
# Spread into each protected endpoint's ``responses`` mapping.
ADMIN_RESPONSES = {
    401: {"description": "Unauthorized — missing or invalid admin bearer token."},
    503: {"description": "Service Unavailable — admin token is not configured."},
}

# auto_error=False so we control the error shape (and return 503 when the
# server itself is misconfigured instead of a misleading 403).
_bearer = HTTPBearer(auto_error=False)


def require_admin(
    credentials: HTTPAuthorizationCredentials | None = Depends(_bearer),
) -> None:
    """Validate the ``Authorization: Bearer <token>`` header for admin access.

    Declared through ``HTTPBearer`` so Swagger UI shows the global
    **Authorize** button — paste the raw token there once and every
    "Try it out" request carries it; no manual ``Bearer `` prefix needed.

    Fails closed: if ``ADMIN_TOKEN`` is not configured in the environment,
    every protected request is rejected with 503 rather than letting writes
    through unauthenticated.

    Args:
        credentials: Parsed ``Authorization`` header (scheme + token),
            injected by FastAPI; ``None`` when the header is missing or not
            a Bearer scheme.

    Raises:
        HTTPException: 503 if ``ADMIN_TOKEN`` is unset; 401 if the header is
            missing, malformed, or the token does not match.
    """
    expected = getenv("ADMIN_TOKEN", "")
    if not expected:
        raise HTTPException(
            status_code=503, detail="Admin token is not configured on the server."
        )

    if credentials is None or not secrets.compare_digest(
        credentials.credentials.strip(), expected
    ):
        raise HTTPException(
            status_code=401,
            detail="Invalid or missing admin token.",
            headers={"WWW-Authenticate": "Bearer"},
        )
