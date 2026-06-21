from fastapi import APIRouter

show_router = APIRouter()


@show_router.get("/")
def get_show():
    return {"message": "hola mi rey"}
