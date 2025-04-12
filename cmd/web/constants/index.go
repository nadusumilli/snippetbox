package constants

var PORT = ":4000"

// Database connectiong string for local.
var DATABASE_CONNECTION_STRING = "user=web password=snippet@123 dbname=snippetbox sslmode=disable"

type SuccessResponse struct {
	Message string                 `json:"message"`
	SDESC   string                 `json:"sdesc"`
	SCODE   string                 `json:"scode"`
	DATA    map[string]interface{} `json:"data,omitempty"`
}

type ErrorResponse struct {
	Error string                 `json:"error"`
	SDESC string                 `json:"sdesc"`
	SCODE string                 `json:"scode"`
	DATA  map[string]interface{} `json:"data,omitempty"`
}

var (
	ErrBadRequest      = "Bad Request"
	ErrCannotBeBlank   = "This field cannot be blank"
	ErrInvalidEmail    = "This field must be a valid email address"
	ErrInvalidUsername = "This field must be a valid username"
	ErrInvalidPassword = "This field must be a valid password"
	ErrMinChars        = "This field must be more than %d characters long"
	ErrMaxChars        = "This field must be less than %d characters long"
)
