package models

import (
	"database/sql"
	"errors"
	"time"
)

type Drive struct {
	ID          int
	Title       string
	Company     string
	Description string
	Date        time.Time
	Published   bool
}

type DriveModel struct {
	DB *sql.DB
}

func (m *DriveModel) Insert(title string, company string, description string, date time.Time) (int, error) {
	stmt := `INSERT INTO "drives" (title, company, description, date) VALUES
	($1, $2, $3, $4) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, company, description, date).Scan(&id)
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *DriveModel) Get(id int) (*Drive, error) {
	stmt := `SELECT id, title, company, description, date, published FROM "drives"
	WHERE id = $1`

	d := &Drive{}

	err := m.DB.QueryRow(stmt, id).Scan(&d.ID, &d.Title, &d.Company, &d.Description, &d.Date, &d.Published)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return d, nil
}

func (m *DriveModel) Latest() ([]*Drive, error) {
	stmt := `SELECT id, title, company, description, date FROM "drives"
	WHERE published=true ORDER BY date DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	drives := []*Drive{}

	for rows.Next() {
		d := &Drive{}

		err = rows.Scan(&d.ID, &d.Title, &d.Company, &d.Description, &d.Date)
		if err != nil {
			return nil, err
		}

		drives = append(drives, d)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return drives, nil
}
