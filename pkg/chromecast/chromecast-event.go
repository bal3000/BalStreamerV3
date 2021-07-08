package chromecast

type ChromecastEvent struct {
	Chromecast  string `json:"chromecast"`
	MessageType string `json:"messageType"`
}
