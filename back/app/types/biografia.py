from pydantic import BaseModel, ConfigDict, Field


class Biography(BaseModel):
    """The band's biography — a single document, not a collection of items."""

    resume: str = Field(
        ...,
        description="The biography text shown on the website's /biografia page.",
        examples=[
            "Kuntur nace en los Andes: charango, quena y bombo al servicio "
            "de la música tradicional."
        ],
    )

    model_config = ConfigDict(
        json_schema_extra={
            "example": {
                "resume": (
                    "Kuntur nace en los Andes: charango, quena y bombo al "
                    "servicio de la música tradicional."
                )
            }
        }
    )
