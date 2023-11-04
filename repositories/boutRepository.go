package repositories

import (
	// "log"
	"ronin/models"

	"github.com/jmoiron/sqlx"
)

type BoutRepository struct {
	DB *sqlx.DB
}

func NewBoutRepository(db *sqlx.DB) *BoutRepository {
	return &BoutRepository{
		DB: db,
	}
}

func (repo *BoutRepository) GetAllBouts() ([]models.Bout, error) {
	var bouts = models.GetBouts()
	sqlStmt := `SELECT * FROM bout`
	rows, err := repo.DB.Queryx(sqlStmt)
	if err != nil {
		return nil, err
	} else {
		for rows.Next() {
			var tempBout models.Bout
			err = rows.StructScan(&tempBout)
			if err != nil {
				return nil, err
			}
			bouts = append(bouts, tempBout)
		}
	}
	return bouts, nil
}

func (repo *BoutRepository) GetBoutById(id string) (models.Bout, error) {
	var bout models.Bout
	sqlStmt := `SELECT * FROM bout WHERE bout_id = $1`
	err := repo.DB.QueryRowx(sqlStmt, id).StructScan(&bout)
	if err != nil {
		return models.Bout{}, err
	}
	return bout, nil
}

func (repo *BoutRepository) CreateBout(bout models.Bout) (int, error) {
	sqlStmt := `INSERT INTO bout (challenger_id, acceptor_id, referee_id, style_id, accepted, completed, cancelled, points) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING bout_id`
	err := repo.DB.QueryRowx(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.RefereeId, bout.StyleId, bout.Accepted, bout.Completed, bout.Cancelled, bout.Points).Scan(&bout.BoutId)
	if err != nil {
		return 0, err
	}
	return bout.BoutId, nil
}

func (repo *BoutRepository) GetOutboundBoutByBoutId(id int) (models.OutboundBout, error) {
	var bout models.OutboundBout
	sqlStmt := `
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		cs.score AS "challengerScore",
		b.acceptor_id AS "acceptorId",
		a.first_name AS "acceptorFirstName",
		a.last_name AS "acceptorLastName",
		ascore.score AS "acceptorScore",
		r.athlete_id AS "refereeId",
		r.first_name AS "refereeFirstName",
		r.last_name AS "refereeLastName",
		s.style_id AS "styleId"
	FROM 
		bout b
	JOIN 
		athlete c ON b.challenger_id = c.athlete_id
	JOIN 
		athlete a ON b.acceptor_id = a.athlete_id
	JOIN 
		athlete_score cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id
	JOIN 
		athlete_score ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.bout_id = $1;`

	//possibly need .Scan(&outboundBout.BoutId, &outboundBout.ChallengerId, &outboundBout.ChallengerFirstName, &outboundBout.ChallengerLastName, &outboundBout.Style, &outboundBout.ChallengerScore, &outboundBout.AcceptorId, &outboundBout.AcceptorFirstName, &outboundBout.AcceptorLastName, &outboundBout.AcceptorScore, &outboundBout.RefereeId, &outboundBout.RefereeFirstName, &outboundBout.RefereeLastName, &outboundBout.StyleId) here
	err := repo.DB.QueryRowx(sqlStmt, id).StructScan(&bout)
	if err != nil {
		return models.OutboundBout{}, err
	}
	return bout, nil
}

func (repo *BoutRepository) UpdateBout(id string, bout models.Bout) error {
	sqlStmt := `UPDATE bout SET challenger_id = $1, acceptor_id = $2, referee_id =$3, accepted = $4, points = $5, style_id = $6 WHERE bout_id = $7`
	_, err := repo.DB.Exec(sqlStmt, bout.ChallengerId, bout.AcceptorId, bout.RefereeId, bout.Accepted, bout.Points, bout.StyleId, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) DeleteBout(id string) error {
	// log.Printf("BoutRepository DB: %v\n", repo.DB)
	sqlStmt := `DELETE FROM bout WHERE bout_id = $1`
	_, err := repo.DB.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) AcceptBout(id string) error {
	sqlStmt := `UPDATE bout SET accepted = true WHERE bout_id = $1`
	_, err := repo.DB.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) DeclineBout(id string) error {
	sqlStmt := `UPDATE bout SET accepted = false, completed = true WHERE bout_id = $1`
	_, err := repo.DB.Exec(sqlStmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) CompleteBoutByBoutId(boutId string) error {
	sqlStmt := `UPDATE bout SET completed = true WHERE bout_id = $1`
	_, err := repo.DB.Exec(sqlStmt, boutId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) CompleteBout(boutId, refereeId string) error {
	sqlStmt := `UPDATE bout SET completed = true WHERE bout_id = $1 and referee_id = $2`
	_, err := repo.DB.Exec(sqlStmt, boutId, refereeId)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BoutRepository) GetPendingBoutsByAthleteId(id string) ([]models.OutboundBout, error) {
	var bouts []models.OutboundBout
	sqlStmt := `WITH latest_scores AS (
		SELECT 
			athlete_id, 
			style_id, 
			score, 
			updated_dt,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		s.style_id AS "styleId",
		cs.score AS "challengerScore",
		b.acceptor_id AS "acceptorId",
		a.first_name AS "acceptorFirstName",
		a.last_name AS "acceptorLastName",
		ascore.score AS "acceptorScore",
		r.athlete_id AS "refereeId",
		r.first_name AS "refereeFirstName",
		r.last_name AS "refereeLastName"
	FROM 
		bout b
	JOIN 
		athlete c ON b.challenger_id = c.athlete_id
	JOIN 
		athlete a ON b.acceptor_id = a.athlete_id
	JOIN 
		latest_scores cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id AND cs.row_num = 1
	JOIN 
		latest_scores ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id AND ascore.row_num = 1
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.accepted = false AND b.cancelled = false AND b.completed = false AND (b.challenger_id = $1 OR b.acceptor_id = $1 OR b.referee_id = 6)`

	rows, err := repo.DB.Queryx(sqlStmt, id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var bout models.OutboundBout
		err = rows.StructScan(&bout)
		if err != nil {
			return nil, err
		}
		bouts = append(bouts, bout)
	}
	return bouts, nil

}

func (repo *BoutRepository) GetIncompleteBoutsByAthleteId(athleteId string) ([]models.OutboundBout, error) {
	var bouts []models.OutboundBout
	sqlStmt := `WITH latest_scores AS (
		SELECT 
			athlete_id, 
			style_id, 
			score, 
			updated_dt,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		s.style_id AS "styleId",
		cs.score AS "challengerScore",
		b.acceptor_id AS "acceptorId",
		a.first_name AS "acceptorFirstName",
		a.last_name AS "acceptorLastName",
		ascore.score AS "acceptorScore",
		r.athlete_id AS "refereeId",
		r.first_name AS "refereeFirstName",
		r.last_name AS "refereeLastName"
	FROM 
		bout b
	JOIN 
		athlete c ON b.challenger_id = c.athlete_id
	JOIN 
		athlete a ON b.acceptor_id = a.athlete_id
	JOIN 
		latest_scores cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id AND cs.row_num = 1
	JOIN 
		latest_scores ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id AND ascore.row_num = 1
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.accepted = true AND b.cancelled = false AND b.completed = false AND (b.challenger_id = $1 OR b.acceptor_id = $1 OR b.referee_id = $1)`

	rows, err := repo.DB.Queryx(sqlStmt, athleteId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var bout models.OutboundBout
		err = rows.StructScan(&bout)
		if err != nil {
			return nil, err
		}
		bouts = append(bouts, bout)
	}
	return bouts, nil
}

func (repo *BoutRepository) GetCompletedBoutsByAthleteId(athleteId string) ([]models.OutboundBout, error) {
	var bouts []models.OutboundBout
	sqlStmt := `WITH latest_scores AS (
		SELECT 
			athlete_id, 
			style_id, 
			score, 
			updated_dt,
			ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT 
		b.bout_id AS "boutId",
		b.challenger_id AS "challengerId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		s.style_name AS "style",
		s.style_id AS "styleId",
		cs.score AS "challengerScore",
		b.acceptor_id AS "acceptorId",
		a.first_name AS "acceptorFirstName",
		a.last_name AS "acceptorLastName",
		ascore.score AS "acceptorScore",
		r.athlete_id AS "refereeId",
		r.first_name AS "refereeFirstName",
		r.last_name AS "refereeLastName"
	FROM 
		bout b
	JOIN 
		athlete c ON b.challenger_id = c.athlete_id
	JOIN 
		athlete a ON b.acceptor_id = a.athlete_id
	JOIN 
		latest_scores cs ON b.challenger_id = cs.athlete_id AND b.style_id = cs.style_id AND cs.row_num = 1
	JOIN 
		latest_scores ascore ON b.acceptor_id = ascore.athlete_id AND b.style_id = ascore.style_id AND ascore.row_num = 1
	JOIN 
		athlete r ON b.referee_id = r.athlete_id
	JOIN 
		style s ON b.style_id = s.style_id
	WHERE 
		b.accepted = true AND b.cancelled = false AND b.completed = true AND (b.challenger_id = $1 OR b.acceptor_id = $1 OR b.referee_id = $1)`

	rows, err := repo.DB.Queryx(sqlStmt, athleteId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var bout models.OutboundBout
		err = rows.StructScan(&bout)
		if err != nil {
			return nil, err
		}
		bouts = append(bouts, bout)
	}
	return bouts, nil
}

func (repo *BoutRepository) CancelBout(boutId, challengerId string) error {
	sqlStmt := `UPDATE bout SET cancelled = true WHERE bout_id = $1 AND challenger_id = $2`
	_, err := repo.DB.Exec(sqlStmt, boutId, challengerId)
	if err != nil {
		return err
	}
	return nil
}
