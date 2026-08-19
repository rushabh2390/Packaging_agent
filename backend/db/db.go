package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) error {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		return fmt.Errorf("open db error: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("ping db error: %w", err)
	}

	return nil
}

type BoxStyle struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func GetStylesLimit(limit int) ([]BoxStyle, error) {
	if DB == nil {
		return nil, fmt.Errorf("database connection is not initialized")
	}

	if limit <= 0 {
		limit = 3
	}

	rows, err := DB.Query("SELECT code, name, category, description FROM box_styles LIMIT ?", limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var styles []BoxStyle
	for rows.Next() {
		var s BoxStyle
		if err := rows.Scan(&s.Code, &s.Name, &s.Category, &s.Description); err != nil {
			return nil, err
		}
		styles = append(styles, s)
	}
	return styles, nil
}
