package repositories

import (
	"ronin/models"

	"github.com/jmoiron/sqlx"
)

type AthleteRepository struct {
	db *sqlx.DB
}

type AthleteUsername struct {
	Username string `json:"username" db:"username"`
}

func NewAthleteRepository(db *sqlx.DB) *AthleteRepository {
	return &AthleteRepository{
		db: db,
	}
}

func (repo *AthleteRepository) GetAllUsernames() ([]string, error) {
	var usernames []string
	var tempUsername AthleteUsername

	sqlStmt := `SELECT username FROM athlete`
	rows, err := repo.db.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempUsername)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, tempUsername.Username)
	}

	return usernames, nil
}

func (repo *AthleteRepository) GetAllAthletes() ([]models.Athlete, error) {
	var athletes []models.Athlete
	var tempAthlete models.Athlete

	sqlStmt := `SELECT * FROM athlete`
	rows, err := repo.db.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempAthlete)
		if err != nil {
			return nil, err
		}
		athletes = append(athletes, tempAthlete)
	}

	return athletes, nil
}

func (repo *AthleteRepository) GetAthleteById(id string) (models.Athlete, error) {
	var tempAthlete models.Athlete

	sqlStmt := `SELECT * FROM athlete where athlete_id = $1`
	err := repo.db.Get(&tempAthlete, sqlStmt, id)
	if err != nil {
		return tempAthlete, err
	}

	return tempAthlete, nil
}

func (repo *AthleteRepository) GetAthleteByUsername(username string) (models.Athlete, error) {
	var tempAthlete models.Athlete

	sqlStmt := `SELECT * FROM athlete where username = $1`
	err := repo.db.Get(&tempAthlete, sqlStmt, username)
	if err != nil {
		return tempAthlete, err
	}

	return tempAthlete, nil
}

func (repo *AthleteRepository) IsAuthorizedUser(athlete models.Athlete) (bool, models.Athlete, error) {
	var athleteId int
	sqlStmt := `SELECT count(*) FROM athlete where username = $1 and password = $2`
	err := repo.db.QueryRow(sqlStmt, athlete.Username, athlete.Password).Scan(&athleteId)
	if err != nil {
		return false, models.Athlete{}, err
	}

	if athleteId == 1 {
		var tempAthlete models.Athlete
		sqlStmt := `SELECT * FROM athlete where username = $1 and password = $2`
		err := repo.db.Get(&tempAthlete, sqlStmt, athlete.Username, athlete.Password)
		if err != nil {
			return true, models.Athlete{}, err
		}

		return true, tempAthlete, nil
	}

	return false, models.Athlete{}, nil
}

func (repo *AthleteRepository) CreateAthlete(athlete models.Athlete) (int, error) {
	var athleteId int
	sqlStmt := `INSERT INTO athlete (first_name, last_name, username, birth_date, email, password)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING athlete_id`
	err := repo.db.QueryRow(sqlStmt, athlete.FirstName, athlete.LastName, athlete.Username, athlete.BirthDate, athlete.Email, athlete.Password).Scan(&athleteId)
	if err != nil {
		return 0, err
	}

	sqlStmt = `INSERT INTO athlete_record (athlete_id, wins, losses, draws) VALUES ($1, 0, 0, 0)`
	_, err = repo.db.Exec(sqlStmt, athleteId)
	if err != nil {
		return 0, err
	}

	return athleteId, nil
}

func (repo *AthleteRepository) UpdateAthlete(athlete models.Athlete) error {
	sqlStmt := `UPDATE athlete SET first_name = $1, last_name = $2, username = $3, birth_date = $4, email = $5, password = $6 WHERE athlete_id = $7`
	_, err := repo.db.Exec(sqlStmt, athlete.FirstName, athlete.LastName, athlete.Username, athlete.BirthDate, athlete.Email, athlete.Password, athlete.AthleteId)
	return err
}

func (repo *AthleteRepository) DeleteAthlete(id string) error {
	sqlStmt := `DELETE FROM athlete WHERE athlete_id = $1`
	_, err := repo.db.Exec(sqlStmt, id)
	return err
}

func (repo *AthleteRepository) GetAthleteRecord(id string) (models.AthleteRecord, error) {
	var record models.AthleteRecord
	sqlStmt := `SELECT * FROM athlete_record where athlete_id = $1`
	err := repo.db.Get(&record, sqlStmt, id)
	if err != nil {
		return record, err
	}
	return record, nil
}

func (repo *AthleteRepository) FollowAthlete(follow models.Follow) error {
	sqlStmt := `INSERT INTO following (follower_id, followed_id) VALUES ($1, $2)`
	_, err := repo.db.Exec(sqlStmt, follow.FollowerId, follow.FollowedId)
	return err
}

func (repo *AthleteRepository) UnfollowAthlete(followerId, followedId int) error {
	sqlStmt := `DELETE FROM following WHERE follower_id = $1 AND followed_id = $2`
	_, err := repo.db.Exec(sqlStmt, followerId, followedId)
	return err
}

func (repo *AthleteRepository) GetAthletesFollowed(id string) ([]int, error) {
	var follows []int
	var tempFollow models.Follow
	sqlStmt := `SELECT * FROM following where follower_id = $1`
	rows, err := repo.db.Queryx(sqlStmt, id)
	if err != nil {
		return follows, err
	}
	defer rows.Close()

	for rows.Next() {
		err2 := rows.StructScan(&tempFollow)
		if err2 != nil {
			return follows, err2
		}
		follows = append(follows, tempFollow.FollowedId)
	}

	return follows, nil
}
