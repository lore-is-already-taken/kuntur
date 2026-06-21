from bson import ObjectId
from fastapi import APIRouter, HTTPException
from pymongo import ReturnDocument

import app.db as db
from app.types.contact import ContactAttendUpdate, ContactResponse, FormPayload

contact_router = APIRouter()


@contact_router.post("/", response_model=ContactResponse, status_code=201)
async def handle_form(payload: FormPayload) -> ContactResponse:
    collection = db.get_collection("contact")
    data = payload.model_dump()
    result = await collection.insert_one(data)
    return ContactResponse(id=str(result.inserted_id), **data)


@contact_router.get("/", response_model=list[ContactResponse])
async def list_contacts() -> list[ContactResponse]:
    collection = db.get_collection("contact")
    docs = await collection.find().to_list(length=None)
    return [
        ContactResponse(
            id=str(doc["_id"]),
            **{k: v for k, v in doc.items() if k != "_id"},
        )
        for doc in docs
    ]


@contact_router.get("/{contact_id}", response_model=ContactResponse)
async def get_contact(contact_id: str) -> ContactResponse:
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


@contact_router.patch("/{contact_id}", response_model=ContactResponse)
async def update_attended(
    contact_id: str, payload: ContactAttendUpdate
) -> ContactResponse:
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
