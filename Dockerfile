FROM golang:1.24.0
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -C ./cmd -o main
EXPOSE 5500
CMD ./cmd/main
