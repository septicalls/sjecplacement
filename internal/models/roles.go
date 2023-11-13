package models

import "database/sql"

type Role struct {
	ID            int
	Profile       string
	Description   string
	Qualification string
	// Nullable
	Cutoff           sql.NullString
	Location         sql.NullString
	Stipend          sql.NullInt32
	CTC              sql.NullFloat64
	ServiceAgreement sql.NullFloat64
	// Reference to Drive Model
	DriveID int
}

type RoleModel struct {
	DB *sql.DB
}

func (m *RoleModel) Insert(r *Role) (int, error) {
	stmt := `INSERT INTO roles (profile, description, qualification, cutoff, location, stipend, ctc,
	serviceagreement, drive_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, r.Profile, r.Description, r.Qualification, r.Cutoff, r.Location,
		r.Stipend, r.CTC, r.ServiceAgreement, r.DriveID).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (m *RoleModel) All(id int) ([]*Role, error) {
	stmt := `SELECT id, profile, description, qualification, cutoff, location, stiped, ctc,
	serviceagreement FROM roles WHERE id = $1 ORDER BY id ASC`

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	roles := []*Role{}

	for rows.Next() {
		r := &Role{}

		err := rows.Scan(
			&r.ID,
			&r.Profile,
			&r.Description,
			&r.Qualification,
			&r.Cutoff,
			&r.Location,
			&r.Stipend,
			&r.CTC,
			&r.ServiceAgreement,
		)
		if err != nil {
			return nil, err
		}

		r.DriveID = id

		roles = append(roles, r)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}
