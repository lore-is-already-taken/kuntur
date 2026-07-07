from pydantic import BaseModel, ConfigDict, Field


class Fecha(BaseModel):
    """A show's month + year bucket (e.g. ``Jun 2026``)."""

    mes: str = Field(
        ...,
        description="Month abbreviation (three letters, capitalized).",
        examples=["Jun", "Ago", "Dic"],
    )
    year: int = Field(
        ...,
        description="Four-digit year.",
        examples=[2026, 2027],
    )

    model_config = ConfigDict(
        json_schema_extra={"example": {"mes": "Jun", "year": 2026}}
    )


class Place(BaseModel):
    """The venue where a show takes place."""

    name: str = Field(
        ...,
        description="Name of the venue (e.g. ``Teatro Municipal``).",
        examples=["Espacio Cultural Centro", "Teatro Municipal"],
    )
    city: str = Field(
        ...,
        description="City where the venue is located.",
        examples=["Santiago", "Valparaíso"],
    )
    country: str = Field(
        ...,
        description="Country where the venue is located.",
        examples=["Chile"],
    )
    direction: str | None = Field(
        default=None,
        description="Optional street address. Not required for the listing view.",
        examples=["Plaza Victoria s/n"],
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "name": "Teatro Municipal",
                "city": "Valparaíso",
                "country": "Chile",
                "direction": "Plaza Victoria s/n",
            }
        }
    )


class Show(BaseModel):
    """A live presentation: when and where."""

    place: Place = Field(
        ...,
        description="The venue and city of the show.",
    )
    fecha: Fecha = Field(
        ...,
        description="The month and year of the show.",
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "place": {
                    "name": "Espacio Cultural Centro",
                    "city": "Santiago",
                    "country": "Chile",
                },
                "fecha": {"mes": "Jun", "year": 2026},
            }
        }
    )


class ShowResponse(Show):
    """A show as returned by the API, with its server-assigned id."""

    id: str = Field(
        ...,
        description="Server-assigned MongoDB ``ObjectId``, string-encoded.",
        examples=["665f1a2b3c4d5e6f7a8b9c0d"],
    )
