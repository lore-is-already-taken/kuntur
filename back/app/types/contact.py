from pydantic import BaseModel, ConfigDict, Field


class FormPayload(BaseModel):
    """A submission to the public contact form."""

    name: str = Field(
        ...,
        description="Sender's name.",
        examples=["Lore Ramirez"],
    )
    email: str = Field(
        ...,
        description="Sender's email address. Used to reply to the message.",
        examples=["lore@example.com"],
    )
    message: str = Field(
        ...,
        description="The body of the message.",
        examples=["Hola, me encantaría contratarlos para un evento."],
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "name": "Lore Ramirez",
                "email": "lore@example.com",
                "message": "Hola, me encantaría contratarlos para un evento.",
            }
        }
    )


class ContactResponse(BaseModel):
    """A contact submission as returned by the API."""

    id: str = Field(
        ...,
        description="Server-assigned MongoDB ``ObjectId``, string-encoded.",
        examples=["665f1a2b3c4d5e6f7a8b9c0d"],
    )
    name: str = Field(..., description="Sender's name.")
    email: str = Field(..., description="Sender's email address.")
    message: str = Field(..., description="The body of the message.")
    attended: bool = Field(
        default=False,
        description="Whether the band has triaged this submission. "
        "Toggled via ``PATCH /contact/{contact_id}``.",
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "id": "665f1a2b3c4d5e6f7a8b9c0d",
                "name": "Lore Ramirez",
                "email": "lore@example.com",
                "message": "Hola, me encantaría contratarlos para un evento.",
                "attended": False,
            }
        }
    )


class ContactAttendUpdate(BaseModel):
    """Payload for ``PATCH /contact/{contact_id}`` — set the attended flag."""

    attended: bool = Field(
        ...,
        description="New value for the ``attended`` flag. Not a toggle — "
        "set to ``true`` to mark handled, ``false`` to revert.",
    )

    model_config = ConfigDict(json_schema_extra={"example": {"attended": True}})
