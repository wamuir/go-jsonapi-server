CREATE TABLE IF NOT EXISTS vertices (
    rowid INTEGER PRIMARY KEY,
    type TEXT NOT NULL,
    id TEXT NOT NULL,
    attributes TEXT,
    meta TEXT
)
