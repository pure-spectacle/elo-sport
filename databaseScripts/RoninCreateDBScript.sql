BEGIN TRANSACTION;

DROP TABLE IF EXISTS athlete, gym, athlete_score, referee, bout, competition, outcome, athlete_competition, style, gym_style, style, athlete_style, referee_style CASCADE;

CREATE TABLE gym (
	gym_id serial,
    gym_name varchar(50) NOT NULL,
	CONSTRAINT PK_gym_id PRIMARY KEY (gym_id));

CREATE TABLE referee (
    referee_id serial,
    gym_id int, 
    style_id int, 
    first_name varchar,
    last_name varchar,
    CONSTRAINT PK_referee_id PRIMARY KEY (referee_id),
    CONSTRAINT FK_gym_id FOREIGN KEY (gym_id) REFERENCES gym(gym_id));

CREATE TABLE style (
	style_id serial,
    style_name varchar(20) NOT NULL,
	CONSTRAINT PK_style_id PRIMARY KEY (style_id));

CREATE TABLE athlete (
	athlete_id serial,
    gym_id int, 
	first_name varchar(20) NOT NULL,
	last_name varchar(30) NOT NULL,
    username varchar(30) NOT NULL,
	birth_date date NOT NULL,
    wins int, 
    losses int,
	CONSTRAINT PK_athlete_id PRIMARY KEY (athlete_id),
	CONSTRAINT FK_gym_id FOREIGN KEY (gym_id) REFERENCES gym(gym_id));
	
CREATE TABLE athlete_score (
    athlete_id int,
    style_id int,
    score int,
    CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id));
	
CREATE TABLE bout (
	bout_id serial,
	challenger_id int,
    acceptor_id int,
    accepted boolean,
    points int,
    CONSTRAINT PK_bout_id PRIMARY KEY (bout_id),
	CONSTRAINT FK_challenger_id FOREIGN KEY (challenger_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_acceptor_id FOREIGN KEY (acceptor_id) REFERENCES athlete(athlete_id));

CREATE TABLE outcome (
	outcome_id serial,
    winner_id int,
    loser_id int,
    disputed boolean,
    CONSTRAINT PK_outcome_id PRIMARY KEY  (outcome_id),
	CONSTRAINT FK_winner_id FOREIGN KEY (winner_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_loser_id FOREIGN KEY (loser_id) REFERENCES athlete(athlete_id));

CREATE TABLE athlete_style (
	athlete_id int,
    style_id int,
	CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id));
	
CREATE TABLE competition (
	competition_id serial,
    referee_id int,
    bout_id int,
    outcome_id int,
    CONSTRAINT PK_competition_id PRIMARY KEY (competition_id),
    CONSTRAINT FK_referee_id FOREIGN KEY (referee_id) REFERENCES referee(referee_id),
	CONSTRAINT FK_bout_id FOREIGN KEY (bout_id) REFERENCES bout(bout_id),
    CONSTRAINT FK_outcome_id FOREIGN KEY (outcome_id) REFERENCES outcome(outcome_id));

CREATE TABLE athlete_competition (
    athlete_id int,
    competition_id int,
	CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_competition_id FOREIGN KEY (competition_id) REFERENCES competition(competition_id));
	

CREATE TABLE referee_style (
    referee_id int,
    style_id int,
	CONSTRAINT FK_referee_id FOREIGN KEY (referee_id) REFERENCES referee(referee_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id));


--rollback
COMMIT;