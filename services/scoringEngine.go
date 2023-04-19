package services

import (
	"math"
	"ronin/models"
)

const K float64 = 32


func CalculateScore(athlete1, athlete2 models.AthleteScore, athlete1Won bool) (float64, float64) {
	expectedOutcome1 := 1 / (1 + math.Pow(10, (athlete2.Score- athlete1.Score)/400))
	expectedOutcome2 := 1 / (1 + math.Pow(10, (athlete1.Score-athlete2.Score)/400))


	var outcome1, outcome2 float64
	if athlete1Won {
		outcome1 = 1
		outcome2 = 0
	} else {
		outcome1 = 0
		outcome2 = 1
	}

	updatedScore1 := athlete1.Score + K*(outcome1-expectedOutcome1)
	updatedScore2 := athlete2.Score + K*(outcome2-expectedOutcome2)

	return updatedScore1, updatedScore2
}
