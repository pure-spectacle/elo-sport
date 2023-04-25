package models

type AthleteUsername struct {
	Username string `json:"username" db:"username"`
}

func GetAllAthleteUsernames() []AthleteUsername {
	var usernames []AthleteUsername
	return usernames
}

func GetAthleteUsername() AthleteUsername {
	var username AthleteUsername
	return username
}
