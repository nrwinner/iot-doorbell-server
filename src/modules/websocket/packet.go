package websocket

const (
	INIT_PACKET    = "init"
	COMMAND_PACKET = "command"
	ERROR_PACKET   = "error"
)

type Packet struct {
	Id         string
	PacketType string
	ResponseId string
}

type InitPacket struct {
	Packet
	Role string
	Name string
}

type CommandPacket struct {
	Packet
	Command string
	Args    []string
}

type ErrorPacket struct {
	Packet
	Error string
}
