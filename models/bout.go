package models

type Bout struct {
	BoutId       int    `json:"boutId" db:"bout_id"`
	ChallengerId int    `json:"challengerId" db:"challenger_id"`
	AcceptorId   int    `json:"acceptorId" db:"acceptor_id"`
	RefereeId    int    `json:"refereeId" db:"referee_id"`
	StyleId      int    `json:"styleId" db:"style_id"`
	Accepted     bool   `json:"accepted" db:"accepted"`
	Completed    bool   `json:"completed" db:"completed"`
	Cancelled    bool   `json:"cancelled" db:"cancelled"`
	Points       int    `json:"points" db:"points"`
	CreatedDate  string `json:"createdDate" db:"created_dt"`
	UpdatedDate  string `json:"updatedDate" db:"updated_dt"`
}

func GetBout() Bout {
	var bout Bout
	return bout
}

func GetBouts() []Bout {
	var bouts []Bout
	return bouts
}

func CreateBout() Bout {
	var bout Bout
	return bout
}

func UpdateBout() Bout {
	var bout Bout
	return bout
}

func DeleteBout() Bout {
	var bout Bout
	return bout
}
