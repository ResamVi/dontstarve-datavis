import re

from bs4 import BeautifulSoup

player_pattern = re.compile(r'(\w+)')
char_pattern = re.compile(r'\[(.+)\]')

def parse_server(server_html):
    server_html = BeautifulSoup(server_html, 'html.parser').stripped_strings

    print("Server Name: " + next(server_html))
    print("Platform: " + next(server_html))
    print("Players: " + next(server_html))
    print("Mode: " + next(server_html))
    print("Season: " + next(server_html))

def parse_player(player_html):
    # strip <a>-tags if it exists and trailing/leading whitespace
    player_html = re.sub(r'<a.+\">', '', player_html)
    player_html = re.sub(r'</a>', '', player_html) 
    player_html = player_html.strip()

    # player name
    player = player_pattern.search(player_html)
    player = "<DIFFICULT TO PARSE>" if player is None else player.group(1)
    
    # character name
    character = char_pattern.search(player_html)
    character = "<TO BE DETERMINED>" if character is None else character.group(1)
    
    print("Player: " + player)
    print("Character: " + character)
    print("")