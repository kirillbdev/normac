package protocol

const (
	RESPONSE_OK  = 0  // Basic OK response (eg for commands PING, GET, SET)
	RESPONSE_ERR = 10 // Base error response
)

type ResponseType int

type Response struct {
	ResponseType ResponseType
	ErrorMessage string
	Value        any
}

func NewErrorResponse(message string) *Response {
	return &Response{
		ResponseType: RESPONSE_ERR,
		ErrorMessage: message,
	}
}

func NewOkResponse(val any) *Response {
	return &Response{
		ResponseType: RESPONSE_OK,
		Value:        val,
	}
}
