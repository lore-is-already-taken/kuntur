"""HTTP routes for the band members (integrantes) feature.

Exposes CRUD operations over the members of the group. Members are stored in
the ``integrantes`` MongoDB collection and returned to the client as
``MemberResponse`` objects, which extend the public ``MemberCreate`` model
with the document ``_id``.
"""

from bson import ObjectId
from fastapi import APIRouter, Depends, HTTPException

import app.db as db
from app.auth import ADMIN_RESPONSES, require_admin
from app.types.integrantes import MemberCreate, MemberResponse

integrantes_router = APIRouter()


@integrantes_router.get(
    "/",
    response_model=list[MemberResponse],
    summary="List all members",
    description="Returns every band member stored in the database, in MongoDB's natural order.",
    response_description="All members in the collection, possibly empty.",
)
async def get_integrantes() -> list[MemberResponse]:
    """List every member stored in the database.

    Returns:
        list[MemberResponse]: All members found in the ``integrantes``
            collection, each carrying its database ``_id`` exposed as ``id``.
            The list is empty when the collection has no documents.
    """
    collection = db.get_collection("integrantes")
    docs = await collection.find().to_list(length=None)
    return [_to_response(doc) for doc in docs]


@integrantes_router.post(
    "/",
    response_model=MemberResponse,
    status_code=201,
    dependencies=[Depends(require_admin)],
    summary="Create a member",
    description=(
        "Admin endpoint. Persists a new band member. The body is validated "
        "against the ``MemberCreate`` schema (name, description, and at "
        "least one instrument). The response includes the server-assigned "
        "``id``."
    ),
    response_description="The created member, with its server-assigned id.",
    responses={
        **ADMIN_RESPONSES,
        422: {
            "description": "Validation Error — payload did not match the MemberCreate schema."
        },
    },
)
async def post_integrantes(payload: MemberCreate) -> MemberResponse:
    """Create a new member entry.

    The request body is validated against the :class:`MemberCreate` schema and
    persisted as a new document in the ``integrantes`` collection. The
    server-assigned ``_id`` is returned in the response so the member can be
    referenced later for updates or deletion.

    Args:
        payload (MemberCreate): Member data to create. See the
            ``MemberCreate`` schema for required and optional fields.

    Returns:
        MemberResponse: The persisted member, including the server-assigned
            ``id`` (string form of the MongoDB ``ObjectId``).
    """
    collection = db.get_collection("integrantes")
    data = payload.model_dump()
    result = await collection.insert_one(data)
    return MemberResponse(id=str(result.inserted_id), **data)


@integrantes_router.put(
    "/{member_id}",
    response_model=MemberResponse,
    dependencies=[Depends(require_admin)],
    summary="Replace a member",
    description=(
        "Admin endpoint. **Full replacement**, not a partial update: every "
        "field of the stored document is overwritten with the contents of "
        "the request body. The internal ``_id`` is preserved."
    ),
    response_description="The member after the replacement.",
    responses={
        **ADMIN_RESPONSES,
        400: {"description": "ID inválido — ``member_id`` is not a valid ObjectId."},
        404: {
            "description": "Integrante no encontrado — no document matches the given id."
        },
        422: {
            "description": "Validation Error — payload did not match the MemberCreate schema."
        },
    },
)
async def put_integrante(member_id: str, payload: MemberCreate) -> MemberResponse:
    """Replace an existing member document with the supplied payload.

    Unlike a partial update, this is a full replacement: every field of the
    stored document is overwritten with the contents of ``payload``. The
    ``_id`` is preserved.

    Args:
        member_id (str): The member's ``id`` as returned by ``GET /`` or
            ``POST /``. Must be a valid MongoDB ``ObjectId`` string.
        payload (MemberCreate): The new member data that will fully replace
            the existing document.

    Raises:
        HTTPException: ``400 ID inválido`` if ``member_id`` is not a valid
            ``ObjectId`` string.
        HTTPException: ``404 Integrante no encontrado`` if no document
            matches the given ``member_id``.

    Returns:
        MemberResponse: The updated member after the replacement.
    """
    collection = db.get_collection("integrantes")
    try:
        oid = ObjectId(member_id)
    except Exception:
        raise HTTPException(status_code=400, detail="ID inválido")
    result = await collection.find_one_and_replace(
        {"_id": oid},
        payload.model_dump(),
        return_document=True,
    )
    if result is None:
        raise HTTPException(status_code=404, detail="Integrante no encontrado")
    return _to_response(result)


def _to_response(doc: dict) -> MemberResponse:
    """Convert a raw MongoDB document into the public ``MemberResponse`` shape.

    Strips the internal ``_id`` field and re-exposes it as a string ``id`` so
    the payload is JSON-serializable and stable for client consumption.

    Args:
        doc (dict): A document as returned by the motor driver, carrying an
            ``_id`` of type :class:`bson.ObjectId`.

    Returns:
        MemberResponse: The same data, ready to be returned to the client.
    """
    return MemberResponse(
        id=str(doc["_id"]), **{k: v for k, v in doc.items() if k != "_id"}
    )
