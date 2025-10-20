package subs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"usersubs/internal/db"
	"usersubs/internal/utils"

	"github.com/google/uuid"
)

type sub struct {
	ServiceName string         `json:"service_name"`
	Price       int32          `josn:"price"`
	UserID      uuid.UUID      `json:"user_id"`
	StartedAt   utils.JSONDate `json:"start_date"`
}

type SubsHandler struct {
	SubsRepo *db.Queries
}

// @Summary GetSubs
// @Description Get all subscriptions
// @Produce json
// @Param user_id query string false "User ID, if need to get all subscriptions of a specific user"
// @Router /api/subs [GET]
func (h SubsHandler) GetSubs(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/subs - Receive request")
	var (
		subs []db.Subscription
		err  error
	)

	user_id := r.URL.Query().Get("user_id")
	if user_id != "" {
		id, err := uuid.Parse(user_id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error: could not parse url query"))
			log.Printf("Error: could not parse url query - %v", err)
			return
		}

		subs, err = h.SubsRepo.GetUserSubs(context.Background(), id)
	} else {
		subs, err = h.SubsRepo.GetSubs(context.Background())
	}

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
}

// @Summary GetSub
// @Description Get a subscription by ID
// @Produce json
// @Param id path int true "ID (int) of specific subscription"
// @Router /api/sub [GET]
func (h SubsHandler) GetSub(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/sub/{id} - Receive request")

	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: could not parse path value"))
		log.Printf("Error: could not parse path value - %v", err)
		return
	}

	sub, err := h.SubsRepo.GetSub(context.Background(), int32(subID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on sql query"))
		log.Printf("Error: something went wrong on sql query - %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(sub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on encoding json"))
		log.Printf("Error: something went wrong on encoding json - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("GET /api/sub - Send response - %v\n", sub)
}

// @Summary PostSub
// @Description Create a subscription
// @Accept json
// @Produce json
// @Param request body sub true "Structure of new subscription"
// @Router /api/sub [POST]
func (h SubsHandler) PostSub(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/sub - Receive request")
	var sub sub
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on decoding json"))
		log.Printf("Error: something went wrong on decoding json - %v", err)
		return
	}

	params := db.AddSubParams{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartedAt:   sub.StartedAt.Time,
	}

	id, err := h.SubsRepo.AddSub(context.Background(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on sql query"))
		log.Printf("Error: something went wrong on sql query - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Added new Subscription with id = %v", id)))
	log.Printf("POST /api/sub - Send response - %v\n", id)
}

// @Summary PutSub
// @Description Update a subscription
// @Accept json
// @Produce json
// @Param id path int true "ID of subscription"
// @Param request body sub true "Structure of subscription"
// @Router /api/sub [PUT]
func (h SubsHandler) PutSub(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /api/sub/{id} - Receive request")
	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: could not parse path value"))
		log.Printf("Error: could not parse path value - %v", err)
		return
	}

	var sub sub
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on decoding json"))
		log.Printf("Error: something went wrong on decoding json - %v", err)
		return
	}

	params := db.UpdateSubParams{
		ID:          int32(subID),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartedAt:   sub.StartedAt.Time,
		UpdatedAt:   time.Now(),
	}

	id, err := h.SubsRepo.UpdateSub(context.Background(), params)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on sql query"))
		log.Printf("Error: something went wrong on sql query - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Updated Subscription with id = %v", id)))
	log.Printf("PUT /api/sub/{id} - Send response - %v\n", id)
}

// @Summary DeleteSub
// @Description Delete a subscription
// @Produce json
// @Param id path int true "ID of subscription"
// @Router /api/sub [DELETE]
func (h SubsHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /api/sub/{id} - Receive request")
	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: could not parse path value"))
		log.Printf("Error: could not parse path value - %v", err)
		return
	}

	id, err := h.SubsRepo.DeleteSub(context.Background(), int32(subID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on sql query"))
		log.Printf("Error: something went wrong on sql query - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Deleted Subscription with id = %v", id)))
	log.Printf("DELETE /api/sub/{id} - Send response - %v\n", id)
}

// @Summary DeleteUserSubs
// @Description Delete all subscriptions of specific user
// @Produce json
// @Param user_id query string true "User ID"
// @Router /api/subs [DELETE]
func (h SubsHandler) DeleteUserSubs(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /api/subs - Receive request")
	user_id := r.URL.Query().Get("user_id")
	if user_id == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Error: query param `user_id` is not provided"))
		log.Println("Error: query param `user_id` is not provided")
		return
	}

	id, err := uuid.Parse(user_id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: could not parse url query"))
		log.Printf("Error: could not parse url query - %v", err)
		return
	}

	ids, err := h.SubsRepo.DeleteUserSubs(context.Background(), id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on sql query"))
		log.Printf("Error: something went wrong on sql query - %v", err)
		return
	}

	if err := json.NewEncoder(w).Encode(ids); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error: something went wrong on encoding json"))
		log.Printf("Error: something went wrong on encoding json - %v", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	log.Printf("DELETE /api/subs - Send response - %v\n", ids)
}
