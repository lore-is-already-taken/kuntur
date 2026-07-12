from typing import List

from pydantic import BaseModel, ConfigDict, Field


class Instrument(BaseModel):
    """An instrument a band member plays."""

    type: str = Field(
        ...,
        description="Kind of instrument (e.g. ``string``, ``percussion``, ``voice``).",
        examples=["string", "percussion", "voice"],
    )
    name: str = Field(
        ...,
        description="Specific instrument name (e.g. ``charango``, ``quena``).",
        examples=["charango", "quena", "bombo"],
    )

    model_config = ConfigDict(
        json_schema_extra={"example": {"type": "string", "name": "charango"}}
    )


class MemberCreate(BaseModel):
    """Payload to create a new band member."""

    name: str = Field(
        ...,
        description="Member's full name.",
        examples=["Camila Quispe"],
    )
    description: str = Field(
        ...,
        description="Short bio or role description shown on the website.",
        examples=["Vocalist and charango player."],
    )
    instrument: List[Instrument] = Field(
        ...,
        description="Instruments the member plays (at least one).",
    )
    photo: str | None = Field(
        None,
        description=(
            "Public path of the member's portrait, served by the frontend "
            "(e.g. ``/static/img/integrantes/camila-quispe.webp``). ``None`` "
            "when no photo is available. The file must already exist among "
            "the frontend's embedded static assets — they ship at build "
            "time, so a new image requires a frontend rebuild."
        ),
        examples=["/static/img/integrantes/camila-quispe.webp"],
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "name": "Camila Quispe",
                "description": "Vocalist and charango player.",
                "instrument": [
                    {"type": "voice", "name": "lead vocals"},
                    {"type": "string", "name": "charango"},
                ],
                "photo": "/static/img/integrantes/camila-quispe.webp",
            }
        }
    )


class MemberResponse(MemberCreate):
    """A band member as returned by the API, with its server-assigned id."""

    id: str = Field(
        ...,
        description="Server-assigned MongoDB ``ObjectId``, string-encoded.",
        examples=["665f1a2b3c4d5e6f7a8b9c0d"],
    )
