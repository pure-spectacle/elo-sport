package services

import (
	"encoding/json"
	"net/http"
	"ronin/models"
	"ronin/repositories"

	"github.com/gorilla/mux"
)

var boutRepo *repositories.BoutRepository

func SetBoutRepo(r *repositories.BoutRepository) {
	boutRepo = r
}

func GetAllBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var bouts = models.GetBouts()
	bouts, err := boutRepo.GetAllBouts()
	if err == nil {
		json.NewEncoder(w).Encode(&bouts)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts = models.GetBout()
	vars := mux.Vars(r)
	id := vars["bout_id"]
	bouts, err := boutRepo.GetBoutById(id)
	if err == nil {
		json.NewEncoder(w).Encode(&bouts)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func CreateBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bout = models.GetBout()
	_ = json.NewDecoder(r.Body).Decode(&bout)
	if bout.ChallengerId != bout.AcceptorId {
		boutId, err := boutRepo.CreateBout(bout)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		var outboundBout models.OutboundBout
		outboundBout, err = boutRepo.GetOutboundBoutByBoutId(boutId)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		json.NewEncoder(w).Encode(&outboundBout)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		errorMessage := map[string]string{
			"error": "ChallengerId and AcceptorId must be different. You cannot create a bout against yourself.",
		}
		json.NewEncoder(w).
			Encode(errorMessage)
	}
}

func UpdateBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bout = models.GetBout()
	_ = json.NewDecoder(r.Body).Decode(&bout)
	vars := mux.Vars(r)
	id := vars["bout_id"]
	err := boutRepo.UpdateBout(id, bout)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&bout)
}

func DeleteBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	err := boutRepo.DeleteBout(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func AcceptBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	err := boutRepo.AcceptBout(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func DeclineBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["bout_id"]
	err := boutRepo.DeclineBout(id)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&id)
}

func CompleteBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	refereeId := vars["referee_id"]
	err := boutRepo.CompleteBout(boutId, refereeId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&boutId)
}

func GetPendingBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts []models.OutboundBout
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	bouts, err := boutRepo.GetPendingBoutsByAthleteId(id)
	if err == nil {
		json.NewEncoder(w).Encode(&bouts)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func GetIncompleteBouts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var bouts []models.OutboundBout
	vars := mux.Vars(r)
	id := vars["athlete_id"]
	bouts, err := boutRepo.GetIncompleteBoutsByAthleteId(id)
	if err == nil {
		json.NewEncoder(w).Encode(&bouts)
	} else {
		http.Error(w, err.Error(), 400)
		return
	}
}

func CancelBout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	boutId := vars["bout_id"]
	athleteId := vars["challenger_id"]
	err := boutRepo.CancelBout(boutId, athleteId)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	json.NewEncoder(w).Encode(&boutId)
}
