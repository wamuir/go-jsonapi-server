CREATE UNIQUE INDEX IF NOT EXISTS edges_idx1
    ON edges (from_rowid, key, to_rowid)
