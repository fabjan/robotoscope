// Package database manages the robot database
package database

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v4/stdlib" // for sql.Open("pgx", ...)

	"github.com/fabjan/robotoscope/core"
)

// PgStore tracks robots in Postgres
type PgStore struct {
	db        *sql.DB
	tableName string
}

func (ts *PgStore) createTableSQL() string {
	sql := `
CREATE TABLE IF NOT EXISTS %s (
	user_agent text not null,
	seen numeric,
	unique(user_agent)
);
`
	return fmt.Sprintf(sql, ts.tableName)
}

func (ts *PgStore) selectSQL(limit int) string {
	return fmt.Sprintf("SELECT user_agent, seen from %s limit %d", ts.tableName, limit)
}

func (ts *PgStore) insertSQL() string {
	sql := `
INSERT INTO %s (user_agent, seen) VALUES ($1, 1)
ON CONFLICT (user_agent) DO UPDATE SET seen = %s.seen + 1
`
	return fmt.Sprintf(sql, ts.tableName, ts.tableName)
}

// Count increases the seen count for the given bot.
func (ts *PgStore) Count(name string) error {
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

// List returns a list showing how many times each robot has been seen.
func (ts *PgStore) List() ([]core.RobotInfo, error) {
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

// OpenPg opens a connection to the Postgres database with the given URL.
func OpenPg(rawURL string) (*sql.DB, error) {
	return sql.Open("pgx", rawURL)
}

// GetPgStore creates a PgStore backed by the given table and DB.
// The table is created if it does not exist.
func GetPgStore(db *sql.DB, name string) (*PgStore, error) {
	s := PgStore{
		tableName: name,
		db:        db,
	}

	_, err := db.Exec(s.createTableSQL())
	if err != nil {
		return nil, err
	}

	return &s, nil
}
