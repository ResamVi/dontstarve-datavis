import time
import re

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.firefox.options import Options

from webdriver_manager.firefox import GeckoDriverManager

# Headless mode to not open a browser and run this on a vps
options = Options()
options.headless = False

# Use firefox because chrome throws a bluetooth error lol
driver = webdriver.Firefox(options=options, executable_path=GeckoDriverManager().install())

# Loading times on the website are yuge
driver.implicitly_wait(20)
driver.get("https://dstserverlist.appspot.com/")

player_pattern = re.compile(r'(\w+)')
char_pattern = re.compile(r'\[(.+)\]')

# Iterate all pages
while True:
    servers = driver.find_elements_by_class_name("serverlist-entry") # 50 servers on a page

    # Click on every server to load the list of players into the modal
    for server in servers:
        server.click()

        # Retrieve players from modal
        players = driver.find_elements_by_xpath("//div[@id='players']//div[@class='col s12 m6 l3']")

        for player in players:
            player_html = player.get_attribute("innerHTML")

            # strip <a>-tags if it exists and trailing/leading whitespace
            player_html = re.sub(r'<a.+\">', '', player_html)
            player_html = re.sub(r'</a>', '', player_html) 
            player_html = player_html.strip()

            # player name
            player = player_pattern.search(player_html).group(1)
            print("Player: " + player)
            
            # character name
            character = char_pattern.search(player_html)
            character = "TBD" if character is None else character.group(1)
            
            print("Character: " + character)
            print("")

        # navigate back to main menu
        webdriver.ActionChains(driver).send_keys(Keys.ESCAPE).perform()

driver.close()