from recommender import RecommendModel
from mongodb import MongoModel
from flask import Flask, jsonify, request
from bson.objectid import ObjectId
from bson import json_util
from datetime import datetime
import logging
import json

app = Flask(__name__)
app.config['JSON_AS_ASCII'] = False

@app.route('/recom/<_id>')
def get_recom(_id):
    formatter = '%(levelname)s : %(asctime)s : %(message)s'
    logging.basicConfig(level=logging.INFO, format=formatter)
    _logger = logging.getLogger(__name__)

    f = open('config_mongo.json', 'r')
    _mongo_conf = json.load(f)
    hostName = str(_mongo_conf['host']) + ':' + str(_mongo_conf['port'])
    mongo = MongoModel(hostName, 'newsdb', 'article_col')
    _logger.info('%s', 'Finish Setup Mongodb')
    recom = RecommendModel()
    _logger.info('%s', 'Finish Setup RecommendModel')
    recom_items, dist_list = recom.get_recom_items(_id)
    if recom_items:
        recom_records = [mongo.find_one(filter={'_id': ObjectId(recom_id)}) for recom_id in recom_items]
        for recom_record, record_dist in zip(recom_records, dist_list):
            try:
                recom_record['_id'] = str(recom_record['_id'])
                # TODO: 2021-03-07T20:00:04+09:00に合わせる
                tm = datetime.strptime(recom_record['publishedAt'], '%Y-%m-%dT%H:%M:%S%z')
                recom_record['publishedAt'] = tm.strftime('%Y-%m-%d %H:%M')
                recom_record['distance'] = record_dist
                _logger.info('%s', 'recom_record[\'titles\']:' + recom_record['titles'])
            except TypeError:
                pass
    else:
        recom_records = None
    _logger.info('%s', 'Finish Load recom_items')
    return jsonify({'data': recom_records}) # json.dumps({'data': recom_records}, default=json_util.default)

@app.route('/personal')
def get_personal():
    _ids_str = request.args.get('ids')
    formatter = '%(levelname)s : %(asctime)s : %(message)s'
    logging.basicConfig(level=logging.INFO, format=formatter)
    _logger = logging.getLogger(__name__)

    f = open('config_mongo.json', 'r')
    _mongo_conf = json.load(f)
    hostName = str(_mongo_conf['host']) + ':' + str(_mongo_conf['port'])
    mongo = MongoModel(hostName, 'newsdb', 'article_col')
    _logger.info('%s', 'Finish Setup Mongodb')
    recom = RecommendModel()
    _logger.info('%s', 'Finish Setup PersonalRecommendModel from: ' + _ids_str)
    personal_items = recom.get_personal_items(_ids_str)
    if personal_items:
        personal_records = [mongo.find_one(filter={'_id': ObjectId(personal_id)}) for personal_id in personal_items]
        for personal_record in personal_records:
            try:
                personal_record['_id'] = str(personal_record['_id'])
                tm = datetime.strptime(personal_record['publishedAt'], '%Y-%m-%dT%H:%M:%S%z')
                personal_record['publishedAt'] = tm.strftime('%Y-%m-%d %H:%M')
                # _logger.info('%s', 'personal_record[\'titles\']:' + personal_record['titles'])
            except TypeError:
                pass
    else:
        personal_records = None
    _logger.info('%s', 'Finish Load personal_items')
    # Delete null record
    return jsonify({'data': [x for x in personal_records if x]})

if __name__ == "__main__":
    app.run(host='0.0.0.0')

    
