import datetime
import requests
import logging
import json
import time
import os
import re
import geoip2.database

from pprint import pprint
from dotenv import load_dotenv
from pony   import orm
from model  import db

import psycopg2

# static variables
endpoints = [
    "https://lobby-us.kleientertainment.com/lobby/read",
    "https://lobby-eu.kleientertainment.com/lobby/read",
    "https://lobby-china.kleientertainment.com/lobby/read",
    "https://lobby-sing.kleientertainment.com/lobby/read"
]

# Convert platform number to name (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
platforms = lambda i : {1: 'Steam', 4: 'WeGame', 19: 'Console'}.get(i, str(i))

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
    password=os.getenv('POSTGRES_PASSWORD'),
    database=os.getenv('POSTGRES_DB'),
    host=os.getenv('HOST')
)

# Create Tables
db.generate_mapping(create_tables=True)

# Create Views (we temporarily create a connection to execute raw sql)
connection = psycopg2.connect(
    port="5432",
    user=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    database=os.getenv('POSTGRES_DB'),
    host=os.getenv('HOST'))

cursor = connection.cursor()
cursor.execute(open("views.sql", "r").read())
connection.commit()

cursor.close()
connection.close()

# GeoIP
reader = geoip2.database.Reader('./GeoLite2-Country.mmdb')

@orm.db_session
def main(endpoint, cycle):

    payload = '{"__token": "%s", "__gameId": "DST", "query": {}}' % os.getenv("TOKEN")
    r = requests.post(endpoint, data=payload)
    servers = r.json()["GET"]

    for server in servers:

        # Get origin of server via IP
        try:
            origin = reader.country(server["__addr"]).country.name
        except:
            origin = "None"

        elapsed = re.search("(\d+)", server["data"])
        elapsed = elapsed.group() if elapsed is not None else -1

        srv = db.Server(
            name=server["name"],
            origin=origin,
            platform=platforms(server["platform"]),
            connected=server["connected"],
            maxconnections=server["maxconnections"],
            elapsed=elapsed,
            mode=server["mode"],
            season=server["season"],
            intent=server["intent"],
            mods="vanilla" if server["mods"] else "modded",
            cycle=cycle,
            date=datetime.datetime.now()
        )
        logging.info("New Server: '%s'", srv.name)

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
        
        # Some player names give me a headache
        try:
            players = json.loads(players)
        except:
            continue
        
        for player in players:
            pl = db.Player(
                cycle=cycle,
                name=player["name"],
                character=player["prefab"],
                origin=srv.origin,
                server=srv
            )
            logging.info("New Player: '%s'", pl.name)
    
    # Add activity
    activ = db.Activity(
        date=datetime.datetime.now(),
        countbyorigin=db.select("SELECT * FROM count_player")
    )
    logging.info("New Activity created")
    

cycle = 1
while True:
    for endpoint in endpoints:
        logging.warning("Endpoint: " + endpoint)
        logging.warning("Cycle: " + str(cycle))
        
        main(endpoint, cycle)
    
    logging.warning("Finished Cycle " + str(cycle))
    cycle += 1
    time.sleep(60 * 15) # Update every 15 minutes
    # TODO: Clear tables

# TODO: multiple series chart of activity
# TODO: Send better HTTP Codes isntead of panic
# TODO: Allow client to set LIMIT
# TODO: Collapse some handlers, repeated code...
# TODO: date of last fetch 