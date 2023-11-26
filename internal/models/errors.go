package models

import "errors"

var (
	ErrNoRecord           = errors.New("models: no matching record found")
	ErrPublish            = errors.New("drive: already published")
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	ErrDuplicateEmail     = errors.New("models: duplicate email")
)
