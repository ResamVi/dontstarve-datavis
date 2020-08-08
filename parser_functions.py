import re
import datetime
import logging

from pony   import orm
from bs4    import BeautifulSoup
from model  import db

# We only need the text of the <div> that makes up a server entry in the list
# Sequence is: name, platform, online, mode, season
def parse_server(server_html):
    soup = BeautifulSoup(server_html, 'html.parser').stripped_strings
    return next(soup), next(soup), next(soup), next(soup), next(soup)

player_pattern = re.compile(r'(\w+)')
char_pattern = re.compile(r'\[(.+)\]')

# Remove <a>-tags and leading/trailing whitespace
def parse_player(player_html):
    player_html = re.sub(r'<a.+\">', '', player_html)
    player_html = re.sub(r'</a>', '', player_html) 
    player_html = player_html.strip()

    # player name
    name = player_pattern.search(player_html)
    name = "<DIFFICULT TO PARSE>" if name is None else name.group(1)
    
    # character name
    character = char_pattern.search(player_html)
    character = "<TO BE DETERMINED>" if character is None else character.group(1)

    return name, character


@orm.db_session
def parse(server, players, cycle):
    # parse server
    server_html = server.get_attribute("innerHTML")
    name, platform, online, mode, season = parse_server(server_html)

    srv = db.Server(
        cycle=cycle,
        name=name,
        platform=platform,
        online=online,
        mode=mode,
        season=season,
        date=datetime.datetime.now()
    )
    logging.info("New Server: '%s'", srv.name)

    # parse players of server
    for player in players:
        player_html = player.get_attribute("innerHTML")
        name, character = parse_player(player_html)
        
        pl = db.Player(
            cycle=cycle,
            name=name,
            character=character,
            server=srv
        )
        logging.info("New Player: '%s'", pl.name)
    
    orm.commit()

