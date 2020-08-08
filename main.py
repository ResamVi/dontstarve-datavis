import time
import logging

from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.firefox.options import Options

from bs4 import BeautifulSoup

from webdriver_manager.chrome import ChromeDriverManager

from parser_functions import parse
from model import db

# Logging init
logging.basicConfig(format="[%(asctime)s]%(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.INFO)

# Database init
db.bind(provider='postgres', user='root', password='password', host='localhost', database='mydatabase')
db.generate_mapping(create_tables=True)

# Selenium init
options = webdriver.ChromeOptions()
options.headless = False
options.add_argument("--start-maximized")
driver = webdriver.Chrome(ChromeDriverManager().install(), options=options)

# Loading times on the website are yuge
driver.implicitly_wait(5)
driver.get("https://dstserverlist.appspot.com/")

# Remove cookie notification
driver.find_element_by_xpath("//a[@aria-label='dismiss cookie message']").click()

# Determine cycle (whether or not previous queries exist)

def start_scraping():
    
    # Iterate all pages
    page_index = 1
    
    while True:
        pageX = 50
        start_time = time.time()
        
        # Click on every server to load the list of players into the modal
        servers = driver.find_elements_by_class_name("serverlist-entry")
        for server in servers:
            try:
                server.click()
            except:
                continue

            players = driver.find_elements_by_xpath("//div[@id='players']//div[@class='col s12 m6 l3']")
            
            # start parsing
            parse(server, players, 1)

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