FROM golang

RUN mkdir /app
WORKDIR /app
ADD . .
ENV PROD="TRUE"
RUN go build -o main api/main.go

CMD ["./main"]