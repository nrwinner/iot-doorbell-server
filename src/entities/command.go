package entities

type Command struct {
	Path           string
	Args           map[string]string
	Responder      Responder
	FromId         string
	TargetDeviceId string
}
