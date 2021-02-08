package db

import (
	"database/sql"
	"log"
	"os"
)

const queryBadIpsPerJail = "SELECT j.name, (SELECT COUNT(1) FROM bips b WHERE j.name = b.jail) FROM jails j"
const queryBannedIpsPerJail = "SELECT j.name, (SELECT COUNT(1) FROM bans b WHERE j.name = b.jail) FROM jails j"

type Fail2BanDB struct {
	DatabasePath string
	sqliteDB     *sql.DB
}

func MustConnectToDb(databasePath string) *Fail2BanDB {
	if _, err := os.Stat(databasePath); os.IsNotExist(err) {
		log.Fatalf("database path does not exist: %v", err)
	}
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		log.Fatal(err)
	}
	return &Fail2BanDB{
		DatabasePath: databasePath,
		sqliteDB:     db,
	}
}

func (db *Fail2BanDB) CountBannedIpsPerJail() (map[string]int, error) {
	return db.RunJailNameToCountQuery(queryBannedIpsPerJail)
}

func (db *Fail2BanDB) CountBadIpsPerJail() (map[string]int, error) {
	return db.RunJailNameToCountQuery(queryBadIpsPerJail)
}

func (db *Fail2BanDB) RunJailNameToCountQuery(query string) (map[string]int, error) {
	stmt, err := db.sqliteDB.Prepare(query)
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
	if stmt != nil {
		err := stmt.Close()
		if err != nil {
			log.Fatal(err)
		}
	}
}
