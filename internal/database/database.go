package database

import (
    "database/sql"
    "fmt"
    "sync"
	"errors"
    _ "github.com/glebarez/go-sqlite"
)

var (
    db   *sql.DB
    once sync.Once
)

type Stats struct {
	WPM float64
}

// ConnectToDatabase initializes the connection if it hasn't been initialized
func connectToDatabase() (*sql.DB, error) {
    var err error
    once.Do(func() {
        db, err = sql.Open("sqlite", "./typtea.db")
        if err != nil {
            return
        }

        _, err = createTable(db)
        if err == nil {
            fmt.Println("Connected to the SQLite database successfully.")
        }
    })

    return db, err
}

// CreateTable ensures the stats table exists
func createTable(db *sql.DB) (sql.Result, error) {
    sqlStmt := `CREATE TABLE IF NOT EXISTS stats (
        id INTEGER PRIMARY KEY,
        wpm REAL NOT NULL
    );`
    return db.Exec(sqlStmt)
}

// Insert adds a new stats entry
func InsertStats(stats *Stats) (int64, error) {
    db, err := connectToDatabase()
    if err != nil {
        return 0, err
    }

    sqlStmt := `INSERT INTO stats (wpm) VALUES (?);`
    result, err := db.Exec(sqlStmt, stats.WPM)
    if err != nil {
        return 0, err
    }

    return result.LastInsertId()
}

// Retrives the max wpm
func MaxWPM() (float64, error) {
    db, err := connectToDatabase()
    if err != nil {
        return 0, err
    }

    sqlStmt := `SELECT MAX(wpm) FROM stats;`
    var maxWPM sql.NullFloat64

    err = db.QueryRow(sqlStmt).Scan(&maxWPM)
    if err != nil {
        return 0, err
    }

    if !maxWPM.Valid {
        return 0, errors.New("no wpm records found")
    }

    return maxWPM.Float64, nil
}
