package store

import (
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/crema-labs/gitgo/pkg/model"
	_ "github.com/mattn/go-sqlite3"
)

var ErrGrantNotFound = errors.New("grant not found")

type Store interface {
	GetGrant(id string) (*model.Grant, error)
	UpdateGrantStatus(id, status string) error
	InsertGrant(grant *model.Grant) error
	GetAllGrants() ([]*model.Grant, error)
	Close() error
	DB() *sql.DB
}

type SQLiteStore struct {
	db *sql.DB
}

func (s *SQLiteStore) DB() *sql.DB {
	return s.db
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := createTable(db); err != nil {
		return nil, err
	}

	return &SQLiteStore{db: db}, nil
}

func createTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS grants (
			grantid TEXT PRIMARY KEY,
			grant_amount TEXT,
			status TEXT,
			contributions TEXT
		)
	`)
	return err
}

func (s *SQLiteStore) GetGrant(id string) (*model.Grant, error) {
	var grant model.Grant
	var contributionsJSON string

	err := s.db.QueryRow("SELECT grantid, grant_amount, status, contributions FROM grants WHERE grantid = ?", id).Scan(
		&grant.GrantID, &grant.GrantAmount, &grant.Status, &contributionsJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrGrantNotFound
		}
		return nil, err
	}

	err = json.Unmarshal([]byte(contributionsJSON), &grant.Contributions)
	if err != nil {
		return nil, err
	}

	return &grant, nil
}

func (s *SQLiteStore) UpdateGrantStatus(id, status string) error {
	tx, err := s.db.Exec("UPDATE grants SET status = ? WHERE grantid = ?", status, id)
	if err != nil {
		return err
	}
	re, err := tx.RowsAffected()
	if err != nil {
		return err
	}

	if re == 0 {
		return ErrGrantNotFound
	}

	return nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

func (s *SQLiteStore) InsertGrant(grant *model.Grant) error {
	contributionsJSON, err := json.Marshal(grant.Contributions)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		"INSERT INTO grants (grantid, grant_amount, status, contributions) VALUES (?, ?, ?, ?)",
		grant.GrantID, grant.GrantAmount, grant.Status, string(contributionsJSON),
	)
	return err
}

func (s *SQLiteStore) GetAllGrants() ([]*model.Grant, error) {
	rows, err := s.db.Query("SELECT grantid, grant_amount, status, contributions FROM grants")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grants []*model.Grant
	for rows.Next() {
		var grant model.Grant
		var contributionsJSON string

		err := rows.Scan(&grant.GrantID, &grant.GrantAmount, &grant.Status, &contributionsJSON)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(contributionsJSON), &grant.Contributions)
		if err != nil {
			return nil, err
		}

		grants = append(grants, &grant)
	}

	return grants, nil
}
