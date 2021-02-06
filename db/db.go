package db

import (
	"database/sql"
	"log"
)

const queryBadIpsPerJail = "SELECT jail, COUNT(1) FROM bips GROUP BY jail"

type Fail2BanDB struct {
	DatabasePath string
	sqliteDB *sql.DB
}

func MustConnectToDb(databasePath string) *Fail2BanDB {
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	return &Fail2BanDB{
		DatabasePath: databasePath,
		sqliteDB:     db,
	}
}

func (db *Fail2BanDB) CountBadIpsPerJail() (map[string]int, error) {
	stmt, err := db.sqliteDB.Prepare(queryBadIpsPerJail)
	defer db.mustCloseStatement(stmt)

	if err != nil {
		return nil, err
	}

	jailNameToCountMap := map[string]int{}
	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	if rows == nil {
		return jailNameToCountMap, nil
	}

	for rows.Next() {
		if rows.Err() != nil {
			return nil, err
		}
		jailName := ""
		count := 0
		err = rows.Scan(&jailName, &count)
		if err != nil {
			return nil, err
		}

		jailNameToCountMap[jailName] = count
	}
	return jailNameToCountMap, nil
}

func (db *Fail2BanDB) mustCloseStatement(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}
