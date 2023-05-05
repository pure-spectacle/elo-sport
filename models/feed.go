package models

type Feed struct {
	BoutId              int    `json:"boutId" db:"boutId"`
	ChallengerId        int    `json:"challengerId" db:"challengerId"`
	ChallengerFirstName string `json:"challengerFirstName" db:"challengerFirstName"`
	ChallengerLastName  string `json:"challengerLastName" db:"challengerLastName"`
	ChallengerScore     int    `json:"challengerScore" db:"challengerScore"`
	ChallengerWins      int    `json:"challengerWins" db:"challengerWins"`
	ChallengerLosses    int    `json:"challengerLosses" db:"challengerLosses"`
	ChallengerDraws     int    `json:"challengerDraws" db:"challengerDraws"`
	Style               string `json:"style" db:"style"`
	StyleId             int    `json:"styleId" db:"styleId"`
	AcceptorId          int    `json:"acceptorId" db:"acceptorId"`
	AcceptorFirstName   string `json:"acceptorFirstName" db:"acceptorFirstName"`
	AcceptorLastName    string `json:"acceptorLastName" db:"acceptorLastName"`
	AcceptorScore       int    `json:"acceptorScore" db:"acceptorScore"`
	AcceptorWins        int    `json:"acceptorWins" db:"acceptorWins"`
	AcceptorLosses      int    `json:"acceptorLosses" db:"acceptorLosses"`
	AcceptorDraws       int    `json:"acceptorDraws" db:"acceptorDraws"`
	RefereeId           int    `json:"refereeId" db:"refereeId"`
	RefereeFirstName    string `json:"refereeFirstName" db:"refereeFirstName"`
	RefereeLastName     string `json:"refereeLastName" db:"refereeLastName"`
	Date                string `json:"date" db:"updateDate"`
	WinnerOldScore      int    `json:"winnerOldScore" db:"winnerOldScore"`
	WinnerNewScore      int    `json:"winnerNewScore" db:"winnerNewScore"`
	LoserOldScore       int    `json:"loserOldScore" db:"loserOldScore"`
	LoserNewScore       int    `json:"loserNewScore" db:"loserNewScore"`
	WinnerId            int    `json:"winnerId" db:"winnerId"`
	LoserId             int    `json:"loserId" db:"loserId"`
}

func GetFeed() Feed {
	var feed Feed
	return feed
}
