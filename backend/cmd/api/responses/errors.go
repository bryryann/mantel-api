package responses

import (
	"fmt"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/app"
)

func logError(r *http.Request, err error) {
	app := app.Get()
	app.Logger.Info(err.Error())
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	app := app.Get()

	env := map[string]any{"error": message}
	err := app.WriteJSON(w, status, env, nil)
	if err != nil {
		logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logError(r, err)

	message := "the server encountered a problem and could not process your request"
	ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	ErrorResponse(w, r, http.StatusNotFound, message)
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	ErrorResponse(w, r, http.StatusUnauthorized, message)
}
