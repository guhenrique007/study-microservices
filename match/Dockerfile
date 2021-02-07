FROM golang:1.15

WORKDIR /Code/src/
COPY . .
RUN GOOS=linux go build -ldflags="-s -w"
CMD ["./match"]