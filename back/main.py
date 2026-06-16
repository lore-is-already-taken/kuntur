from contextlib import asynccontextmanager

from bson import ObjectId
from dotenv import load_dotenv
from fastapi import FastAPI, HTTPException
from pydantic import BaseModel

import app.db as db
from app.types.integrantes import MemberCreate, MemberResponse

load_dotenv()


@asynccontextmanager
async def lifespan(_: FastAPI):
    await db.connect()
    yield
    await db.disconnect()


app = FastAPI(lifespan=lifespan)


class FormPayload(BaseModel):
    name: str
    email: str
    message: str


@app.get("/")
async def root():
    return {"message": "Hello World"}


@app.post("/contacto")
async def handle_form(payload: FormPayload):
    print(payload)
    return {"received": payload}


@app.get("/bio")
async def bio():
    return {"message": "This is important information"}


@app.post("/bio")
async def update_bio():
    pass


@app.get("/integrantes", response_model=list[MemberResponse])
async def get_integrantes():
    collection = db.get_collection("integrantes")
    docs = await collection.find().to_list(length=None)
    return [_to_response(doc) for doc in docs]


@app.post("/integrantes", response_model=MemberResponse, status_code=201)
async def post_integrantes(payload: MemberCreate):
    collection = db.get_collection("integrantes")
    data = payload.model_dump()
    result = await collection.insert_one(data)
    return MemberResponse(id=str(result.inserted_id), **data)


@app.put("/integrantes/{member_id}", response_model=MemberResponse)
async def put_integrante(member_id: str, payload: MemberCreate):
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
    return MemberResponse(
        id=str(doc["_id"]), **{k: v for k, v in doc.items() if k != "_id"}
    )
