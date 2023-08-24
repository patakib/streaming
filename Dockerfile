FROM golang:latest
RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o streamingapp .
RUN chmod +x wait-for-it.sh