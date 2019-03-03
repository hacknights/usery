package requests

type Request struct {
}

type validator interface {
	validate() error
}
