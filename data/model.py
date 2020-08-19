import datetime
import logging

from pony import orm

db = orm.Database()

class Server(db.Entity):
    id              = orm.PrimaryKey(int, auto=True)
    name            = orm.Optional(str)
    origin          = orm.Required(str)
    platform        = orm.Required(str)
    connected       = orm.Required(int)
    maxconnections  = orm.Required(int)
    elapsed         = orm.Required(int)
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
    name        = orm.Optional(str)
    character   = orm.Optional(str)
    origin      = orm.Required(str)
    server      = orm.Required(Server)

# Track active player over time by their origin
class Activity(db.Entity):
    id              = orm.PrimaryKey(int, auto=True)
    date            = orm.Required(datetime.datetime)
    countbyorigin   = orm.Required(orm.Json) # {China: 2991, USA: 320, Russia: 245, ...}