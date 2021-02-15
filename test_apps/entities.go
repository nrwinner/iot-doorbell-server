package main

type CommandPacket struct {
	Id             string
	PacketType     string
	Command        string
	Args           map[string]string
	FromId         string
	TargetDeviceId string
}

type ErrorPacket struct {
	Id             string
	PacketType     string
	Error          string
	FromId         string
	TargetDeviceId string
}
