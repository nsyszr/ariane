package api

type Object interface {
	Encode() ([]byte, error)
}
