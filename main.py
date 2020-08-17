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
endpoints = [
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

# Database init
db.bind(
    provider='postgres',
    user=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASS'),
    database=os.getenv('POSTGRES_DB'),
    host=os.getenv('HOST')
)

db.generate_mapping(create_tables=True)

@orm.db_session
def main(endpoint, cycle):

    payload = '{"__token": "%s", "__gameId": "DST", "query": {}}' % os.getenv("TOKEN")
    r = requests.post(endpoint, data=payload)
    servers = r.json()["GET"]

    f = open("output.txt", "a")
    f.write(json.dumps(r.json(), indent=4, sort_keys=True))
    f.close()

    for server in servers:

        # Get origin of server via IP
        origin = geolite2.lookup(server["__addr"])
        origin = origin.country if origin is not None else "None"

        srv = db.Server(
            name=server["name"],
            origin=origin,
            platform=platforms(server["platform"]),
            connected=server["connected"],
            maxconnections=server["maxconnections"],
            mode=server["mode"],
            season=server["season"],
            intent=server["intent"],
            mods=server["mods"],
            cycle=cycle,
            date=datetime.datetime.now()
        )
        logging.info("New Server: '%s'", srv.name)
        
        #print("---")
        #print("BEFORE")
        #print(server["players"])

        # Player list empty
        if server["players"] == "return {  }":
            continue

        # List of players is not json ('return {...}' for whatever reason) 
        # so we fix that with a lot of duct tape
        players = server["players"] \
            .replace("return {", "[") \
            .replace('colour=', '"colour":') \
            .replace('prefab=', '"prefab":') \
            .replace('eventlevel=', '"eventlevel":') \
            .replace('name=', '"name":') \
            .replace('netid=', '"netid":') \
            .replace("\n}", "]") \
            .replace('["', '"') \
            .replace('"]=', '":')

        #print("AFTER")
        #print(players)
        #print("---")
        players = json.loads(players)
        
        for player in players:
            pl = db.Player(
                cycle=cycle,
                name=player["name"],
                character=player["prefab"],
                server=srv
            )
            logging.info("New Player: '%s'", pl.name)
    

cycle = 1
for endpoint in endpoints:
    logging.warning("Endpoint: " + endpoint)
    main(endpoint, cycle)

# TODO: Web-Server: REST API
# TODO: Check how many days in
# TODO: Players inherit server's origin
# TODO: time_characters: Date<wendy, wigfrid, wilson, ...>
# TODO: time_origins: Date<China, USA>

# TODO: Chart of previous 24h time_characters
# TODO: multiple series chart of previous 24h time_origins
# TODO: bar Chart of Steam/TGP and
# TODO: geo chart of player