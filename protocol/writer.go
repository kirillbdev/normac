package protocol

type Writer interface {
	Write(response *Response) []byte
}
