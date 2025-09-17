package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

func getFriendsById(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	userIdParam := ps.ByName("user_id")
	id, err := strconv.Atoi(userIdParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	query := r.URL.Query()

	page := helpers.ParseIntOrDefault(query.Get("page"), 1)
	pageSize := helpers.ParseIntOrDefault(query.Get("page_size"), 20)

	_, err = app.Models.Users.Exists(int64(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrUserNotFound):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	paginationData := data.Pagination{
		Page:     page,
		PageSize: pageSize,
	}

	friends, err := app.Models.Friendships.GetFriends(int64(id), paginationData)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	jsonResponse := envelope{
		"friends": friends,
		"meta": map[string]any{
			"page":      page,
			"page_size": pageSize,
		},
	}
	err = jsonhttp.WriteJSON(w, http.StatusAccepted, jsonResponse, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func sendFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		ReceiverID int `json:"receiver_id"`
	}

	err := jsonhttp.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	fs := &data.Friendship{
		SenderID:   user.ID,
		ReceiverID: int64(input.ReceiverID),
	}

	err = app.Models.Friendships.SendRequest(fs)
	if err != nil {
		switch err {
		case data.ErrFriendshipRequestToSelf:
			res.BadRequestResponse(w, r, err)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	env := envelope{
		"message":    "friend request sent",
		"created_at": fs.CreatedAt,
		"status":     fs.Status,
	}
	err = jsonhttp.WriteJSON(w, http.StatusCreated, env, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func patchPendingFriendRequest(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	idParam := ps.ByName("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	err = jsonhttp.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	switch input.Status {
	case "accepted", "blocked", "pending":
		// do nothing
	default:
		res.BadRequestResponse(w, r, data.ErrInvalidFriendshipStatus)
		return
	}

	fs := &data.Friendship{
		ID:         int64(id),
		ReceiverID: user.ID,
		Status:     data.FriendshipStatus(input.Status),
	}

	patched, err := app.Models.Friendships.PatchFriendship(fs)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrNoSuchRequest):
			res.NotFoundResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusAccepted, envelope{"friendship": patched}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func listPendingRequests(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	query := r.URL.Query()

	by := query.Get("by")
	switch by {
	case "received", "sent":
		// do nothing
	default:
		by = "sent"
	}

	page := helpers.ParseIntOrDefault(query.Get("page"), 1)
	pageSize := helpers.ParseIntOrDefault(query.Get("page_size"), 20)

	user := app.Context.GetUser(r)

	var (
		reqs []data.Friendship
		err  error
	)

	paginationData := data.Pagination{
		Page:     page,
		PageSize: pageSize,
	}
	if by == "received" {
		reqs, err = app.Models.Friendships.GetReceivedPendingRequests(user.ID, paginationData)
	} else {
		reqs, err = app.Models.Friendships.GetSentPendingRequests(user.ID, paginationData)
	}
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	jsonResponse := envelope{
		"requests": reqs,
		"meta": map[string]any{
			"page":      page,
			"page_size": pageSize,
		},
	}
	err = jsonhttp.WriteJSON(w, http.StatusAccepted, jsonResponse, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
