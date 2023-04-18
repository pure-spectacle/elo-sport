package models

type Athlete struct {
	AthleteId int    `json:"athlete_id" db:"athlete_id"`
	GymId     int    `json:"gymId" db:"gym_id"`
	FirstName string `json:"firstName" db:"first_name"`
	LastName  string `json:"lastName" db:"last_name"`
	Username  string `json:"username" db:"username"`
	BirthDate string `json:"birthDate" db:"birth_date"`
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

func CreateAthlete() Athlete {
	var athlete Athlete
	return athlete
}

func UpdateAthlete() Athlete {
	var athlete Athlete
	return athlete
}

func DeleteAthlete() Athlete {
	var athlete Athlete
	return athlete
}
