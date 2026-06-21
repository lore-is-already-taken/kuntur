from pydantic import BaseModel, Field


class FormPayload(BaseModel):
    name: str
    email: str
    message: str


class ContactResponse(BaseModel):
    id: str
    name: str
    email: str
    message: str
    attended: bool = Field(default=False)


class ContactAttendUpdate(BaseModel):
    attended: bool

