package handlers

import (
	"errors"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/bryryann/mantel/backend/internal/validator"
)

type envelope map[string]any

func registerUserHandler(w http.ResponseWriter, r *http.Request) {
	application := app.Get()

	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := application.ReadJSON(w, r, &input)
	if err != nil {
		responses.BadRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		responses.ServerErrorResponse(w, r, err)
		return
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		responses.FailedValidationResponse(w, r, v.Errors)
		return
	}

	err = application.Models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			responses.FailedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddError("username", "a user has already registered with this username")
			responses.FailedValidationResponse(w, r, v.Errors)
		default:
			responses.ServerErrorResponse(w, r, err)
		}
		return
	}

	err = application.WriteJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		responses.ServerErrorResponse(w, r, err)
	}
}
