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

from fastapi import Header, HTTPException

# OpenAPI documentation for the responses this dependency can produce.
# Spread into each protected endpoint's ``responses`` mapping.
ADMIN_RESPONSES = {
    401: {"description": "Unauthorized — missing or invalid admin bearer token."},
    503: {"description": "Service Unavailable — admin token is not configured."},
}


def require_admin(authorization: str | None = Header(default=None)) -> None:
    """Validate the ``Authorization: Bearer <token>`` header for admin access.

    Fails closed: if ``ADMIN_TOKEN`` is not configured in the environment,
    every protected request is rejected with 503 rather than letting writes
    through unauthenticated.

    Args:
        authorization: Raw ``Authorization`` header, injected by FastAPI.

    Raises:
        HTTPException: 503 if ``ADMIN_TOKEN`` is unset; 401 if the header is
            missing, malformed, or the token does not match.
    """
    expected = getenv("ADMIN_TOKEN", "")
    if not expected:
        raise HTTPException(
            status_code=503, detail="Admin token is not configured on the server."
        )

    scheme, _, token = (authorization or "").partition(" ")
    if scheme.lower() != "bearer" or not secrets.compare_digest(
        token.strip(), expected
    ):
        raise HTTPException(
            status_code=401,
            detail="Invalid or missing admin token.",
            headers={"WWW-Authenticate": "Bearer"},
        )
