import json
import sys

from pymongo.mongo_client import MongoClient

from .models import Output

MONGO_URL = ""
MAX_SIZE = 20


if __name__ == "__main__":

    client: MongoClient = MongoClient(MONGO_URL)

    results = []

    for line in sys.stdin:

        try:
            results.append(Output(**json.loads(line)).dict())
        except Exception:
            continue

        if len(results) == MAX_SIZE:
            client.db.servers.insert_many(results)
            results = []
