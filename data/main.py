import geoip2.database
import psycopg2
import datetime
import requests
import logging
import json
import time
import os
import re

import shortterm
import longterm

from pprint import pprint
from dotenv import load_dotenv
from pony   import orm

# Klei endpoints
endpoints = [
    "https://lobby-us.kleientertainment.com/lobby/read",
    "https://lobby-eu.kleientertainment.com/lobby/read",
    "https://lobby-china.kleientertainment.com/lobby/read",
    "https://lobby-sing.kleientertainment.com/lobby/read"
]

# Logging
logging.basicConfig(format="[%(asctime)s] %(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.WARNING)

# Load .env file
load_dotenv()

# Silly hack to wait for docker container to initialize
logging.warning("Waiting 5s")
time.sleep(5)

# Database init
shortterm.db.bind(
    provider='postgres',
    user=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    database=os.getenv('DB_SHORT'),
    host=os.getenv('DB_HOST')
)

longterm.db.bind(
    provider='postgres',
    user=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    database=os.getenv('DB_LONG'),
    host=os.getenv('DB_HOST')
)

longterm.db.generate_mapping(create_tables=True)
shortterm.db.generate_mapping(create_tables=True)
shortterm.createViews()

cycle = 1
while True:
    logging.warning("Starting Cycle " + str(cycle))
    
    with orm.db_session:
        for endpoint in endpoints:    
            
            logging.warning("Endpoint: " + endpoint)
            
            payload = '{"__token": "%s", "__gameId": "DST", "query": {}}' % os.getenv("TOKEN")
            
            with requests.post(endpoint, data=payload) as answer:
                servers = answer.json()["GET"]

            # Insert short-term data
            for server in servers:
                srv = shortterm.createServer(server, cycle)
                shortterm.createPlayer(server, srv, cycle)

            # Insert long-term data
            for server in servers:
                srv = longterm.createServer(server, cycle)
                longterm.createPlayer(server, srv, cycle, shortterm.getLastUpdate())

        # Insert series data
        snapshot = shortterm.prepareSnapshot()
        longterm.createSnapshot(snapshot)
    
    logging.warning("Finished Cycle " + str(cycle))
    
    cycle += 1
    time.sleep(3 * 60) # Update every 5 minutes
    
    shortterm.clearTables()
    shortterm.createViews()


# Sample Query
# 
# {
#     "GET": [
#         {
#             "Users": null,
#             "__addr": "------",
#             "__lastPing": 1597668970,
#             "__rowId": "f8e4880052416ffffc93e7e48267bcdb",
#             "allownewplayers": true,
#             "clanonly": false,
#             "clienthosted": false,
#             "clientmodsoff": false,
#             "connected": 4,
#             "data": "return { day=105, dayselapsedinseason=14, daysleftinseason=1 }",
#             "dedicated": true,
#             "desc": "\u4e00\u8d77\u6765\u8bed\u97f3\u5f00\u9ed1\u5427!QQ\u7fa4:6430 69994 \u5bc6\u7801\u7fa4\u5185\u63d0\u4f9b\u54e6\uff01\u7d20\u8d28\u6e38\u620f\uff01",
#             "event": false,
#             "fo": false,
#             "guid": "2067472753315897219",
#             "host": "-----",
#             "intent": "social",
#             "lanonly": false,
#             "maxconnections": 8,
#             "mode": "endless",
#             "mods": true,
#             "mods_info": [
#                 "workshop-2199027653598523456",
#                 "\u989d\u5916\u88c5\u5907\u680fExtra Equip Slots (Updated)",
#                 "1.8.0.2",
#                 "1.8.0.2",
#                 true,
#                 "workshop-2199027653598518755",
#                 "DST\u5783\u573e\u81ea\u52a8\u6e05\u7406",
#                 "1.6",
#                 "1.6",
#                 true,
#                 "workshop-2199027653598517004",
#                 "\u751f\u7269\u52a0\u5f3a(Monster strengthened)",
#                 "1.31",
#                 "1.31",
#                 true,
#                 "workshop-2199027653598523217",
#                 "\u5c0f\u767d\u670d\u52a1\u7aef",
#                 "4.6",
#                 "4.6",
#                 true,
#                 "workshop-2199027653598523720",
#                 "The Forge Item Pack",
#                 "1.0.7",
#                 "1.0.7",
#                 true
#             ],
#             "name": "!\u840c\u795e\u670d\u52a1\u5668! \u751f\u7269\u52a0\u5f3a \u534a\u7eaf\u51c0\u6863 \u8bed\u97f3\u623f Q\u7fa4:6430 69994 \u4e00\u8d77\u73a9\u5416\uff01",
#             "nat": 7,
#             "password": true,
#             "platform": 4,
#             "players": "return {\n  {\n    colour=\"8B668B\",\n    eventlevel=0,\n    name=\"\u6211\u53c8\u4e0d\u662f\u9648\u5955\u8fc5\",\n    netid=\"R:76561197982261011\",\n    prefab=\"wilson\" \n  },\n  {\n    colour=\"CD853F\",\n    eventlevel=0,\n    name=\"\u517b\u4e00\u53ea\u5948\u5948\",\n    netid=\"R:76561197982331087\",\n    prefab=\"wathgrithr\" \n  },\n  {\n    colour=\"CDAA7D\",\n    eventlevel=0,\n    name=\"\u840c\u795e\u54df\",\n    netid=\"R:76561197973675353\",\n    prefab=\"wathgrithr\" \n  },\n  {\n    colour=\"FFA54F\",\n    eventlevel=0,\n    name=\"\u6728\u8111\u58f3\",\n    netid=\"R:76561197972551185\",\n    prefab=\"wathgrithr\" \n  } \n}",
#             "port": 10999,
#             "pvp": false,
#             "season": "winter",
#             "secondaries": {
#                 "2308447579": {
#                     "__addr": "---",
#                     "__lastPing": 1597668782,
#                     "id": "2308447579",
#                     "port": 10998,
#                     "steamid": "R:90071997884764293"
#                 }
#             },
#             "session": "24CE17BEE4098F2B",
#             "slaves": {
#                 "2308447579": {
#                     "__addr": "---",
#                     "__lastPing": 1597668782,
#                     "id": "2308447579",
#                     "port": 10998,
#                     "steamid": "R:90071997884764293"
#                 }
#             },
#             "steamid": "R:90071997884764299",
#             "steamroom": "0",
#             "tags": "endless,\u6295\u7968,\u6d1e\u7a74,extra equip slot eqs,clean,\u6e05\u7406,the forge,dst,items,pack,event,forge,battlemaster pugna",
#             "tick": 15,
#             "v": 422664,
#             "valvecloudserver": false,
#             "valvepopid": "",
#             "valveroutinginfo": "",
#             "worldgen": "return {\n  {\n    desc=\"\u6807\u51c6\u300a\u9965\u8352\u300b\u4f53\u9a8c\u3002\",\n    hideminimap=false,\n    id=\"SURVIVAL_TOGETHER\",\n    location=\"forest\",\n    max_playlist_position=999,\n    min_playlist_position=0,\n    name=\"\u9ed8\u8ba4\",\n    numrandom_set_pieces=4,\n    override_level_string=false,\n    overrides={ frograin=\"never\" },\n    random_set_pieces={\n      \"Sculptures_2\",\n      \"Sculptures_3\",\n      \"Sculptures_4\",\n      \"Sculptures_5\",\n      \"Chessy_1\",\n      \"Chessy_2\",\n      \"Chessy_3\",\n      \"Chessy_4\",\n      \"Chessy_5\",\n      \"Chessy_6\",\n      \"Maxwell1\",\n      \"Maxwell2\",\n      \"Maxwell3\",\n      \"Maxwell4\",\n      \"Maxwell6\",\n      \"Maxwell7\",\n      \"Warzone_1\",\n      \"Warzone_2\",\n      \"Warzone_3\" \n    },\n    required_prefabs={ \"multiplayer_portal\" },\n    required_setpieces={ \"Sculptures_1\", \"Maxwell5\" },\n    substitutes={  },\n    version=4 \n  },\n  {\n    background_node_range={ 0, 1 },\n    desc=\"\u63a2\u67e5\u6d1e\u7a74\u2026\u2026 \u4e00\u8d77\uff01\",\n    hideminimap=false,\n    id=\"DST_CAVE\",\n    location=\"cave\",\n    max_playlist_position=999,\n    min_playlist_position=0,\n    name=\"\u6d1e\u7a74\",\n    numrandom_set_pieces=0,\n    override_level_string=false,\n    overrides={ start_location=\"caves\", task_set=\"cave_default\" },\n    required_prefabs={ \"multiplayer_portal\" },\n    substitutes={  },\n    version=4 \n  } \n}"
#         }
#     ]
# }