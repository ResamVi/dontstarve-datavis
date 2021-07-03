FROM golang

RUN mkdir /app
WORKDIR /app
ADD . .
ENV PROD="TRUE"
RUN go build -o main fetch/main.go

CMD ["./main"]