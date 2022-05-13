import pydantic

from typing import List


class PlayersSampleRow(pydantic.BaseModel):
    """Players sample row model."""

    id: str
    name: str


class Players(pydantic.BaseModel):
    """Players model."""

    online: int
    max: int
    sample: List[PlayersSampleRow] = []


class Output(pydantic.BaseModel):
    """Output model."""

    host: str
    port: int
    version: str
    protocol: int
    players: Players
    description: str
    date: str
