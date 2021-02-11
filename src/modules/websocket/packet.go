package websocket

import "doorbell-server/src/entities"

const (
	INIT_PACKET    = "init"
	COMMAND_PACKET = "command"
	ERROR_PACKET   = "error"
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

type ErrorPacket struct {
	Packet
	Error string
}

func CommandFromPacket(packet CommandPacket) entities.Command {
	return entities.Command{
		Path: packet.Command,
		Args: packet.Args,
	}
}

func PacketFromCommand(command entities.Command) CommandPacket {
	return CommandPacket {
		Packet: Packet{
			PacketType: COMMAND_PACKET,
		},
		Command: command.Path,
		Args: command.Args,
	}
}