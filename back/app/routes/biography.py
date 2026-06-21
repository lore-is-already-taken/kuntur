from fastapi import APIRouter

bio_router = APIRouter()


@bio_router.get("/")
async def bio():
    return {"message": "this is important information"}


@bio_router.post("/")
async def update_bio():
    pass
