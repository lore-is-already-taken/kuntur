"""Kuntur API — FastAPI entry point.

This module wires the MongoDB lifespan, registers the four feature routers,
and configures the public documentation surface (Swagger UI at ``/docs`` and
ReDoc at ``/redoc``). The OpenAPI specification is exposed at
``/openapi.json``.

The API is intentionally public and has no authentication. The contact form
endpoint is open to anyone; the rest are admin-style reads/writes intended
for the band to manage their own content.
"""
from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI

import app.db as db

from app.routes.biography import bio_router
from app.routes.contact import contact_router
from app.routes.members import member_router
from app.routes.show import show_router

load_dotenv()


@asynccontextmanager
async def lifespan(_: FastAPI):
    await db.connect()
    yield
    await db.disconnect()


app = FastAPI(
    title="Kuntur API",
    description=(
        "Backend HTTP API for the [Kuntur](https://kunturkantu.cl) band "
        "website.\n\n"
        "## Resources\n\n"
        "- **Shows** — upcoming and past live presentations.\n"
        "- **Members** — the band's integrantes and their instruments.\n"
        "- **Biography** — the band's bio (stub at the moment).\n"
        "- **Contact** — public form submissions and admin triage.\n\n"
        "## Interactive documentation\n\n"
        "- Swagger UI: [`/docs`](/docs)\n"
        "- ReDoc: [`/redoc`](/redoc)\n"
        "- OpenAPI JSON: [`/openapi.json`](/openapi.json)\n"
    ),
    version="0.1.0",
    contact={
        "name": "Kuntur",
        "url": "https://kunturkantu.cl",
        "email": "hola@kunturkantu.cl",
    },
    license_info={
        "name": "Private — all rights reserved.",
    },
    lifespan=lifespan,
    docs_url="/docs",
    redoc_url="/redoc",
    openapi_url="/openapi.json",
)

app.include_router(contact_router, prefix="/contact", tags=["Contact"])
app.include_router(member_router, prefix="/integrantes", tags=["Members"])
app.include_router(bio_router, prefix="/bio", tags=["Biography"])
app.include_router(show_router, prefix="/show", tags=["Shows"])


@app.get(
    "/",
    summary="Health check",
    description=(
        "Static hello-world response. Use ``/redoc`` or ``/docs`` for the "
        "real API documentation; use this endpoint to verify the service is "
        "up."
    ),
    response_description="A static JSON object with a ``message`` key.",
)
async def root():
    return {"message": "Hello World"}
