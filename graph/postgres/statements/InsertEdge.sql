INSERT INTO edges(from_rowid, to_rowid, key, position, meta)
SELECT a.rowid,
       b.rowid,
       $1,
       $2,
       $3
  FROM vertices a,
       vertices b
 WHERE (a.type=$4 AND a.id=$5 AND b.type=$6 AND b.id=$7)
