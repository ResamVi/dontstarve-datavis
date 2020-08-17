import datetime
import requests
import logging
import json
import time
import os

from pprint import pprint
from dotenv import load_dotenv
from geoip  import geolite2
from pony   import orm
from model  import db

# static variables
servers = [
    "https://lobby-us.kleientertainment.com/lobby/read",
    "https://lobby-eu.kleientertainment.com/lobby/read",
    "https://lobby-china.kleientertainment.com/lobby/read",
    "https://lobby-sing.kleientertainment.com/lobby/read"
]

# Convert platform number to name (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
platforms = lambda i : {0: 'None', 1: 'Steam', 2: 'PSN', 4: 'TGP', 8: 'WeGame / QQgame', 10: 'XBLIVE'}.get(i, 'Invalid')

# Logging
logging.basicConfig(format="[%(asctime)s] %(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.INFO)

# lazy solution to waiting on db docker container to finish loading
#logging.warning("Waiting ten seconds before starting") 
#time.sleep(10)

# Load .env file
load_dotenv()
#payload = '{"__token": "%s", "__gameId": "DST", "query": {"__rowId":"f8e4880052416ffffc93e7e48267bcdb"}}' % os.getenv("TOKEN")
payload = '{"__token": "%s", "__gameId": "DST", "query": {}}' % os.getenv("TOKEN")

r = requests.post(servers[2], data=payload)

servers = r.json()["GET"]

cycle = 1

for server in servers:
    
    #print(server["connected"])
    #print(server["name"])
    #print(server["players"])
    #print(server["platform"])
    #print(server["connected"])
    #print(server["maxconnections"])
    #print(server["mode"])
    #print(server["season"])
    #print(type(server["platform"]))

    origin = geolite2.lookup(server["__addr"])
    origin = origin.country if origin is not None else "None"

    srv = db.Server(
        name=server["connected"],
        origin=origin,
        platform=platforms(server["platform"]),
        connected=server["connected"],
        maxconnections=server["maxconnections"],
        mode=server["mode"],
        season=server["season"],
        cycle=cycle,
        date=datetime.datetime.now()
    )
    logging.info("New Server: '%s'", srv.name)
    
    break

#print(answer["GET"][0]["__rowId"])
#print(answer["GET"][0]["data"])
#print(answer["GET"][0]["connected"])
#print(answer["GET"][0]["name"])
#print(answer["GET"][0]["players"])

# TODO: Query-er speichert Ergebnisse routinemäßig in json-Dateien bzw. parst sie und in die Datenbank
# TODO: Web-Server: REST API

# Write to file
# f = open("output.txt", "a")
# f.write(json.dumps(r.json(), indent=4, sort_keys=True))
# f.close()