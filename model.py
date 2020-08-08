import datetime
import logging

from pony import orm

db = orm.Database()

class Server(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Required(str)
    platform    = orm.Required(str)
    online      = orm.Required(str)
    mode        = orm.Required(str)
    season      = orm.Required(str)
    date        = orm.Required(datetime.datetime)
    players     = orm.Set("Player")
    
class Player(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Required(str)
    character   = orm.Required(str)
    server      = orm.Required(Server)
