FROM golang:latest

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -o main .

COPY wait-for-it.sh .

EXPOSE 8000

CMD ["./main"]

CMD ["./wait-for-it.sh", "db:5432", "--", "./main"]