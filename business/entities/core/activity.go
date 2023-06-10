package core

import "apiingolang/activity/business/entities/dto"

type Activities []Activity

type Activity struct {
	Key      string `json:"key"`
	Activity string `json:"activity"`
}

func GetActivity(ac dto.Activity) Activity {
	return Activity{Key: ac.Key, Activity: ac.Activity}
}

func GetActivities(acts dto.Activities) Activities {
	activities := Activities{}
	for _, ac := range acts {
		activities = append(activities, GetActivity(ac))
	}
	return activities
}
