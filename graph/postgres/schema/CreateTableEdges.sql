CREATE TABLE IF NOT EXISTS edges (
    rowid SERIAL PRIMARY KEY,
    from_rowid INTEGER NOT NULL
        REFERENCES vertices(rowid)
	ON DELETE CASCADE,
    to_rowid INTEGER NOT NULL
        REFERENCES vertices(rowid)
	ON DELETE CASCADE,
    key TEXT NOT NULL,
    position INTEGER,
    meta TEXT
)
