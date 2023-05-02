BEGIN TRANSACTION;

DROP TABLE IF EXISTS athlete, gym, referee, style, athlete, athlete_record, athlete_gym, athlete_score, bout, outcome, athlete_style, competition, athlete_competition, referee_style CASCADE; 

CREATE TABLE gym (
	gym_id serial PRIMARY KEY,
    gym_name varchar(100) NOT NULL,
    gym_address varchar(500),
    gym_city varchar(100),
    gym_state varchar(2),
    gym_zip varchar(5),
    gym_phone varchar(10) NOT NULL,
    gym_email varchar(100),
    gym_website varchar(100),
    gym_description varchar(1000) NOT NULL);

CREATE TABLE referee (
    referee_id serial PRIMARY KEY,
    gym_id int, 
    style_id int, 
    first_name varchar,
    last_name varchar,
    CONSTRAINT FK_gym_id FOREIGN KEY (gym_id) REFERENCES gym(gym_id));

CREATE TABLE style (
	style_id serial PRIMARY KEY,
    style_name varchar(20) NOT NULL);

CREATE TABLE athlete (
	athlete_id serial PRIMARY KEY,
	first_name varchar(20) NOT NULL,
	last_name varchar(30) NOT NULL,
    username varchar(30) NOT NULL,
	birth_date date NOT NULL,
    password varchar(30) NOT NULL,
    email varchar(100) NOT NULL);

CREATE TABLE athlete_record (
    athlete_id int,
    wins int,
    losses int,
    draws int,
    CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id)
);

CREATE TABLE athlete_gym (
    athlete_id int,
    gym_id int,
    CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_gym_id FOREIGN KEY (gym_id) REFERENCES gym(gym_id),
    CONSTRAINT unique_athlete_gym UNIQUE (athlete_id, gym_id));
	
CREATE TABLE athlete_score (
    athlete_id serial,
    style_id int,
    score int,
    CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id));
	
CREATE TABLE bout (
	bout_id serial PRIMARY KEY,
	challenger_id int NOT NULL,
    acceptor_id int NOT NULL,
    referee_id int NOT NULL,
    style_id int NOT NULL,
    accepted boolean,
    completed boolean,
    cancelled boolean,
    points int,
    CONSTRAINT FK_referee_id FOREIGN KEY (referee_id) REFERENCES athlete(athlete_id),
	CONSTRAINT FK_challenger_id FOREIGN KEY (challenger_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id),
    CONSTRAINT FK_acceptor_id FOREIGN KEY (acceptor_id) REFERENCES athlete(athlete_id));

CREATE TABLE outcome (
    outcome_id serial PRIMARY KEY,
    bout_id int UNIQUE,
    winner_id int,
    loser_id int,
    is_draw boolean,
    CONSTRAINT FK_winner_id FOREIGN KEY (winner_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_loser_id FOREIGN KEY (loser_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_bout_id FOREIGN KEY (bout_id) REFERENCES bout(bout_id));


CREATE TABLE athlete_style (
	athlete_id int,
    style_id int,
	CONSTRAINT FK_athlete_id FOREIGN KEY (athlete_id) REFERENCES athlete(athlete_id),
    CONSTRAINT FK_style_id FOREIGN KEY (style_id) REFERENCES style(style_id),
    CONSTRAINT unique_athlete_style UNIQUE (athlete_id, style_id));
	
CREATE TABLE competition (
	competition_id serial PRIMARY KEY,
    referee_id int,
    bout_id int,
    outcome_id int,
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