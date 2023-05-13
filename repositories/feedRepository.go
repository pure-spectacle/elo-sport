package repositories

import (
	"ronin/models"

	"github.com/jmoiron/sqlx"
)

type FeedRepository struct {
	DB *sqlx.DB
}

func NewFeedRepository(db *sqlx.DB) *FeedRepository {
	return &FeedRepository{
		DB: db,
	}
}

func (fr *FeedRepository) GetFeedByAthleteId(id string) ([]models.Feed, error) {
	var feed []models.Feed

	sqlStmt := `WITH athlete_following AS (
		SELECT followed_id
		FROM following
		WHERE follower_id = $1
	),
	latest_scores AS (
		SELECT athlete_id, style_id, score, updated_dt,
			   ROW_NUMBER() OVER (PARTITION BY athlete_id, style_id ORDER BY updated_dt DESC) as row_num
		FROM athlete_score
	)
	SELECT
		DISTINCT b.bout_id AS "boutId",
		c.first_name AS "challengerFirstName",
		c.last_name AS "challengerLastName",
		c.username AS "challengerUsername",
		c.athlete_id AS "challengerId",
		a.first_name AS "acceptorFirstName",
		a.athlete_id AS "acceptorId",
		a.last_name AS "acceptorLastName",
		a.username AS "acceptorUsername",
		w.first_name AS "winnerFirstName",
		w.last_name AS "winnerLastName",
		w.username AS "winnerUsername",
		l.first_name AS "loserFirstName",
		l.last_name AS "loserLastName",
		l.username AS "loserUsername",
		o.is_draw AS "isDraw",
		r.first_name AS "refereeFirstName",
		r.last_name AS "refereeLastName",
		r.athlete_id AS "refereeId",
		s.style_id AS "styleId",
		s.style_name AS "style",
		o.winner_id AS "winnerId",
		o.loser_id AS "loserId",
		b.updated_dt AS "updatedDt",
		ww.wins AS "winnerWins",
		ww.losses AS "winnerLosses",
		ww.draws AS "winnerDraws",
		ll.wins AS "loserWins",
		ll.losses AS "loserLosses",
		ll.draws AS "loserDraws",
		ws.score AS "winnerScore",
		ls.score AS "loserScore"
	FROM
		bout b
	JOIN
		athlete_following af ON b.challenger_id = af.followed_id OR b.acceptor_id = af.followed_id
	JOIN
		athlete c ON b.challenger_id = c.athlete_id
	JOIN
		athlete a ON b.acceptor_id = a.athlete_id
	LEFT JOIN
		outcome o ON b.bout_id = o.bout_id
	LEFT JOIN
		athlete w ON o.winner_id = w.athlete_id
	LEFT JOIN
		athlete l ON o.loser_id = l.athlete_id
	JOIN
		athlete r ON b.referee_id = r.athlete_id
	JOIN
		style s ON b.style_id = s.style_id
	LEFT JOIN
		athlete_record ww ON o.winner_id = ww.athlete_id
	LEFT JOIN
		athlete_record ll ON o.loser_id = ll.athlete_id
	LEFT JOIN
		latest_scores ws ON o.winner_id = ws.athlete_id AND ws.style_id = b.style_id AND ws.row_num = 1
	LEFT JOIN
		latest_scores ls ON o.loser_id = ls.athlete_id AND ls.style_id = b.style_id AND ls.row_num = 1
	WHERE b.cancelled != true and ((accepted = false and completed = false) OR (b.completed = true and b.accepted = true))
	ORDER BY b.updated_dt DESC;`

	rows, err := fr.DB.Queryx(sqlStmt, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tempFeed = models.GetFeed()
	for rows.Next() {
		err = rows.StructScan(&tempFeed)
		feed = append(feed, tempFeed)
	}

	return feed, err
}
