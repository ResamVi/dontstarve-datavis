import time
import logging
import sys

from os import environ as env
from distutils import util

from selenium                                       import webdriver
from selenium.webdriver.common.desired_capabilities import DesiredCapabilities
from selenium.webdriver.common.keys                 import Keys
from selenium.webdriver.common.by                   import By
from selenium.webdriver.support.ui                  import WebDriverWait
from selenium.webdriver.support                     import expected_conditions as expected

from bs4                        import BeautifulSoup
from pony                       import orm
from webdriver_manager.chrome   import ChromeDriverManager

from parser_functions import parse
from model import db

# Logging init
logging.basicConfig(format="[%(asctime)s] %(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.INFO)

# lazy solution to waiting on db container 
logging.warning("Waiting ten seconds before starting") 
time.sleep(10)

# Database init
db.bind(
    provider='postgres',
    user=env['POSTGRES_USER'],
    password=env['POSTGRES_PASS'],
    host='db',
    database=env['POSTGRES_DB']
)

db.generate_mapping(create_tables=True)

# Selenium init
options = webdriver.ChromeOptions()
options.headless = True
options.add_argument("--start-maximized")
options.add_argument('--no-sandbox')       

driver = webdriver.Remote("http://selenium:4444/wd/hub", DesiredCapabilities.CHROME)
# driver = webdriver.Chrome(ChromeDriverManager().install(), options=options)
logging.warning("Remote Selenium connection established")

# Loading times on the website are yuge
driver.get("https://dstserverlist.appspot.com/")

# Remove cookie notification
driver.find_element_by_xpath("//a[@aria-label='dismiss cookie message']").click()

# Determine cycle (whether or not previous queries exist)

def start_scraping(cycle):
    
    # Iterate all pages
    page_index = 1
    
    while True:
        pageX = 50
        start_time = time.time()
        
        WebDriverWait(driver, 10).until(expected.element_to_be_clickable((By.CLASS_NAME, "serverlist-entry")))
        servers = driver.find_elements_by_class_name("serverlist-entry")

        # Click on every server to load the list of players into the modal
        for server in servers:
            server.click()

            # Wait for Player List to load
            # If "Error: Server does not exist" will be displayed a timeout occurs
            try:
                WebDriverWait(driver, 5).until(expected.presence_of_element_located((By.XPATH, "//div[@id='players']//div[@class='col s12 m6 l3']")))
                players = driver.find_elements_by_xpath("//div[@id='players']//div[@class='col s12 m6 l3']")
                
                # start parsing
                parse(server, players, cycle)
            except:
                pass
            
            # navigate back to main menu
            webdriver.ActionChains(driver).send_keys(Keys.ESCAPE).perform()

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

# Start and continue forever
with orm.db_session:
    logging.warning("Start scraping")

    cycle = 1 + db.select('MAX(cycle) FROM player')[0] # Continue cycle number if previous exist

    while True:
        start_scraping(cycle)
        cycle += 1
        logging.warning("Cycle finished. Now starting Cycle " + str(cycle))

# TODO: Show progress (3/165) of tables