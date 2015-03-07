DROP TABLE IF EXISTS users;

CREATE TABLE users (
	username char(20) NOT NULL,
	password char(255) NOT NULL,
	PRIMARY KEY (username)
);

DROP TABLE IF EXISTS notes;

CREATE TABLE notes (
	username char(20) NOT NULL,
	note TEXT NOT NULL,
	PRIMARY KEY (username)
);
