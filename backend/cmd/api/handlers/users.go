package handlers

import (
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
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
		// TODO: implement proper error handling code.
		// Until then, use panic for testing purposes.
		panic(err)
	}

	user := &data.User{
		Username: input.Username,
		Email:    input.Email,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		panic(err)
	}

	v := validator.New()
	if data.ValidateUser(v, user); !v.Valid() {
		for error := range v.Errors {
			w.Write([]byte(error))
		}

		return
	}

	err = application.Models.Users.Insert(user)
	if err != nil {
		panic(err)
	}

	err = application.WriteJSON(w, http.StatusCreated, envelope{"user": user}, nil)
	if err != nil {
		panic(err)
	}
}
