package responses

import (
	"fmt"
	"net/http"

	"github.com/bryryann/mantel/backend/cmd/api/helpers"
)

// logError logs the provided error using the logger associated with the Responses struct.
func (res *Responses) logError(_ *http.Request, err error) {
	res.logger.Info(err.Error())
}

// ErrorResponse sends a JSON-formatted error response with the specified status code and message.
// If an error occurs while writing the JSON, it logs the error and sends an internal server error response.
func (res *Responses) ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := map[string]any{"error": message}
	err := helpers.WriteJSON(w, status, env, nil)
	if err != nil {
		res.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

// ServerErrorResponse sends a 500 Internal Server Error response with a generic error message.
// It also logs the provided error.
func (res *Responses) ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.logError(r, err)

	message := "the server encountered a problem and could not process your request"
	res.ErrorResponse(w, r, http.StatusInternalServerError, message)
}

// NotFoundResponse sends a 404 Not Found response with a generic error message.
func (res *Responses) NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	res.ErrorResponse(w, r, http.StatusNotFound, message)
}

// MethodNotAllowedResponse sends a 405 Method Not Allowed response with a message indicating the unsupported method.
func (res *Responses) MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request, err error) {
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	res.ErrorResponse(w, r, http.StatusMethodNotAllowed, message)
}

// BadRequestResponse sends a 400 Bad Request response with the provided error message.
func (res *Responses) BadRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	res.ErrorResponse(w, r, http.StatusBadRequest, err.Error())
}

// FailedValidationResponse sends a 422 Unprocessable Entity response with the provided validation errors.
func (res *Responses) FailedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	res.ErrorResponse(w, r, http.StatusUnprocessableEntity, errors)
}

// InvalidCredentialsResponse sends a 401 Unauthorized response with a generic invalid credentials message.
func (res *Responses) InvalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// InvalidAuthenticationTokenResponse sends a 401 Unauthorized response with a message about an invalid or missing token.
// It also sets the WWW-Authenticate header to "Bearer".
func (res *Responses) InvalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authenticate", "Bearer")

	message := "invalid or missing authentication token"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}

// AuthenticationRequiredResponse sends a 401 Unauthorized response with a message indicating authentication is required.
func (res *Responses) AuthenticationRequiredResponse(w http.ResponseWriter, r *http.Request) {
	message := "you must be authenticated to access this resource"
	res.ErrorResponse(w, r, http.StatusUnauthorized, message)
}
