package router

import (
	"errors"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
)

func updateBio(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	res := responses.Get()

	user := app.Context.GetUser(r)

	var input struct {
		Content *string `json:"content"`
	}

	err := jsonhttp.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	if input.Content == nil {
		res.BadRequestResponse(w, r, errors.New("no fields provided"))
		return
	}

	userBio, err := app.Models.Bios.GetByUserID(user.ID)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	userBio.Content = *input.Content

	err = app.Models.Bios.Update(userBio)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusNoContent, nil, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
