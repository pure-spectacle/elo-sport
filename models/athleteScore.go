package models

type AthleteScore struct {
	AthleteId   int     `json:"athleteId" db:"athlete_id"`
	StyleId     int     `json:"styleId" db:"style_id"`
	OutcomeId   int     `json:"outcomeId" db:"outcome_id"`
	Score       float64 `json:"score" db:"score"`
	CreatedDate string  `json:"createdDate" db:"created_dt"`
	UpdatedDate string  `json:"updatedDate" db:"updated_dt"`
}

func GetAthleteScore() AthleteScore {
	var athleteScore AthleteScore
	return athleteScore
}

func GetAthleteScores() []AthleteScore {
	var athleteScores []AthleteScore
	return athleteScores
}

func CreateAthleteScore() AthleteScore {

	var athleteScore AthleteScore
	return athleteScore
}

func UpdateAthleteScore() AthleteScore {
	var athleteScore AthleteScore
	return athleteScore
}

func DeleteAthleteScore() AthleteScore {
	var athleteScore AthleteScore
	return athleteScore
}
