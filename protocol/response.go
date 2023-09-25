package protocol

const (
	RESPONSE_OK  = 200 // Basic OK response (eg for commands PING, GET, SET)
	RESPONSE_ERR = 400 // Base error response
)

type ResponseType int

type Response struct {
	ResponseType ResponseType
	ErrorMessage string
	Value        any
}
