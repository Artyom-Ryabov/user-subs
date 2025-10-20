package internal

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"usersubs/internal/db"
	"usersubs/internal/subs"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/swaggo/http-swagger"
	_ "usersubs/docs"
)

func Start() error {
	log.SetPrefix("[INFO]: ")
	if err := godotenv.Load(); err != nil {
		log.Printf("Error: something wrong with getting .env - %v\n", err)
		log.Printf("Trying to get environment variables not from .env\n")
	}

	query, err := startDB()
	if err != nil {
		return err
	}

	err = startServer(query)
	if err != nil {
		return err
	}

	return nil
}

func startDB() (*db.Queries, error) {
	connection, exist := os.LookupEnv("DB_CONNECTION")
	if !exist {
		return nil, errors.New("Error: cant get field DB_CONNECTION from .env")
	}

	sqlDB, err := sql.Open("postgres", connection)

	err = sqlDB.Ping()
	if err != nil {
		return nil, fmt.Errorf("CONNECT DB - something went wrong - %v", err)
	}

	return db.New(sqlDB), nil
}

func startServer(query *db.Queries) error {
	port, exist := os.LookupEnv("SERVER_PORT")
	if !exist {
		return errors.New("Error: cant get field PORT from .env")
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	handler := subs.SubsHandler{SubsRepo: query}

	mux.HandleFunc("GET /api/subs", handler.GetSubs)
	mux.HandleFunc("POST /api/sub", handler.PostSub)
	mux.HandleFunc("PUT /api/sub/{id}", handler.PutSub)
	mux.HandleFunc("DELETE /api/sub/{id}", handler.DeleteSub)
	mux.HandleFunc("DELETE /api/subs", handler.DeleteUserSubs)
	mux.HandleFunc("GET /swagger/", httpSwagger.Handler(httpSwagger.URL(fmt.Sprintf("http://localhost:%s/swagger/doc.json", port))))

	log.Printf("Server starts at port: %v\n", port)
	log.Fatal(server.ListenAndServe())
	return nil
}
