package main

type CommandPacket struct {
	Id         string
	PacketType string
	Command    string
	Args       []string
}
