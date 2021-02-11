package entities

type Client interface {
	GetId() string
	GetRole() string
	SendMessage(message string)
	SendCommand(command Command)
	SendError(error string)
}
