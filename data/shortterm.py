import os
import re
import json
import psycopg2
import datetime
import logging
import geoip2.database

from pony import orm

db = orm.Database()

# Convert platform number to name (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
platforms = lambda i : {1: 'Steam', 4: 'WeGame', 19: 'Console'}.get(i, str(i))

# GeoIP
reader = geoip2.database.Reader('./GeoLite2-Country.mmdb')

class Server(db.Entity):
    id              = orm.PrimaryKey(int, auto=True)
    name            = orm.Optional(str)
    country         = orm.Required(str)
    iso             = orm.Required(str)
    continent       = orm.Required(str)
    platform        = orm.Required(str)
    connected       = orm.Required(int)
    maxconnections  = orm.Required(int)
    elapsed         = orm.Required(int)
    mode            = orm.Optional(str)
    season          = orm.Optional(str)
    intent          = orm.Optional(str)
    mods            = orm.Optional(str)
    cycle           = orm.Required(int)
    date            = orm.Required(datetime.datetime)
    players         = orm.Set("Player")

def createServer(data, cycle):
    # Get origin of server via IP
    try:
        geoip       = reader.country(data["__addr"])
        country     = geoip.country.name
        continent   = geoip.continent.names['en']
        iso         = geoip.country.iso_code
    except:
        country     = "Antarctica"
        continent   = "Antarctica"
        iso         = "AQ"

    elapsed = re.search(r"(\d+)", data["data"])
    elapsed = elapsed.group() if elapsed is not None else -1

    srv = db.Server(
        name            = data["name"],
        country         = country,
        iso             = iso,
        continent       = continent,
        platform        = platforms(data["platform"]),
        connected       = data["connected"],
        maxconnections  = data["maxconnections"],
        elapsed         = elapsed,
        mode            = data["mode"],
        season          = data["season"],
        intent          = data["intent"],
        mods            = "vanilla" if data["mods"] else "modded",
        cycle           = cycle,
        date            = datetime.datetime.now()
    )
    logging.info("New Server: '%s'", srv.name)

    return srv

class Player(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Optional(str)
    character   = orm.Optional(str)
    country     = orm.Required(str)
    iso         = orm.Required(str)
    continent   = orm.Required(str)
    server      = orm.Required(Server)

def createPlayer(data, server, cycle):
    if data["players"] == "return {  }":
        return

    # List of players is not json ('return {...}' for whatever reason) 
    # so we fix that with a lot of duct tape
    players = data["players"] \
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
        return
    
    for player in players:
        pl = db.Player(
            cycle           = cycle,
            name            = player["name"],
            character       = player["prefab"],
            country         = server.country,
            iso             = server.iso,
            continent       = server.continent,
            server          = server
        )
        logging.info("\tNew Player: '%s'", pl.name)

# Create Views (we temporarily create a connection to execute raw sql)
def createViews():
    connection = psycopg2.connect(
        port="5432",
        user=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        database=os.getenv('DB_SHORT'),
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

def getLastUpdate():
    last_update = db.select("SELECT date FROM last_update")[0]
    duration = datetime.datetime.now() - last_update

    print(duration)
    return duration

def prepareSnapshot():
    server_count, player_count = db.select("SELECT server_count, player_count FROM count")[0]
    character_count = db.select("SELECT character, count FROM count_character WHERE character IN ('wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly')")
    country_count = db.select("SELECT country, count FROM count_player")

    characters = ['wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly']
    
    topfive_percentage = {}
    for character in characters:
        list = db.select("SELECT country, percent FROM percentage_character_by_country WHERE total_count > 30 AND character = $character ORDER BY percent DESC LIMIT 5")
        topfive_percentage[character] = list

    return {
        "player_count": player_count,
        "server_count": server_count,
        "character_count": dict(character_count),
        "country_count": dict(country_count),
        "topfive_percentage": topfive_percentage
    }

    