# DST DataViz
This website is dedicated to players who have asked themselves "But who is the most played character now?" Along with answering this question more information (as far as what's documented and allowed) are shown - this covers mostly players, servers and server metadata, regional differences, platforms.

# Running 
```
docker-compose build
docker-compose up
```

Requires `.env` in project root that contains:
```
TOKEN=<paste token here, see [here](https://forums.kleientertainment.com/forums/topic/115578-retrieving-dst-server-data)>
POSTGRES_USER=<select postgres username>
POSTGRES_PASSWORD=<select postgress password>
POSTGRES_DB=<select database name>
DB_HOST=<"db" if running with docker-compose else most likely "localhost">
```