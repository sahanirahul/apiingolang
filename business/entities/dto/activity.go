package dto

type Activity struct {
	Key      string `json:"key"`
	Activity string `json:"activity"`
}

type Activities []Activity
