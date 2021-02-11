package entities

type Command struct {
	Path string
	Args []string
	Responder Responder
}
