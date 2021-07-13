DELETE
  FROM edges
  WHERE from_rowid IN (
	SELECT rowid
	FROM vertices
	WHERE type=$1
	AND id=$2
 )
 AND to_rowid IN (
	SELECT rowid
	FROM vertices
	WHERE type=$3
	AND id=$4
 )
 AND key=$5
