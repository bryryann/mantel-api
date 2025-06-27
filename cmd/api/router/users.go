package router

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/helpers"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/bryryann/mantel/backend/internal/validator"
	"github.com/julienschmidt/httprouter"
)

func userData(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	app := app.Get()
	res := responses.Get()

	id, err := strconv.Atoi(ps.ByName("user_id"))
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	user, err := app.Models.Users.Get(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	followers, err := app.Models.Follows.GetFollowers(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	followees, err := app.Models.Follows.GetFollowees(int64(id))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if followers == nil {
		followers = []data.UserPublic{}
	}
	if followees == nil {
		followees = []data.UserPublic{}
	}

	env := envelope{"user": user, "followers": followers, "followees": followees}

	helpers.WriteJSON(w, http.StatusAccepted, env, nil)
}

// registerUser handles the HTTP request for registering a new user.
// It reads the JSON payload from the request, validates the input, and inserts the user into the database.
func registerUser(w http.ResponseWriter, r *http.Request) {
	application := app.Get() // Retrieve the application instance.
	res := responses.Get()   // Retrieve the responses instance.

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := helpers.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		res.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = application.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			res.FailedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddError("username", "a user has already registered with this username")
			res.FailedValidationResponse(w, r, v.Errors)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = helpers.WriteJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
