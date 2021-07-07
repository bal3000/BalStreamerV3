package storage

// db model for current playing event
type CurrentlyPlaying struct {
	Fixture    string `json:"fixture" bson:"fixture,omitempty"`
	Chromecast string `json:"chromecast" bson:"chromecast,omitempty"`
}
