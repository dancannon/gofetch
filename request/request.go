package request

type Requester interface {
	Send(url string) (string, error)
}
