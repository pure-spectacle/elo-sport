package models

type Outcome struct {
	OutcomeId int  `json:"outcomeId" db:"outcome_id"`
	BoutId    int  `json:"boutId" db:"bout_id"`
	WinnerId  int  `json:"winnerId" db:"winner_id"`
	LoserId   int  `json:"loserId" db:"loser_id"`
	StyleId   int  `json:"styleId" db:"style_id"`
	IsDraw    bool `json:"isDraw" db:"is_draw"`
	CreatedDate string `json:"createdDate" db:"created_dt"`
	UpdatedDate string `json:"updatedDate" db:"updated_dt"`
}

func GetOutcome() Outcome {
	var outcome Outcome
	return outcome
}

func GetOutcomes() []Outcome {
	var outcomes []Outcome
	return outcomes
}

func CreateOutcome() Outcome {
	var outcome Outcome
	return outcome
}

func UpdateOutcome() Outcome {
	var outcome Outcome
	return outcome
}

func DeleteOutcome() Outcome {
	var outcome Outcome
	return outcome
}
