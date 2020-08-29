import os
import psycopg2
import datetime
import logging

from pony import orm

db = orm.Database()

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
    
class Player(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Optional(str)
    character   = orm.Optional(str)
    country     = orm.Required(str)
    iso         = orm.Required(str)
    continent   = orm.Required(str)
    server      = orm.Required(Server)

# Track active player over time by their origin
class Activity(db.Entity): # Rename Snapshot
    id              = orm.PrimaryKey(int, auto=True)
    date            = orm.Required(datetime.datetime)
    countbyorigin   = orm.Required(orm.Json) # {China: 2991, USA: 320, Russia: 245, ...}

# Create Views (we temporarily create a connection to execute raw sql)
def createViews():
    connection = psycopg2.connect(
        port="5432",
        user=os.getenv('POSTGRES_USER'),
        password=os.getenv('POSTGRES_PASSWORD'),
        database=os.getenv('POSTGRES_DB'),
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