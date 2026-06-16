from os import getenv

from motor.motor_asyncio import AsyncIOMotorClient

_client: AsyncIOMotorClient | None = None


def _build_uri() -> str:
    user = getenv("DB_USER", "")
    password = getenv("DB_PASSWORD", "")
    host = getenv("DB_HOST", "127.0.0.1")
    port = getenv("DB_PORT", "27017")
    name = getenv("DB_NAME", "kuntur")
    auth_source = getenv("DB_AUTH_SOURCE", name)
    if user and password:
        return (
            f"mongodb://{user}:{password}@{host}:{port}/{name}?authSource={auth_source}"
        )
    return f"mongodb://{host}:{port}"


async def connect() -> None:
    global _client
    _client = AsyncIOMotorClient(_build_uri())


async def disconnect() -> None:
    if _client:
        _client.close()


def get_collection(name: str):
    assert _client is not None, "DB not connected — call connect() first"
    return _client[getenv("DB_NAME", "kuntur")][name]
