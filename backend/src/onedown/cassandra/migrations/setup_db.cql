CREATE TABLE IF NOT EXISTS users (
    ID uuid,
    Email text,
    PRIMARY KEY(ID, email)
);

CREATE TABLE IF NOT EXISTS puzzle_sessions (
    ID uuid,
    SessionData blob,
    PRIMARY KEY(ID)
);