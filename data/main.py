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
from model  import db, createViews, clearTables
from query  import query


# static variables
endpoints = [
    "https://lobby-us.kleientertainment.com/lobby/read",
    "https://lobby-eu.kleientertainment.com/lobby/read",
    "https://lobby-china.kleientertainment.com/lobby/read",
    "https://lobby-sing.kleientertainment.com/lobby/read"
]

# Logging
logging.basicConfig(format="[%(asctime)s] %(levelname)s — %(message)s")
logging.getLogger().setLevel(logging.WARNING)

# Load .env file
load_dotenv()

# Silly hack to wait for docker container to initialize
logging.warning("Waiting 5s")
# time.sleep(5)

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

cycle = 1
while True:
    logging.warning("Starting Cycle " + str(cycle))
    
    for endpoint in endpoints:
        try:
            logging.warning("Endpoint: " + endpoint)
            query(endpoint, cycle)
        except Exception as e:
            print(e)
    
    logging.warning("Finished Cycle " + str(cycle))
    
    cycle += 1
    time.sleep(5 * 60) # Update every 5 minutes
    
    clearTables()
    createViews()
