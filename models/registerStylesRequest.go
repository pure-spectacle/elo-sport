package models

type RegisterStylesRequest struct {
	AthleteID int   `json:"athleteId"`
	Styles    []int `json:"styles"`
}

func GetRegisterStylesRequest() RegisterStylesRequest {
	var registerStylesRequest RegisterStylesRequest
	return registerStylesRequest
}

func GetRegisterStylesRequests() []RegisterStylesRequest {
	var registerStylesRequests []RegisterStylesRequest
	return registerStylesRequests
}

func (registerStylesRequest *RegisterStylesRequest) SetAthleteID(athleteID int) {
	registerStylesRequest.AthleteID = athleteID
}

func (registerStylesRequest *RegisterStylesRequest) SetStyles(styles []int) {
	registerStylesRequest.Styles = styles
}

func (registerStylesRequest *RegisterStylesRequest) GetAthleteID() int {
	return registerStylesRequest.AthleteID
}

func (registerStylesRequest *RegisterStylesRequest) GetStyles() []int {
	return registerStylesRequest.Styles
}

func (registerStylesRequest *RegisterStylesRequest) GetRegisterStylesRequest() RegisterStylesRequest {
	return *registerStylesRequest
}

func (registerStylesRequest *RegisterStylesRequest) GetRegisterStylesRequests() []RegisterStylesRequest {
	return []RegisterStylesRequest{*registerStylesRequest}
}

func (registerStylesRequest *RegisterStylesRequest) CreateRegisterStylesRequest() RegisterStylesRequest {
	return *registerStylesRequest
}

func (registerStylesRequest *RegisterStylesRequest) UpdateRegisterStylesRequest() RegisterStylesRequest {
	return *registerStylesRequest
}

func (registerStylesRequest *RegisterStylesRequest) DeleteRegisterStylesRequest() RegisterStylesRequest {
	return *registerStylesRequest
}
