DELETE
  FROM edges
  WHERE from_rowid IN (
	SELECT rowid
	FROM vertices
	WHERE type=?
	AND id=?
 )
 AND to_rowid IN (
	SELECT rowid
	FROM vertices
	WHERE type=?
	AND id=?
 )
 AND key=?
