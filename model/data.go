package model

type Attr struct {
	Topic    string `json:"topic"`
	LogPath  string `json:"log_path"`
	Service  string `json:"service"`
	SendRate int    `json:"send_rate"`
}
type Evalue struct {
	Key  string `json:"key"`
	Attr []Attr `json:"attr"`
}
