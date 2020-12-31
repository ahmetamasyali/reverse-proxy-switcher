FROM golang:latest
WORKDIR /app
COPY . .
RUN go get -d
RUN go build -o main .
EXPOSE 9090
EXPOSE 9091
EXPOSE 9092
CMD ["./main"]