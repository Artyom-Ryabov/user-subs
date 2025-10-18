FROM golang:1.24.0
WORKDIR /app
# ENV PORT=3000
# ENV DB_CONNECTION="host=db port=5432 user=postgres password=admin sslmode=disable dbname=UserSubs"
COPY . .
RUN go mod download
RUN go build -o main
EXPOSE 5500
CMD ./main
