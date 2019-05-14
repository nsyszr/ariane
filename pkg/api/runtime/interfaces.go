package runtime

type Object interface {
	Encode() ([]byte, error)
}
