"""HTTP routes for the public contact form.

Stores submissions from the website's contact form in the ``contact`` MongoDB
collection and exposes a small admin surface to read them back and to mark
them as attended.

Public surface (anyone can submit):
    - ``POST /`` — submit the form.

Admin surface (intended for the band to triage submissions):
    - ``GET  /``             — list every submission.
    - ``GET  /{contact_id}`` — fetch a single submission.
    - ``PATCH /{contact_id}`` — toggle the ``attended`` flag.
"""
from bson import ObjectId
from fastapi import APIRouter, HTTPException
from pymongo import ReturnDocument

import app.db as db
from app.types.contact import ContactAttendUpdate, ContactResponse, FormPayload

contact_router = APIRouter()


@contact_router.post(
    "/",
    response_model=ContactResponse,
    status_code=201,
    summary="Submit the contact form",
    description=(
        "Public endpoint. Persists a new contact form submission. The body is "
        "validated against the ``FormPayload`` schema (name, email, message). "
        "The response includes the server-assigned ``id`` and the initial "
        "``attended: false`` flag."
    ),
    response_description="The created submission, with its server-assigned id and ``attended=false``.",
    responses={
        422: {
            "description": "Validation Error — payload did not match the FormPayload schema."
        },
    },
)
async def handle_form(payload: FormPayload) -> ContactResponse:
    """Persist a new contact form submission.

    The request body is validated against the :class:`FormPayload` schema and
    inserted as a new document in the ``contact`` collection. The
    server-assigned ``_id`` is returned in the response so the submission can
    be referenced later.

    Args:
        payload (FormPayload): The form fields submitted by the visitor. See
            the ``FormPayload`` schema for required and optional fields.

    Returns:
        ContactResponse: The persisted submission, including the
            server-assigned ``id`` (string form of the MongoDB ``ObjectId``).
    """
    collection = db.get_collection("contact")
    data = payload.model_dump()
    result = await collection.insert_one(data)
    return ContactResponse(id=str(result.inserted_id), **data)


@contact_router.get(
    "/",
    response_model=list[ContactResponse],
    summary="List all contact submissions",
    description=(
        "Admin endpoint. Returns every submission stored in the ``contact`` "
        "collection, in MongoDB's natural order. Intended for the band to "
        "triage incoming messages."
    ),
    response_description="All submissions in the collection, possibly empty.",
)
async def list_contacts() -> list[ContactResponse]:
    """List every contact submission stored in the database.

    Returns:
        list[ContactResponse]: All submissions in the ``contact`` collection,
            each carrying its database ``_id`` exposed as ``id``. The list is
            empty when the collection has no documents.
    """
    collection = db.get_collection("contact")
    docs = await collection.find().to_list(length=None)
    return [
        ContactResponse(
            id=str(doc["_id"]),
            **{k: v for k, v in doc.items() if k != "_id"},
        )
        for doc in docs
    ]


@contact_router.get(
    "/{contact_id}",
    response_model=ContactResponse,
    summary="Fetch a single contact submission",
    description=(
        "Admin endpoint. Returns the submission identified by ``contact_id``."
    ),
    response_description="The requested submission.",
    responses={
        400: {"description": "ID inválido — ``contact_id`` is not a valid ObjectId."},
        404: {"description": "Contacto no encontrado — no document matches the given id."},
    },
)
async def get_contact(contact_id: str) -> ContactResponse:
    """Fetch a single contact submission by id.

    Args:
        contact_id (str): The submission's ``id`` as returned by ``GET /`` or
            ``POST /``. Must be a valid MongoDB ``ObjectId`` string.

    Raises:
        HTTPException: ``400 ID inválido`` if ``contact_id`` is not a valid
            ``ObjectId`` string.
        HTTPException: ``404 Contacto no encontrado`` if no document matches
            the given ``contact_id``.

    Returns:
        ContactResponse: The requested submission.
    """
    collection = db.get_collection("contact")
    try:
        oid = ObjectId(contact_id)
    except Exception:
        raise HTTPException(status_code=400, detail="ID inválido")
    doc = await collection.find_one({"_id": oid})
    if doc is None:
        raise HTTPException(status_code=404, detail="Contacto no encontrado")
    return ContactResponse(
        id=str(doc["_id"]),
        **{k: v for k, v in doc.items() if k != "_id"},
    )


@contact_router.patch(
    "/{contact_id}",
    response_model=ContactResponse,
    summary="Set the attended flag on a submission",
    description=(
        "Admin endpoint. Sets the ``attended`` field to the value in the "
        "payload (not a toggle). Use ``true`` to mark handled, ``false`` to "
        "revert."
    ),
    response_description="The submission after the update, with the new ``attended`` value reflected.",
    responses={
        400: {"description": "ID inválido — ``contact_id`` is not a valid ObjectId."},
        404: {"description": "Contacto no encontrado — no document matches the given id."},
        422: {
            "description": "Validation Error — payload did not match the ContactAttendUpdate schema."
        },
    },
)
async def update_attended(
    contact_id: str, payload: ContactAttendUpdate
) -> ContactResponse:
    """Set the ``attended`` flag of a contact submission.

    This is a single-field, full-value update: the field is set to the value
    in ``payload.attended`` (not toggled). Use it to mark a submission as
    handled, or to revert a previous mark.

    Args:
        contact_id (str): The submission's ``id``. Must be a valid MongoDB
            ``ObjectId`` string.
        payload (ContactAttendUpdate): Carries the desired ``attended`` value.

    Raises:
        HTTPException: ``400 ID inválido`` if ``contact_id`` is not a valid
            ``ObjectId`` string.
        HTTPException: ``404 Contacto no encontrado`` if no document matches
            the given ``contact_id``.

    Returns:
        ContactResponse: The submission after the update, with the new
            ``attended`` value reflected.
    """
    collection = db.get_collection("contact")
    try:
        oid = ObjectId(contact_id)
    except Exception:
        raise HTTPException(status_code=400, detail="ID inválido")
    doc = await collection.find_one_and_update(
        {"_id": oid},
        {"$set": {"attended": payload.attended}},
        return_document=ReturnDocument.AFTER,
    )
    if doc is None:
        raise HTTPException(status_code=404, detail="Contacto no encontrado")
    return ContactResponse(
        id=str(doc["_id"]),
        **{k: v for k, v in doc.items() if k != "_id"},
    )
