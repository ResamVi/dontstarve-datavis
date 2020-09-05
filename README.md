![image](https://user-images.githubusercontent.com/6261556/90813771-54418b00-e328-11ea-98c9-f31e26f69a50.png)

# DST DataViz
This website is dedicated to players who have asked themselves "But who is the most played character now?" Along with answering this question more information (as far as what's documented and allowed) are shown - this covers mostly players, servers and server metadata, regional differences, platforms.

*[Check it out here](https://dst.resamvi.io/#/)*

# Running 
```
docker-compose build
docker-compose up
```

Requires `.env` in project root that contains:
```
TOKEN               = <paste token here, see https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data>
POSTGRES_USER       = <select postgres username>
POSTGRES_PASSWORD   = <select postgress password>
POSTGRES_MULTIPLE_DATABASES = <what DB_SHORT is>,<what DB_LONG is>
DB_SHORT            = <db name>
DB_LONG             = <db name>
DB_HOST             = <"db" if running with docker-compose else most likely "localhost">
ROCKET_ENV          = <"dev" or "prod">
ROCKET_DATABASES    = { shortterm = { url = "postgres://<POSTGRES_USER>:<POSTGRES_PASSWORD>@<DB_HOST>:5432/<DB_SHORT>" }, longterm = { url = "postgres://<POSTGRES_USER>:<POSTGRES_PASSWORD>@d<DB_HOST>:5432/<DB_LONG>" } }
VUE_APP_ENDPOINT    = <address the frontend has to call for data e.g. "http://localhost:3000">
```
