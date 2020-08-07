import time
import datetime
import logging

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.firefox.options import Options

from bs4 import BeautifulSoup

from webdriver_manager.chrome import ChromeDriverManager

from parser_functions import parse_server, parse_player

logging.basicConfig(format="[%(asctime)s] %(message)s")
logging.getLogger().setLevel(logging.INFO)

"""
from pony.orm import *

db = Database()

class Server(db.Entity):
    id          = PrimaryKey(int, auto=True)
    group       = Required(int)
    name        = Required(str)
    lp          = Required(int)
    wins        = Required(int)
    losses      = Required(int)
    date        = Required(datetime.datetime)
    summonerId  = Required(str)

class Player(db.Entity):


db.bind(provider='postgres', user='root', password='password', host='localhost', database='mydatabase')
db.generate_mapping(create_tables=True)
"""
# ----------------------------------------------------------------

# Headless mode to not open a browser and run this on a vps
options = Options()
options.headless = False

# Use firefox because chrome throws a bluetooth error lol
driver = webdriver.Chrome(ChromeDriverManager().install())

# Loading times on the website are yuge
driver.implicitly_wait(5)
driver.get("https://dstserverlist.appspot.com/")

# Remove cookie notification
driver.find_element_by_xpath("//a[@aria-label='dismiss cookie message']").click()

driver.execute_script("window.scrollTo(0,3000);")
driver.find_elements_by_class_name("page")[-1].click()

def start_scraping():
    
    # Iterate all pages
    page_index = 1
    
    while True:
        pageX = 50
        start_time = time.time()
        
        # Click on every server to load the list of players into the modal
        servers = driver.find_elements_by_class_name("serverlist-entry")
        for server in servers:
            server_html = server.get_attribute("innerHTML")
            parse_server(server_html)
            
            try:
                server.click()
            except:
                continue

            # Retrieve players from modal
            players = driver.find_elements_by_xpath("//div[@id='players']//div[@class='col s12 m6 l3']")
            for player in players:
                player_html = player.get_attribute("innerHTML")
                parse_player(player_html)

            # navigate back to main menu
            webdriver.ActionChains(driver).send_keys(Keys.ESCAPE).pause(2).perform()

            # Scroll down a bit
            driver.execute_script("window.scrollTo(0," + str(pageX) + ");")
            pageX += 50

        # Get button to next page
        for button in driver.find_elements_by_class_name("page"):
            if button.text == str(page_index + 1):
                button.click()
                break
        else:
            driver.close() # No more pages found
            logging.info("No more pages found: %d pages parsed", page_index)
            return
        
        driver.execute_script("window.scrollTo(0,0);")
        
        # Track and log elapsed time
        elapsed_time = time.time() - start_time
        logging.warning("Parsed page %d in %s", page_index, time.strftime("%H:%M:%S", time.gmtime(elapsed_time)))
        
        page_index += 1


start_scraping()
