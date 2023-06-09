package dto

type Activity map[string]string

type Activities []Activity

type BoredApiActivityResponse struct {
	Key      string `json:"key"`
	Activity string `json:"activity"`
}

func CreateActivity(req *BoredApiActivityResponse) Activity {
	ac := Activity{}
	ac[req.Key] = req.Activity
	return ac
}
