package dto

type Activities []Activity

type Activity struct {
	Key      string `json:"key"`
	Activity string `json:"activity"`
}
