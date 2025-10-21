package subs

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"
	"usersubs/internal/db"
	"usersubs/internal/utils"

	"github.com/google/uuid"
)

type subJSON struct {
	ID          int32          `json:"id,omitempty"`
	ServiceName string         `json:"service_name"`
	Price       int32          `json:"price"`
	UserID      uuid.UUID      `json:"user_id"`
	StartedAt   utils.JSONDate `json:"start_date"`
	EndedAt     utils.JSONDate `json:"end_date,omitzero"`
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
		subsDB []db.Subscription
		err    error
	)

	user_id := r.URL.Query().Get("user_id")
	if user_id != "" {
		id, err := uuid.Parse(user_id)
		if err != nil {
			utils.SendError(w, "Error: could not parse url query", http.StatusInternalServerError, err)
			return
		}

		subsDB, err = h.SubsRepo.GetUserSubs(context.Background(), id)
	} else {
		subsDB, err = h.SubsRepo.GetSubs(context.Background())
	}

	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	subs := []subJSON{}
	for _, sub := range subsDB {
		subs = append(subs, subJSON{
			ID:          sub.ID,
			ServiceName: sub.ServiceName,
			Price:       sub.Price,
			UserID:      sub.UserID,
			StartedAt:   utils.JSONDate(sub.StartedAt),
			EndedAt:     utils.JSONDate(sub.EndedAt.Time),
		})
	}

	if err := utils.SendData(w, subs, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
}

// @Summary GetSub
// @Description Get a subscription by ID
// @Produce json
// @Param id path int true "ID (int) of specific subscription"
// @Router /api/sub/{id} [GET]
func (h SubsHandler) GetSub(w http.ResponseWriter, r *http.Request) {
	log.Println("GET /api/sub/{id} - Receive request")

	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		utils.SendError(w, "Error: could not parse path value", http.StatusInternalServerError, err)
		return
	}

	sub, err := h.SubsRepo.GetSub(context.Background(), int32(subID))
	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	s := subJSON{
		ID:          int32(subID),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartedAt:   utils.JSONDate(sub.StartedAt),
		EndedAt:     utils.JSONDate(sub.EndedAt.Time),
	}

	if err := utils.SendData(w, s, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
}

// @Summary PostSub
// @Description Create a subscription
// @Accept json
// @Produce json
// @Param request body subJSON true "Structure of new subscription"
// @Router /api/sub [POST]
func (h SubsHandler) PostSub(w http.ResponseWriter, r *http.Request) {
	log.Println("POST /api/sub - Receive request")
	var sub subJSON
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		utils.SendError(w, "Error: something went wrong on decoding json", http.StatusInternalServerError, err)
		return
	}

	params := db.AddSubParams{
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartedAt:   time.Time(sub.StartedAt),
		EndedAt:     sql.NullTime{Time: time.Time(sub.EndedAt), Valid: true},
	}

	id, err := h.SubsRepo.AddSub(context.Background(), params)
	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	sub.ID = id

	if err := utils.SendData(w, sub, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
}

// @Summary PutSub
// @Description Update a subscription
// @Accept json
// @Produce json
// @Param id path int true "ID of subscription"
// @Param request body subJSON true "Structure of subscription"
// @Router /api/sub/{id} [PUT]
func (h SubsHandler) PutSub(w http.ResponseWriter, r *http.Request) {
	log.Println("PUT /api/sub/{id} - Receive request")
	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		utils.SendError(w, "Error: could not parse path value", http.StatusInternalServerError, err)
		return
	}

	var sub subJSON
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		utils.SendError(w, "Error: something went wrong on decoding json", http.StatusInternalServerError, err)
		return
	}

	params := db.UpdateSubParams{
		ID:          int32(subID),
		ServiceName: sub.ServiceName,
		Price:       sub.Price,
		UserID:      sub.UserID,
		StartedAt:   time.Time(sub.StartedAt),
		EndedAt:     sql.NullTime{Time: time.Time(sub.EndedAt), Valid: true},
		UpdatedAt:   time.Now(),
	}

	id, err := h.SubsRepo.UpdateSub(context.Background(), params)
	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	sub.ID = id

	if err := utils.SendData(w, sub, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
}

// @Summary DeleteSub
// @Description Delete a subscription
// @Produce json
// @Param id path int true "ID of subscription"
// @Router /api/sub/{id} [DELETE]
func (h SubsHandler) DeleteSub(w http.ResponseWriter, r *http.Request) {
	log.Println("DELETE /api/sub/{id} - Receive request")
	pathID := r.PathValue("id")
	subID, err := strconv.ParseInt(pathID, 10, 32)
	if err != nil {
		utils.SendError(w, "Error: could not parse path value", http.StatusInternalServerError, err)
		return
	}

	id, err := h.SubsRepo.DeleteSub(context.Background(), int32(subID))
	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	if err := utils.SendData(w, id, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
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
		utils.SendError(w, "Error: query param `user_id` is not provided", http.StatusBadRequest, errors.New("no user_id"))
		return
	}

	id, err := uuid.Parse(user_id)
	if err != nil {
		utils.SendError(w, "Error: could not parse url query", http.StatusInternalServerError, err)
		return
	}

	ids, err := h.SubsRepo.DeleteUserSubs(context.Background(), id)

	if err != nil {
		utils.SendError(w, "Error: something went wrong on sql query", http.StatusInternalServerError, err)
		return
	}

	if err := utils.SendData(w, ids, http.StatusOK); err != nil {
		utils.SendError(w, "Error: something went wrong on encoding json", http.StatusInternalServerError, err)
		return
	}
}
