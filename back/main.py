"""Kuntur API — FastAPI entry point.

This module wires the MongoDB lifespan, registers the four feature routers,
and configures the public documentation surface (Swagger UI at ``/docs`` and
ReDoc at ``/redoc``). The OpenAPI specification is exposed at
``/openapi.json``.

Content reads and the contact form submission are public. Admin-style writes
and contact triage require a bearer token (see ``app.auth.require_admin``),
validated against the ``ADMIN_TOKEN`` environment variable.
"""

from contextlib import asynccontextmanager

from dotenv import load_dotenv
from fastapi import FastAPI

import app.db as db
from app.routes.biografia import biografia_router
from app.routes.contacto import contacto_router
from app.routes.integrantes import integrantes_router
from app.routes.presentaciones import presentaciones_router

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
        "- **Presentaciones** — upcoming and past live presentations.\n"
        "- **Integrantes** — the band's members and their instruments.\n"
        "- **Biografía** — the band's biography.\n"
        "- **Contacto** — public form submissions and admin triage.\n\n"
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

app.include_router(contacto_router, prefix="/contacto", tags=["Contacto"])
app.include_router(integrantes_router, prefix="/integrantes", tags=["Integrantes"])
app.include_router(biografia_router, prefix="/biografia", tags=["Biografía"])
app.include_router(
    presentaciones_router, prefix="/presentaciones", tags=["Presentaciones"]
)


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
