package models

type AthleteRecord struct {
	AthleteId int `json:"athlete_id" db:"athlete_id"`
	Wins      int `json:"wins" db:"wins"`
	Losses    int `json:"losses" db:"losses"`
	Draws     int `json:"draws" db:"draws"`
	CreatedDate string `json:"createdDate" db:"created_dt"`
	UpdatedDate string `json:"updatedDate" db:"updated_dt"`
}

func GetAthleteRecord() AthleteRecord {
	var athleteRecord AthleteRecord
	return athleteRecord
}

func GetAthleteRecords() []AthleteRecord {
	var athleteRecords []AthleteRecord
	return athleteRecords
}

func CreateAthleteRecord() AthleteRecord {
	var athleteRecord AthleteRecord
	return athleteRecord
}

func UpdateAthleteRecord() AthleteRecord {
	var athleteRecord AthleteRecord
	return athleteRecord
}

func DeleteAthleteRecord() AthleteRecord {
	var athleteRecord AthleteRecord
	return athleteRecord
}
