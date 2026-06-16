from typing import List

from pydantic import BaseModel


class Instrument(BaseModel):
    type: str
    name: str


class MemberCreate(BaseModel):
    name: str
    description: str
    instrument: List[Instrument]


class MemberResponse(MemberCreate):
    id: str
