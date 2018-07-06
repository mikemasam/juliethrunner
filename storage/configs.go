package storage

type Configs struct {
	Every    string   `json:"every"`
	At       []string `json:"at"`
	Commands []string `json:"tasks"`
	Id       string   `json:"id"`
}
