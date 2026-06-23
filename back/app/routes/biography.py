"""HTTP routes for the band's biography (bio).

.. warning::
    This module is currently a **stub**. Both endpoints exist so the router
    surface and the frontend wiring can be exercised, but neither reads nor
    writes a real biography document. The handlers return placeholder data and
    must be replaced with the real implementation before going to production.
"""
from fastapi import APIRouter

bio_router = APIRouter()


@bio_router.get(
    "/",
    summary="Get the band's biography (stub)",
    description=(
        "**Stub.** Returns a hard-coded placeholder message instead of the "
        "real biography. Replace with a MongoDB read once the bio schema "
        "and collection exist."
    ),
    response_description="Placeholder payload — not the real biography.",
)
async def bio():
    """Return the band's biography (placeholder).

    .. warning::
        Stub. Returns a hard-coded message instead of the real biography.
        Replace with a MongoDB read once the biography schema and collection
        exist.

    Returns:
        dict: Placeholder payload with a single ``message`` key.
    """
    return {"message": "this is important information"}


@bio_router.post(
    "/",
    summary="Update the band's biography (stub)",
    description=(
        "**Stub.** Accepts no input and writes nothing. Replace with a "
        "MongoDB upsert keyed by a single biography document once the schema "
        "exists."
    ),
    response_description="No content — endpoint is a stub.",
    responses={
        501: {"description": "Not Implemented — this endpoint is a placeholder."},
    },
)
async def update_bio():
    """Update the band's biography (placeholder).

    .. warning::
        Stub. Accepts no input and writes nothing. Replace with a MongoDB
        upsert keyed by a single biography document once the schema exists.
    """
    pass
