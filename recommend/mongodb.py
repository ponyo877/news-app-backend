
from pymongo import MongoClient
from pymongo import DESCENDING
from pymongo import ASCENDING

class MongoModel(object):

    def __init__(self, hostName, dbName, collectionName):
        self.client = MongoClient(hostName)
        self.db = self.client[dbName]
        self.collection = self.db.get_collection(collectionName)

    def find_one(self, projection=None, filter=None, sort=None):
        return self.collection.find_one(projection=projection, filter=filter,sort=sort)

    def find(self, projection=None, filter=None, sort=None):
        return self.collection.find(projection=projection, filter=filter, sort=sort)

    def update_many(self, filter, update):
        return self.collection.update_many(filter,update)
