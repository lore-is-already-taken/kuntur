"""HTTP routes for the presentations (shows) feature.

Exposes the live presentations of the band. Each show is stored in the
``shows`` MongoDB collection and returned to the client as a ``ShowResponse``
object, which extends the public ``Show`` model with the document ``_id``.
"""

from typing import List

from fastapi import APIRouter, Depends

import app.db as db
from app.auth import ADMIN_RESPONSES, require_admin
from app.types.presentaciones import Show, ShowResponse

presentaciones_router = APIRouter()


@presentaciones_router.get(
    "/",
    summary="List all shows",
    description=(
        "Returns every show stored in the database, in MongoDB's natural "
        "order. Use the ``fecha.year`` and ``fecha.mes`` fields to sort "
        "client-side (oldest or newest first)."
    ),
    response_description="All shows in the collection, possibly empty.",
)
async def get_show() -> List[ShowResponse]:
    """List every show stored in the database.

    Returns:
        List[ShowResponse]: All shows found in the ``shows`` collection, each
        carrying its database ``_id`` exposed as ``id``. The list is empty
        when the collection has no documents.
    """
    collection = db.get_collection("shows")
    docs = await collection.find().to_list(length=None)
    return [_to_response(doc) for doc in docs]


@presentaciones_router.post(
    "/",
    status_code=201,
    response_model=ShowResponse,
    dependencies=[Depends(require_admin)],
    summary="Create a show",
    description=(
        "Admin endpoint. Persists a new show entry. The body is validated "
        "against the ``Show`` schema (``place`` + ``fecha``). The response "
        "includes the server-assigned ``id`` so the show can be referenced "
        "later."
    ),
    response_description="The created show, with its server-assigned id.",
    responses={
        **ADMIN_RESPONSES,
        422: {
            "description": "Validation Error — payload did not match the Show schema."
        },
    },
)
async def post_show(payload: Show) -> ShowResponse:
    """Create a new show entry.

    The request body is validated against the :class:`Show` schema (venue and
    date) and persisted as a new document in the ``shows`` collection. The
    server-assigned ``_id`` is returned in the response so the show can be
    referenced later.

    Args:
        payload (Show): Show to create. ``place`` and ``fecha`` are required;
            ``place.direction`` is optional.

    Returns:
        ShowResponse: The persisted show, including the server-assigned
            ``id`` (string form of the MongoDB ``ObjectId``).
    """
    collection = db.get_collection("shows")
    data = payload.model_dump()
    result = await collection.insert_one(data)
    return ShowResponse(id=str(result.inserted_id), **data)


def _to_response(doc: dict) -> ShowResponse:
    """Convert a raw MongoDB document into the public ``ShowResponse`` shape.

    Strips the internal ``_id`` field and re-exposes it as a string ``id`` so
    the payload is JSON-serializable and stable for client consumption.

    Args:
        doc (dict): A document as returned by the motor driver, carrying an
            ``_id`` of type :class:`bson.ObjectId`.

    Returns:
        ShowResponse: The same data, ready to be returned to the client.
    """
    return ShowResponse(
        id=str(doc["_id"]), **{k: v for k, v in doc.items() if k != "_id"}
    )
