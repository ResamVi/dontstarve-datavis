FROM golang

RUN mkdir /app
WORKDIR /app
ADD . .
ADD GeoLite2-Country.mmdb .
ENV PROD="TRUE"
RUN go build -o main fetcher/main.go

CMD ["./main"]