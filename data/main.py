import geoip2.database
import psycopg2
import datetime
import requests
import logging
import json
import time
import os
import re

from pprint import pprint
from dotenv import load_dotenv
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
platforms = lambda i : {1: 'Steam', 4: 'WeGame', 19: 'Console'}.get(i, str(i))

@orm.db_session
def getData(endpoint, cycle):

    try:
        payload = '{"__token": "%s", "__gameId": "DST", "query": {}}' % os.getenv("TOKEN")
        r = requests.post(endpoint, data=payload)
        servers = r.json()["GET"]
    except:
        return

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

# Create Views (we temporarily create a connection to execute raw sql)
def createViews():
    connection = psycopg2.connect(
        port="5432",
        user=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        database=os.getenv('POSTGRES_DB'),
        host=os.getenv('DB_HOST')
    )

    cursor = connection.cursor()
    cursor.execute(open("views.sql", "r").read())
    connection.commit()

    cursor.close()
    connection.close()

def clearTables():
    db.drop_all_tables(with_all_data=True)
    db.create_tables()

# -----
# Logging
logging.basicConfig(format="[%(asctime)s] %(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.WARNING)

# Load .env file
load_dotenv()

# Silly hack to wait for docker container to initialize
logging.warning("Waiting 10s")
time.sleep(5)

# Database init
db.bind(
    provider='postgres',
    user=os.getenv('POSTGRES_USER'),
    password=os.getenv('POSTGRES_PASSWORD'),
    database=os.getenv('POSTGRES_DB'),
    host=os.getenv('DB_HOST')
)

# Create Tables and Views
db.generate_mapping(create_tables=True)
createViews()

# GeoIP
reader = geoip2.database.Reader('./GeoLite2-Country.mmdb')

cycle = 1
while True:
    logging.warning("Starting Cycle " + str(cycle))
    for endpoint in endpoints:
        logging.warning("Endpoint: " + endpoint)
        
        getData(endpoint, cycle)
    
    logging.warning("Finished Cycle " + str(cycle))
    cycle += 1
    time.sleep(5 * 60) # Update every 5 minutes
    
    clearTables()
    createViews()

# TODO: multiple series chart of activity
