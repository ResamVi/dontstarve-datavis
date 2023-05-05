FROM golang

RUN mkdir /app
WORKDIR /app
ADD . .
RUN go build -o main fetch/main.go

CMD ["./main"]