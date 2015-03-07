DROP TABLE IF EXISTS users;

CREATE TABLE users (
	username char(20),
	password char(50)
);

DROP TABLE IF EXISTS notes;

CREATE TABLE notes (
	username char(20),
	note char(500)
);
