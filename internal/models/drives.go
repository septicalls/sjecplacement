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
	Roles       int
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
	stmt := `SELECT d.id, d.title, d.company, d.description, d.date, d.published,
	(SELECT COUNT(r.id) FROM roles r WHERE r.drive_id = d.id) FROM drives d WHERE d.id = $1`

	d := &Drive{}

	err := m.DB.QueryRow(stmt, id).Scan(&d.ID, &d.Title, &d.Company, &d.Description, &d.Date, &d.Published, &d.Roles)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		}
		return nil, err
	}

	return d, nil
}

func (m *DriveModel) Latest() ([]*Drive, error) {
	stmt := `SELECT d.id, d.title, d.company, d.description, d.date, COUNT(r.id)
	FROM drives d LEFT JOIN roles r ON r.drive_id = d.id
	WHERE d.published = true
	GROUP BY d.id, d.title, d.company, d.description, d.date
	ORDER BY d.date DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	drives := []*Drive{}

	for rows.Next() {
		d := &Drive{}

		err = rows.Scan(&d.ID, &d.Title, &d.Company, &d.Description, &d.Date, &d.Roles)
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
