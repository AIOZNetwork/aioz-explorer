package domain

import "net/http"

type Error struct {
	HttpStatus   int
	Code         string // for i18N, later
	ErrorMessage string `json:"messages"`
	Success      bool   `json:"success"`
}

// Error is required by error interface.
func (e Error) Error() string {
	return e.ErrorMessage
}

var (
	// ErrInternalServerError will throw if any the Internal Server Error happen
	ErrInternalServerError = Error{http.StatusInternalServerError, "EPIC_FAIL", "Internal Server Error", false}

	// ErrNotFound will throw if the requested item is not exists
	ErrNotFound = Error{http.StatusNotFound, "NOT_FOUND", "Your requested Item is not found", false}

	// ErrConflict will throw if the current action already exists
	ErrConflict = Error{http.StatusConflict, "ALREADY_EXISTS", "Your Item already exist", false}

	// ErrBadParamInput will throw if the given request-body or params is not valid
	ErrBadParamInput = Error{http.StatusBadRequest, "BAD_REQUEST", "Given Param is not valid", false}

	ErrInvalidSignupToken = Error{http.StatusConflict, "INVALID_SIGNUP_TOKEN", "Invalid token for signup", false}

	ErrNotAuthorisedForOperation = Error{http.StatusBadRequest, "NOT_AUTHORISED_FOR_OPERATION",
		"User does not have authority to carry out operation", false}

	ErrMaxProfileReached = Error{http.StatusBadRequest, "MAX_PROFILE_REACHED", "User reached maximum number of profile", false}

	ErrUserNotFound = Error{http.StatusNotFound, "USER_NOT_FOUND", "User not found", false}

	ErrLoginFailed = Error{http.StatusBadRequest, "LOGIN_FAILED", "Login failed", false}

	// ErrUnreachable service communication has failed (502).
	ErrUnreachable = Error{http.StatusBadGateway, "UNREACHABLE", "Unreachable", false}

	ErrIncorrectPassword = Error{http.StatusUnauthorized, "INCORRECT_PASSWORD", "Incorrect password", false}

	ErrUserProfileNotFound = Error{http.StatusNotFound, "USER_PROFILE_NOT_FOUND", "User profile not found", false}

	ErrInvalidEmail = Error{http.StatusBadRequest, "INVALID_EMAIL", "The specified email is not unique", false}

	ErrInvalidPhone = Error{http.StatusBadRequest, "INVALID_PHONE", "The specified phone is not unique", false}

	ErrIPNotFound = Error{http.StatusNotFound, "IP_NOT_FOUND", "IP address not found", false}
)

type ErrorPasswordValidation struct {
	ErrorMessage string
}

// ErrorPasswordValidation is required by error interface.
func (e *ErrorPasswordValidation) Error() string {
	return e.ErrorMessage
}

type ErrorUsernameValidation struct {
	ErrorMessage string
}

// ErrorUsernameValidation is required by error interface.
func (e *ErrorUsernameValidation) Error() string {
	return e.ErrorMessage
}

// Error to inform that the logged user does not have permission to execute an action
type ErrorNoAuthorization struct {
	ErrorMessage string
}

// ErrorNoAuthorization is required by error interface.
func (e *ErrorNoAuthorization) Error() string {
	return e.ErrorMessage
}
