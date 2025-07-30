package router

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/bryryann/mantel/backend/cmd/api/app"
	"github.com/bryryann/mantel/backend/cmd/api/config"
	"github.com/bryryann/mantel/backend/cmd/api/jsonhttp"
	"github.com/bryryann/mantel/backend/cmd/api/responses"
	"github.com/bryryann/mantel/backend/internal/data"
	"github.com/bryryann/mantel/backend/internal/validator"
	"github.com/pascaldekloe/jwt"
)

// authenticateToken handles the generation of authentication tokens for users.
// It validates the input credentials, checks the user's password, and generates a JWT token
// if the credentials are valid.
func authenticateToken(w http.ResponseWriter, r *http.Request) {
	app := app.Get()
	cfg := config.Load()
	res := responses.Get()

	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := jsonhttp.ReadJSON(w, r, &input)
	if err != nil {
		res.BadRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		res.FailedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.Models.Users.GetByUsername(input.Username)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			res.InvalidCredentialsResponse(w, r)
		default:
			res.ServerErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	if !match {
		res.InvalidCredentialsResponse(w, r)
		return
	}

	var claims jwt.Claims
	claims.Subject = strconv.FormatInt(user.ID, 10)
	claims.Issued = jwt.NewNumericTime(time.Now())
	claims.NotBefore = jwt.NewNumericTime(time.Now())
	claims.Expires = jwt.NewNumericTime(time.Now().Add(24 * time.Hour))
	claims.Issuer = cfg.JWT.Issuer
	claims.Audiences = []string{cfg.JWT.Audience}

	jwtBytes, err := claims.HMACSign(jwt.HS256, []byte(cfg.JWT.Secret))
	if err != nil {
		res.ServerErrorResponse(w, r, err)
		return
	}

	err = jsonhttp.WriteJSON(w, http.StatusCreated, envelope{"authentication_token": string(jwtBytes)}, nil)
	if err != nil {
		res.ServerErrorResponse(w, r, err)
	}
}
