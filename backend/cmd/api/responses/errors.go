package responses

import (
	"fmt"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/helpers"
)

func (res *Responses) logError(r *http.Request, err error) {
	res.logger.Info(err.Error())
}

func (res *Responses) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := map[string]any{"error": message}
	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		res.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (res *Responses) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	res.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

func (res *Responses) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	res.ErrorResponse(w, r, http.StatusNotFound, message)
}

func (res *Responses) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	res.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (res *Responses) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

func (res *Responses) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	res.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (res *Responses) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (res *Responses) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

func (res *Responses) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}
