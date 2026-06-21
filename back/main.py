from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI
from pydantic import BaseModel

import app.db as db

from app.routes.biography import bio_router
from app.routes.members import member_router
from app.routes.shows import show_router

load_dotenv()


@asynccontextmanager
async def lifespan(_: FastAPI):
    await db.connect()
    yield
    await db.disconnect()


app = FastAPI(lifespan=lifespan)
app.include_router(member_router, prefix="/integrantes", tags=["Members"])
app.include_router(bio_router, prefix="/bio", tags=["Biography"])
app.include_router(show_router, prefix="/show", tags=["Shows"])


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