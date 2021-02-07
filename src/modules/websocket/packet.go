package websocket

const (
	INIT_PACKET    = "init"
	COMMAND_PACKET = "command"
)

type Packet struct {
	Id         string
	PacketType string
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
