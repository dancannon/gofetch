package request

type Requester interface {
	Send(url string) (map[string]interface{}, error)
}
