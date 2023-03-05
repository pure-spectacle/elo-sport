package models

type Athlete struct {
	AthleteId int `json:"athlete_id" db:"athlete_id"`
	GymId     int    `json:"gym_id" db:"gym_id"`
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Username  string `json:"username" db:"username"`
	BirthDate string `json:"birth_date" db:"birth_date"`
	Wins      int    `json:"wins" db:"wins"`
	Losses    int    `json:"losses" db:"losses"`
}

func GetAthlete() Athlete {
	var athlete Athlete
	return athlete
}

func GetAthletes() []Athlete {
	var athletes []Athlete
	return athletes
}
