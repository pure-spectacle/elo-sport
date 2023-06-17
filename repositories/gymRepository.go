package repositories

import (
	"ronin/models"

	"github.com/jmoiron/sqlx"
)

type GymRepository struct {
	DB *sqlx.DB
}

func NewGymRepository(db *sqlx.DB) *GymRepository {
	return &GymRepository{
		DB: db,
	}
}

func (repo *GymRepository) GetAllGyms() ([]models.Gym, error) {
	var gyms []models.Gym
	var tempGym models.Gym

	sqlStmt := `SELECT * FROM gym`
	rows, err := repo.DB.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.StructScan(&tempGym)
		if err != nil {
			return nil, err
		}
		gyms = append(gyms, tempGym)
	}

	return gyms, nil
}

func (repo *GymRepository) GetGymById(id string) (models.Gym, error) {
	var gym models.Gym
	sqlStmt := `SELECT * FROM gym WHERE gym_id = $1`
	err := repo.DB.QueryRowx(sqlStmt, id).StructScan(&gym)
	if err != nil {
		return models.Gym{}, err
	}
	return gym, nil
}

func (repo *GymRepository) CreateGym(gym models.Gym) (models.Gym, error) {
	sqlStmt := `INSERT INTO gym (gym_name, gym_address, gym_city, gym_state, gym_zip, gym_phone, gym_email, gym_website, gym_description) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING gym_id`
	_, err := repo.DB.Exec(sqlStmt, gym.Name, gym.Address, gym.City, gym.State, gym.Zip, gym.Phone, gym.Email, gym.Website, gym.Description)
	if err != nil {
		return models.Gym{}, err
	}
	return gym, nil
}
