package websocket

type Controller struct {
	MessageEventHandler    func(message string, client Client)
	ConnectEventHandler    func(client Client)
	DisconnectEventHandler func(client Client)
}

func (controller *Controller) HandleMessage(message string, client Client) {
	if controller.MessageEventHandler != nil {
		controller.MessageEventHandler(message, client)
	}
}

func (controller *Controller) HandleConnect(client Client) {
	if controller.ConnectEventHandler != nil {
		controller.ConnectEventHandler(client)
	}
}

func (controller *Controller) HandleDisconnect(client Client) {
	if controller.DisconnectEventHandler != nil {
		controller.DisconnectEventHandler(client)
	}
}
