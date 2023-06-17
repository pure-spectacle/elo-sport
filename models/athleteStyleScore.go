package models

type AthleteStyleScore struct {
	Score     string `json:"score" db:"score"`
	StyleName string `json:"styleName" db:"style_name"`
}

func GetAthleteStyleScore() AthleteStyleScore {
	var athleteStyleScore AthleteStyleScore
	return athleteStyleScore
}

func GetAthleteStyleScores() []AthleteStyleScore {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores
}

func (athleteStyleScore *AthleteStyleScore) SetScore(score string) {
	athleteStyleScore.Score = score
}

func (athleteStyleScore *AthleteStyleScore) SetStyleName(styleName string) {
	athleteStyleScore.StyleName = styleName
}

func CreateAthleteStyleScore() AthleteStyleScore {

	var athleteStyleScore AthleteStyleScore
	return athleteStyleScore
}

func UpdateAthleteStyleScore() AthleteStyleScore {
	var athleteStyleScore AthleteStyleScore
	return athleteStyleScore
}

func DeleteAthleteStyleScore() AthleteStyleScore {

	var athleteStyleScore AthleteStyleScore
	return athleteStyleScore
}

func (athleteStyleScore *AthleteStyleScore) GetScore() string {
	return athleteStyleScore.Score
}

func (athleteStyleScore *AthleteStyleScore) GetStyleName() string {
	return athleteStyleScore.StyleName
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScore() AthleteStyleScore {
	return *athleteStyleScore
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScores() []AthleteStyleScore {

	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores
}

func (athleteStyleScore *AthleteStyleScore) SetAthleteStyleScore(athleteStyleScore2 AthleteStyleScore) {
	*athleteStyleScore = athleteStyleScore2
}

func (athleteStyleScore *AthleteStyleScore) SetAthleteStyleScores(athleteStyleScores []AthleteStyleScore) {
	*athleteStyleScore = athleteStyleScores[0]
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoreById(id string) (AthleteStyleScore, error) {
	var athleteStyleScore2 AthleteStyleScore
	return athleteStyleScore2, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresById(id string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresByAthleteId(id string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresByStyleId(id string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresByScore(score string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresByStyleName(styleName string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}

func (athleteStyleScore *AthleteStyleScore) GetAthleteStyleScoresByAthleteIdAndStyleId(athleteId string, styleId string) ([]AthleteStyleScore, error) {
	var athleteStyleScores []AthleteStyleScore
	return athleteStyleScores, nil
}
