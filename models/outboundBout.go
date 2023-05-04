package models

type OutboundBout struct {
	BoutId              int    `json:"boutId" db:"boutId"`
	ChallengerId        int    `json:"challengerId" db:"challengerId"`
	ChallengerFirstName string `json:"challengerFirstName" db:"challengerFirstName"`
	ChallengerLastName  string `json:"challengerLastName" db:"challengerLastName"`
	Style               string `json:"style" db:"style"`
	StyleId             int    `json:"styleId" db:"styleId"`
	ChallengerScore     int    `json:"challengerScore" db:"challengerScore"`
	AcceptorId          int    `json:"acceptorId" db:"acceptorId"`
	AcceptorFirstName   string `json:"acceptorFirstName" db:"acceptorFirstName"`
	AcceptorLastName    string `json:"acceptorLastName" db:"acceptorLastName"`
	AcceptorScore       int    `json:"acceptorScore" db:"acceptorScore"`
	RefereeId           int    `json:"refereeId" db:"refereeId"`
	RefereeFirstName    string `json:"refereeFirstName" db:"refereeFirstName"`
	RefereeLastName     string `json:"refereeLastName" db:"refereeLastName"`
}

func GetOutboundBout() OutboundBout {
	var outboundBout OutboundBout
	return outboundBout
}

func GetOutboundBouts() []OutboundBout {
	var outboundBouts []OutboundBout
	return outboundBouts
}

func CreateOutboundBout() OutboundBout {
	var outboundBout OutboundBout
	return outboundBout
}

func UpdateOutboundBout() OutboundBout {
	var outboundBout OutboundBout
	return outboundBout
}

func DeleteOutboundBout() OutboundBout {
	var outboundBout OutboundBout
	return outboundBout
}
