"""HTTP routes for the band's biography (biografía).

The biography is a singleton: the ``biografia`` collection holds at most one
document. ``GET`` returns it (404 while unset) and ``POST`` upserts it, so
there is never more than one biography to reason about.
"""

from fastapi import APIRouter, Depends, HTTPException

import app.db as db
from app.auth import ADMIN_RESPONSES, require_admin
from app.types.biografia import Biography

biografia_router = APIRouter()


@biografia_router.get(
    "/",
    response_model=Biography,
    summary="Get the band's biography",
    description=(
        "Returns the single biography document. Responds 404 until a "
        "biography has been created via ``POST /biografia/``."
    ),
    response_description="The biography text.",
    responses={
        404: {"description": "Not Found — no biography has been created yet."},
    },
)
async def get_biografia() -> Biography:
    """Return the band's biography.

    Returns:
        Biography: The singleton biography document.

    Raises:
        HTTPException: 404 if no biography document exists yet.
    """
    collection = db.get_collection("biografia")
    doc = await collection.find_one({})
    if doc is None:
        raise HTTPException(status_code=404, detail="Biography not set yet.")
    return Biography(resume=doc.get("resume", ""))


@biografia_router.post(
    "/",
    response_model=Biography,
    dependencies=[Depends(require_admin)],
    summary="Create or replace the band's biography",
    description=(
        "Admin endpoint. Upserts the singleton biography document: creates "
        "it on first call, replaces the text on subsequent calls."
    ),
    response_description="The biography as stored.",
    responses={
        **ADMIN_RESPONSES,
        422: {
            "description": "Validation Error — payload did not match the Biography schema."
        },
    },
)
async def post_biografia(payload: Biography) -> Biography:
    """Create or replace the band's biography.

    Args:
        payload: The new biography text.

    Returns:
        Biography: The stored biography (echoes the payload).
    """
    collection = db.get_collection("biografia")
    await collection.update_one({}, {"$set": payload.model_dump()}, upsert=True)
    return payload
