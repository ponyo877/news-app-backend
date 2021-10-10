import json

json_path = 'users2.json'
json_file = open(json_path, 'r')
json_dict = json.load(json_file)
for dup_item in json_dict:
    for delete_target_item in dup_item['items'][1:]:
        print(delete_target_item)

