package validator

import (
	"database/sql"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFieldError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChar(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func ValidDate(date string) (time.Time, bool) {
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil || !parsedDate.After(time.Now()) {
		return time.Time{}, false
	}

	return parsedDate, true
}

func ValidInt(integer string) (sql.NullInt32, bool) {
	if strings.TrimSpace(integer) == "" {
		return sql.NullInt32{Valid: false}, true
	}

	parsedInt, err := strconv.Atoi(integer)
	if err != nil || parsedInt < 0 {
		return sql.NullInt32{}, false
	}

	return sql.NullInt32{Int32: int32(parsedInt), Valid: true}, true
}

func ValidFloat(float string) (sql.NullFloat64, bool) {
	if strings.TrimSpace(float) == "" {
		return sql.NullFloat64{Valid: false}, true
	}

	parsedFloat, err := strconv.ParseFloat(float, 64)
	if err != nil || parsedFloat < 0.0 {
		return sql.NullFloat64{}, false
	}

	return sql.NullFloat64{Float64: parsedFloat, Valid: true}, true
}

func ValidString(str string, n int) (sql.NullString, bool) {
	if !MaxChar(str, n) {
		return sql.NullString{}, false
	}

	if !NotBlank(str) {
		return sql.NullString{Valid: false}, true
	}

	return sql.NullString{String: str, Valid: true}, true
}
