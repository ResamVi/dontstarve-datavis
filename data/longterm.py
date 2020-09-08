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
    try:
        geoip       = reader.country(data["__addr"])
        continent   = geoip.continent.names['en']
        country     = geoip.country.name if geoip.country.name is not None else "Antarctica"
        iso         = geoip.country.iso_code if geoip.country.iso_code is not None else "AQ"
    except:
        country     = "Antarctica"
        continent   = "Antarctica"
        iso         = "AQ"

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
    duration    = orm.Required(datetime .timedelta)
    server      = orm.Required(Server)
    orm.PrimaryKey(name, server)

def createPlayer(data, server, cycle, interval):
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
            pl.duration += interval
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
                duration        = datetime.timedelta(),
                server          = server
            )
            logging.info("\tNew Player: '%s'", pl.name)

class Series_Count(db.Entity):
    date            = orm.Required(datetime.datetime)
    player_count    = orm.Required(int)
    server_count    = orm.Required(int)

class Series_Character_Count(db.Entity):
    date            = orm.Required(datetime.datetime)
    wilson          = orm.Required(int)
    willow          = orm.Required(int)
    wolfgang        = orm.Required(int)
    wendy           = orm.Required(int)
    wx78            = orm.Required(int)
    wickerbottom    = orm.Required(int)
    woodie          = orm.Required(int)
    wes             = orm.Required(int)
    waxwell         = orm.Required(int)
    wathgrithr      = orm.Required(int)
    webber          = orm.Required(int)
    warly           = orm.Required(int)
    wormwood        = orm.Required(int)
    winona          = orm.Required(int)
    wortox          = orm.Required(int)
    wurt            = orm.Required(int)
    walter          = orm.Required(int)

class Series_Player_Count(db.Entity):
    date            = orm.Required(datetime.datetime)
    countries       = orm.Required(orm.Json)

class Series_Character_Ranking(db.Entity):
    date            = orm.Required(datetime.datetime)
    character       = orm.Required(str)
    first           = orm.Required(str)
    first_percent   = orm.Required(float)
    second          = orm.Required(str)
    second_percent  = orm.Required(float)
    third           = orm.Required(str)
    third_percent   = orm.Required(float)
    fourth          = orm.Required(str)
    fourth_percent  = orm.Required(float)
    fifth           = orm.Required(str)
    fifth_percent   = orm.Required(float)

class Series_Continent(db.Entity):
    date            = orm.Required(datetime.datetime)
    asia            = orm.Required(int)
    europe          = orm.Required(int)
    north_america   = orm.Required(int)
    south_america   = orm.Required(int)
    africa          = orm.Required(int)
    oceania         = orm.Required(int)
    

def createSnapshot(snapshot):
    
    db.Series_Count(
        date            = datetime.datetime.now(),
        player_count    = snapshot["player_count"],
        server_count    = snapshot["server_count"]
    )
    
    db.Series_Character_Count(
        date            = datetime.datetime.now(),
        wilson          = snapshot["character_count"]["wilson"],
        willow          = snapshot["character_count"]["willow"],
        wolfgang        = snapshot["character_count"]["wolfgang"],
        wendy           = snapshot["character_count"]["wendy"],
        wx78            = snapshot["character_count"]["wx78"],
        wickerbottom    = snapshot["character_count"]["wickerbottom"],
        woodie          = snapshot["character_count"]["woodie"],
        wes             = snapshot["character_count"]["wes"],
        waxwell         = snapshot["character_count"]["waxwell"],
        wathgrithr      = snapshot["character_count"]["wathgrithr"],
        webber          = snapshot["character_count"]["webber"],
        warly           = snapshot["character_count"]["warly"],
        wormwood        = snapshot["character_count"]["wormwood"],
        winona          = snapshot["character_count"]["winona"],
        wortox          = snapshot["character_count"]["wortox"],
        wurt            = snapshot["character_count"]["wurt"],
        walter          = snapshot["character_count"]["walter"]
    )        
    
    db.Series_Player_Count(date = datetime.datetime.now(), countries = snapshot["country_count"])

    db.Series_Continent(date = datetime.datetime.now(),
        asia            = snapshot["continent_count"]["Asia"],
        europe          = snapshot["continent_count"]["Europe"],
        north_america   = snapshot["continent_count"]["North America"],
        south_america   = snapshot["continent_count"]["South America"],
        africa          = snapshot["continent_count"]["Africa"], 
        oceania         = snapshot["continent_count"]["Oceania"]
    )

    characters = ['wendy', 'wathgrithr', 'wilson', 'woodie', 'wolfgang', 'wickerbottom', 'wx78', 'walter', 'webber', 'winona', 'waxwell', 'wortox', 'wormwood', 'wurt', 'wes', 'willow', 'warly']
    for character in characters:
        
        try:
            preference = db.Series_Character_Ranking(
                date            = datetime.datetime.now(),
                character       = character,
                first           = snapshot["topfive_percentage"][character][0][0],
                first_percent   = snapshot["topfive_percentage"][character][0][1],
                second          = snapshot["topfive_percentage"][character][1][0],
                second_percent  = snapshot["topfive_percentage"][character][1][1],
                third           = snapshot["topfive_percentage"][character][2][0],
                third_percent   = snapshot["topfive_percentage"][character][2][1],
                fourth          = snapshot["topfive_percentage"][character][3][0],
                fourth_percent  = snapshot["topfive_percentage"][character][3][1],
                fifth           = snapshot["topfive_percentage"][character][4][0],
                fifth_percent   = snapshot["topfive_percentage"][character][4][1]
            )
        except:
            print(character)
            print(snapshot["topfive_percentage"][character])
            
