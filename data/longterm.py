import re
import json
import logging
import datetime
import psycopg2
import geoip2.database

from pony import orm

db = orm.Database()

# Convert platform number to name (see: https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data/?do=findComment&comment=1306033)
platforms = lambda i : {1: 'Steam', 4: 'WeGame', 19: 'Console'}.get(i, str(i))

# GeoIP
reader = geoip2.database.Reader('./GeoLite2-Country.mmdb')

class Server(db.Entity):
    name            = orm.Required(str)
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
    orm.PrimaryKey(name, country)
    
def createServer(data, cycle):
    # Get origin of server via IP
    geoip       = reader.country(data["__addr"])
    country     = geoip.country.name
    continent   = geoip.continent.names['en']
    iso         = geoip.country.iso_code

    elapsed = re.search(r"(\d+)", data["data"])
    elapsed = elapsed.group() if elapsed is not None else -1

    if db.Server.exists(name = data["name"], country = country):
        return db.Server.get(name = data["name"], country = country)

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
    cycle       = orm.Required(int)
    name        = orm.Required(str)
    character   = orm.Optional(str)
    country     = orm.Required(str)
    iso         = orm.Required(str)
    continent   = orm.Required(str)
    duration    = orm.Required(int)
    server      = orm.Required(Server)
    orm.PrimaryKey(name, server)

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
        
        # Empty name edge case
        if player["name"].strip() == "":
            player["name"] = "unnamed"

        if db.Player.exists(name = player["name"], server = server):
            pl = db.Player.get(name = player["name"], server = server)
            pl.duration += 5
            db.commit()
            logging.info("\tUpdating Player: '%s' to %d", pl.name, pl.duration)
        else:
            pl = db.Player(
                cycle           = cycle,
                name            = player["name"],
                character       = player["prefab"],
                country         = server.country,
                iso             = server.iso,
                continent       = server.continent,
                duration        = 0,
                server          = server
            )
            logging.info("\tNew Player: '%s'", pl.name)

class SeriesPlayerCount(db.Entity):
    date            = orm.Required(datetime.datetime)
    count           = orm.Required(int)

def createSnapshot():
    db.select('name FROM Person WHERE age > $x')