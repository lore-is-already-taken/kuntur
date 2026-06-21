from pydantic import BaseModel


class FormPayload(BaseModel):
    name: str
    email: str
    message: str


class ContactResponse(BaseModel):
    id: str
    name: str
    email: str
    message: str
