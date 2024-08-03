CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    date CHAR(8),
    title VARCHAR(256),
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS date_index on scheduler (date);