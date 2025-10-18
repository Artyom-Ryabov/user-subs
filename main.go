package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"usersubs/db"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	log.SetPrefix("[INFO]: ")
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error: something wrong with getting .env")
	}
	port, exist := os.LookupEnv("SERVER_PORT")
	if !exist {
		panic("Error: cant get field PORT from .env")
	}
	connection, exist := os.LookupEnv("DB_CONNECTION")
	if !exist {
		panic("Error: cant get field DB_CONNECTION from .env")
	}

	sqlDB, err := sql.Open("postgres", connection)
	query := db.New(sqlDB)

	err = sqlDB.Ping()
	if err != nil {
		log.Printf("CONNECT DB - something went wrong - %v\n", err)
		return 
	}

	mux := http.NewServeMux()
	server := http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	ctx := context.Background()

	mux.HandleFunc("GET /api/subs", func(w http.ResponseWriter, r *http.Request) {
		log.Println("GET /api/subs - Receive request")
		subs, err := query.GetSubs(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: something went wrong on sql query"))
			log.Printf("Error: something went wrong on sql query - %v", err)
			return
		}
		
		if err := json.NewEncoder(w).Encode(subs); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: something went wrong on encoding json"))
			log.Printf("Error: something went wrong on encoding json - %v", err)
			return
		}

		w.WriteHeader(http.StatusOK)
		log.Printf("GET /api/subs - Send response - %v\n", subs)
	})

	mux.HandleFunc("POST /api/sub", func(w http.ResponseWriter, r *http.Request) {
		log.Println("POST /api/sub - Receive request")
		var sub Sub
		if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: something went wrong on decoding json"))
			log.Printf("Error: something went wrong on decoding json - %v", err)
			return
		}

		params := db.AddSubParams{
			ServiceName: sub.ServiceName,
			Price: sub.Price,
			UserID: sub.UserID,
			StartedAt: sub.StartedAt.Time,
		}

		id, err := query.AddSub(ctx, params)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: something went wrong on sql query"))
			log.Printf("Error: something went wrong on sql query - %v", err)
			return 
		}
		
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(fmt.Sprintf("Added new Subscription with id = %v", id)))
		log.Printf("POST /api/sub - Send response - %v\n", id)
	})

	log.Printf("Server starts at port: %v\n", port)
	log.Fatal(server.ListenAndServe())
}
