package models

type Gym struct {
	GymId   	int    `json:"gymId" db:"gym_id"`
	Name    	string `json:"name" db:"gym_name"`
	Address 	string `json:"address" db:"gym_address"`
	City    	string `json:"city" db:"gym_city"`
	State   	string `json:"state" db:"gym_state"`
	Zip    		string `json:"zip" db:"gym_zip"`
	Phone 		string `json:"phone" db:"gym_phone"`
	Email 		string `json:"email" db:"gym_email"`
	Website 	string `json:"website" db:"gym_website"`
	Description string `json:"description" db:"gym_description"`
	CreatedDate string `json:"createdDate" db:"created_dt"`
	UpdatedDate string `json:"updatedDate" db:"updated_dt"`
}

func GetGym() Gym {
	var gym Gym
	return gym
}

func GetGyms() []Gym {
	var gyms []Gym
	return gyms
}

func CreateGym() Gym {
	var gym Gym
	return gym
}

func UpdateGym() Gym {
	var gym Gym
	return gym
}

func DeleteGym() Gym {
	var gym Gym
	return gym
}
