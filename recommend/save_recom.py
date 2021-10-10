from recommender import RecommendModel
from mongodb import MongoModel
from bson.objectid import ObjectId
import json
import logging



if __name__ == "__main__":

    formatter = '%(levelname)s : %(asctime)s : %(message)s'
    logging.basicConfig(level=logging.INFO, format=formatter)
    _logger = logging.getLogger(__name__)

    _logger.info('%s', 'Start Create Instance of RecommendModel')
    f = open('config_mongo.json', 'r')
    _mongo_conf = json.load(f)
    hostName = str(_mongo_conf['host']) + ':' + str(_mongo_conf['port'])
    mongo = MongoModel(hostName, 'newsdb', 'article_col')
    recom = RecommendModel()
    _logger.info('%s', 'Finish Create Instance of RecommendModel')

    target_items = mongo.find(filter={'acquired': False})
    recom.put_recom_items(target_items)
    # success_id_list_tmp = recom.put_recom_items(target_items)
    # success_id_list = [ObjectId(_id) for _id in success_id_list_tmp]

    mongo.update_many({'acquired': False}, {'$set':{'acquired': True}})
    # mongo.update_many({'_id': {'$in': success_id_list}}, {'$set': {'acquired': True}})
    _logger.info('%s', 'Finish Method of mongo.update_many')

