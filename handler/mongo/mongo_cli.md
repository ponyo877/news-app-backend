## elasticsearch Commands
kubectl exec -it quickstart-es-default-0 sh
kubectl get secret quickstart-es-elastic-user -o go-template='{{.data.elastic | base64decode}}'
curl -u "elastic:1m4OuZj3Ap3X5HSc6G0w8B40" -k "https://quickstart-es-http:9200/_search?pretty" -H 'Content-Type: application/json' -d'{"query": {"match_all": {}}, "sort" : [{"publishedAt":"desc"}]}'

curl -u "elastic:1m4OuZj3Ap3X5HSc6G0w8B40" -k "https://quickstart-es-http:9200/_bulk?pretty" -H 'Content-Type: application/json' -d'{"create":{ "_index" : "test_es2" , "_id" : "602fba5ba9e5c03c204c6b22"}}
{"id":"602fba5ba9e5c03c204c6b22","image":"https://img.gitouhon-juku-k8s2.ga/2091160353.jpeg","publishedAt":"2021-02-19T22:15:19+09:00","sitetitle":"ニュー速クオリティ","titles":"【画像】Uber eats(ウーバーイーツ)の闇が深すぎる件・・・","url":"http://news4vip.livedoor.biz/archives/52389213.html"}'

# mongo operation
kubectl exec -it mongo-0 sh
mongo --host mongo.default --port 27017

rs.initiate({_id: "rs0", members: [{_id: 0, host: "mongo.default:27017"}]})
use newsdb
db.createCollection('article_col');
db.article_col.insert({title: 'Yahoo', URL: 'https://www.yahoo.co.jp/', image: 'tbd', updateDate: new Date(), click: 0, siteID: 1});

db.createCollection('site_col');
db.site_col.insert({siteID: 1, sitetitle: '痛いニュース',           rssURL: 'http://blog.livedoor.jp/dqnplus/index.rdf',   latestDate: '2020-01-01 00:00:00'});
db.site_col.insert({siteID: 4, sitetitle: 'ハムスター速報',         rssURL: 'http://hamusoku.com/index.rdf',               latestDate: '2020-01-01 00:00:00'});
db.site_col.insert({siteID: 5, sitetitle: '暇人＼^o^／速報',        rssURL: 'http://himasoku.com/index.rdf',               latestDate: '2020-01-01 00:00:00'});
db.site_col.insert({siteID: 6, sitetitle: 'VIPPERな俺',             rssURL: 'http://blog.livedoor.jp/news23vip/index.rdf', latestDate: '2020-01-01 00:00:00'});
db.site_col.insert({siteID: 3, sitetitle: 'ニュー速クオリティ',     rssURL: 'http://news4vip.livedoor.biz/index.rdf',      latestDate: '2020-01-01 00:00:00'});

db.site_col.deleteMany( { siteID: {$lt: 10} } )