// Package database manages the robot database
package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // for sql.Open("pgx", ...)

	"github.com/fabjan/robotoscope/core"
)

type tableStore struct {
	db        *sql.DB
	tableName string
}

var createRobotsSQL string = `
CREATE TABLE IF NOT EXISTS robots (
	user_agent text not null,
	seen numeric,
    unique(user_agent)
);
`

var createCheatersSQL string = `
CREATE TABLE IF NOT EXISTS cheaters (
	user_agent text not null,
	seen numeric,
    unique(user_agent)
);
`

func (ts *tableStore) selectSQL(limit int) string {
	return fmt.Sprintf("SELECT user_agent, seen from %s limit %d", ts.tableName, limit)
}

func (ts *tableStore) insertSQL() string {
	sql := `
INSERT INTO %s (user_agent, seen) VALUES ($1, 1)
ON CONFLICT (user_agent) DO UPDATE SET seen = %s.seen + 1
`
	return fmt.Sprintf(sql, ts.tableName, ts.tableName)
}

func (ts *tableStore) List() ([]core.RobotInfo, error) {
	rows, err := ts.db.Query(ts.selectSQL(640)) // 640 rows ought to be enough for anyone
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	info := []core.RobotInfo{}
	ri := core.RobotInfo{}
	for rows.Next() {
		err := rows.Scan(&ri.UserAgent, &ri.Seen)
		if err != nil {
			return nil, err
		}
		info = append(info, ri)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return info, nil
}

func (ts *tableStore) Count(name string) error {
	res, err := ts.db.Exec(ts.insertSQL(), name)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("no rows updated")
	}
	return nil
}

// PgStores maintains robot stores robots and cheaters.
type PgStores struct {
	Robots   *tableStore
	Cheaters *tableStore
	db       *sql.DB
}

// GetPostgresStores connects to the given database, initializes the tables, and
// returns the connection.
func GetPostgresStores(dbURL string) (*PgStores, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createRobotsSQL)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createCheatersSQL)
	if err != nil {
		return nil, err
	}

	stores := &PgStores{
		db: db,
		Robots: &tableStore{
			tableName: "robots",
			db:        db,
		},
		Cheaters: &tableStore{
			tableName: "cheaters",
			db:        db,
		},
	}

	return stores, nil
}

// Close closes the database connection.
func (s *PgStores) Close() {
	s.db.Close()
}
