import datetime
import logging

from pony import orm

db = orm.Database()

class Server(db.Entity):
    id              = orm.PrimaryKey(int, auto=True)
    name            = orm.Required(str)
    origin          = orm.Required(str)
    platform        = orm.Required(str)
    current_online  = orm.Required(int)
    max_online      = orm.Required(int)
    mode            = orm.Required(str)
    season          = orm.Required(str)
    cycle           = orm.Required(int)
    date            = orm.Required(datetime.datetime)
    players         = orm.Set("Player")
    
class Player(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Required(str)
    character   = orm.Required(str)
    server      = orm.Required(Server)
