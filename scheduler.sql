CREATE TABLE scheduler (
    id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
    date CHAR(8),
    title VARCHAR(256),
    comment TEXT,
    repeat VARCHAR(128)
);

CREATE INDEX date_index on scheduler (date);