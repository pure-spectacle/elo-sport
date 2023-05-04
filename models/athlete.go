package models

type Athlete struct {
	AthleteId	 int    `json:"athlete_id" db:"athlete_id"`
	FirstName	 string `json:"firstName" db:"first_name"`
	LastName 	 string `json:"lastName" db:"last_name"`
	Username  	string `json:"username" db:"username"`
	BirthDate 	string `json:"birthDate" db:"birth_date"`
	Email     	string `json:"email" db:"email"`
	Password  	string `json:"password" db:"password"`
	CreatedDate string `json:"createdDate" db:"created_dt"`
	UpdatedDate string `json:"updatedDate" db:"updated_dt"`
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
