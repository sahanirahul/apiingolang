package dto

import "encoding/json"

type Activities []Activity

type Activity struct {
	Key      string `json:"key"`
	Activity string `json:"activity"`
}

func (acts *Activities) String() string {
	b, _ := json.Marshal(acts)
	return string(b)
}
