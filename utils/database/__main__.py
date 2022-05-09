import json
import sys

from pymongo.mongo_client import MongoClient

from .models import Output

MONGO_URL = ""


if __name__ == "__main__":

    client: MongoClient = MongoClient(MONGO_URL)

    for line in sys.stdin:
        output = Output(**json.loads(line))
        client.db.output.insert_one(output.dict())
