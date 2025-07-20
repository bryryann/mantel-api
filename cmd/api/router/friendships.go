package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/julienschmidt/httprouter"
)

func sendFriendRequest(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		ReceiverID int `json:"receiver_id"`
	}

	err := helpers.ReadJSON(w, r, &input)
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
	err = helpers.WriteJSON(w, http.StatusCreated, env, nil)
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

	err = helpers.ReadJSON(w, r, &input)
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

	err = helpers.WriteJSON(w, http.StatusAccepted, envelope{"friendship": patched}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}

func listPendingRequests(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	reqs, err := app.Models.Friendships.GetPendingRequests(int64(user.ID))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = helpers.WriteJSON(w, http.StatusAccepted, envelope{"requests": reqs}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
