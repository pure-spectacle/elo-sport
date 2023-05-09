package models

type Feed struct {
	BoutId              int    `json:"boutId" db:"boutId"`
	ChallengerId        int    `json:"challengerId" db:"challengerId"`
	ChallengerFirstName string `json:"challengerFirstName" db:"challengerFirstName"`
	ChallengerLastName  string `json:"challengerLastName" db:"challengerLastName"`
	ChallengerUsername  string `json:"challengerUsername" db:"challengerUsername"`
	Style               string `json:"style" db:"style"`
	StyleId             int    `json:"styleId" db:"styleId"`
	AcceptorId          int    `json:"acceptorId" db:"acceptorId"`
	AcceptorFirstName   string `json:"acceptorFirstName" db:"acceptorFirstName"`
	AcceptorLastName    string `json:"acceptorLastName" db:"acceptorLastName"`
	AcceptorUsername    string `json:"acceptorUsername" db:"acceptorUsername"`
	RefereeId           int    `json:"refereeId" db:"refereeId"`
	RefereeFirstName    string `json:"refereeFirstName" db:"refereeFirstName"`
	RefereeLastName     string `json:"refereeLastName" db:"refereeLastName"`
	Date                string `json:"date" db:"updatedDt"`
	WinnerScore         int    `json:"winnerScore" db:"winnerScore"`
	LoserScore          int    `json:"loserScore" db:"loserScore"`
	WinnerId            int    `json:"winnerId" db:"winnerId"`
	WinnerWins          int    `json:"winnerWins" db:"winnerWins"`
	WinnerLosses        int    `json:"winnerLosses" db:"winnerLosses"`
	WinnerDraws         int    `json:"winnerDraws" db:"winnerDraws"`
	LoserWins           int    `json:"loserWins" db:"loserWins"`
	LoserLosses         int    `json:"loserLosses" db:"loserLosses"`
	LoserDraws          int    `json:"loserDraws" db:"loserDraws"`
	LoserId             int    `json:"loserId" db:"loserId"`
	WinnerFirstName     string `json:"winnerFirstName" db:"winnerFirstName"`
	WinnerLastName      string `json:"winnerLastName" db:"winnerLastName"`
	WinnerUsername      string `json:"winnerUsername" db:"winnerUsername"`
	LoserFirstName      string `json:"loserFirstName" db:"loserFirstName"`
	LoserLastName       string `json:"loserLastName" db:"loserLastName"`
	LoserUsername       string `json:"loserUsername" db:"loserUsername"`
	IsDraw              bool   `json:"isDraw" db:"isDraw"`
}

func GetFeed() Feed {
	var feed Feed
	return feed
}
