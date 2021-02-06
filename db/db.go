package db

import (
	"database/sql"
	"log"
	"strconv"
)

const queryCountTotalBadIps = "SELECT COUNT(1) FROM bips"

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

func (db *Fail2BanDB) CountTotalBadIps() (int, error) {
	stmt, err := db.sqliteDB.Prepare(queryCountTotalBadIps)
	defer db.mustCloseStatement(stmt)

	if err != nil {
		return -1, err
	}

	result := ""
	err = stmt.QueryRow().Scan(&result)

	if err != nil {
		return -1, err
	}

	return strconv.Atoi(result)
}

func (db *Fail2BanDB) mustCloseStatement(stmt *sql.Stmt) {
	err := stmt.Close()
	if err != nil {
		log.Fatal(err)
	}
}
