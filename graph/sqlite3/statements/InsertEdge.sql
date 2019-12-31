INSERT INTO edges(from_rowid, to_rowid, key, position, meta)
SELECT a.rowid,
       b.rowid,
       ?,
       ?,
       ?
  FROM vertices a,
       vertices b
 WHERE (a.type=? AND a.id=? AND b.type=? AND b.id=?)
