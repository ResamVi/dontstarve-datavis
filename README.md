![image](https://user-images.githubusercontent.com/6261556/90813771-54418b00-e328-11ea-98c9-f31e26f69a50.png)

# DST DataViz
This website is dedicated to players who have asked themselves "But who is the most played character now?" Along with answering this question more information (as far as what's documented and allowed) are shown - this covers mostly players, servers and server metadata, regional differences, platforms.

*[Check it out here](https://dst.resamvi.io/#/)*

# Running 
Setup Discord hook to get notifications

`backend/alert/alerts.go`
```
discord.WebhookURL = "https://discord.com/api/webhooks/..."
```

Insert Token 

`docker-compose.yml`
```
TOKEN=<YOUR TOKEN HERE>
```

Set public URL to access backend (from frontend)

`web/Dockerfile`
```
ENV VUE_APP_SERVER_ENDPOINT="https://dststat.resamvi.io"
```

Start Running

```
docker-compose build
docker-compose up
```