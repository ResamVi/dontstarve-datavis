import datetime
import logging

from pony import orm

db = orm.Database()

class Server(db.Entity):
    id              = orm.PrimaryKey(int, auto=True)
    name            = orm.Required(str)
    origin          = orm.Required(str)
    platform        = orm.Required(str)
    connected       = orm.Required(int)
    maxconnections  = orm.Required(int)
    mode            = orm.Optional(str)
    season          = orm.Optional(str)
    intent          = orm.Optional(str)
    mods            = orm.Optional(bool)
    cycle           = orm.Required(int)
    date            = orm.Required(datetime.datetime)
    players         = orm.Set("Player")
    
class Player(db.Entity):
    id          = orm.PrimaryKey(int, auto=True)
    cycle       = orm.Required(int)
    name        = orm.Required(str)
    character   = orm.Optional(str)
    server      = orm.Required(Server)
