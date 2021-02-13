package main

type CommandPacket struct {
	Id         string
	PacketType string
	Command    string
	Args       map[string]string
}
