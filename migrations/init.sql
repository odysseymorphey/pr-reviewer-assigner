CREATE TYPE pull_request_status AS ENUM ('OPEN', 'MERGED');

CREATE TABLE teams (
    team_name TEXT PRIMARY KEY
);

CREATE TABLE users (
    user_id   TEXT PRIMARY KEY,
    username  TEXT NOT NULL,
    team_name TEXT NOT NULL REFERENCES teams(team_name)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE
);

CREATE TABLE pull_requests (
    pull_request_id   TEXT PRIMARY KEY,
    pull_request_name TEXT NOT NULL,
    author_id         TEXT NOT NULL REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    status            pull_request_status NOT NULL,
    created_at        TEXT NOT NULL,
    merged_at         TEXT
);

CREATE TABLE pull_request_reviewers (
    pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id)
        ON UPDATE CASCADE
        ON DELETE CASCADE,
    reviewer_id     TEXT NOT NULL REFERENCES users(user_id)
        ON UPDATE CASCADE
        ON DELETE RESTRICT,
    PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX idx_pull_request_reviewers_reviewer ON pull_request_reviewers (reviewer_id);
CREATE UNIQUE INDEX ux_users_team_name_username ON users(team_name, username);