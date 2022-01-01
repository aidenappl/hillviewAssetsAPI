package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/assetsAPI/db"
)

func CreateNewAsset(db db.Queryable, name string, identifier string, category int, imageURL *string, description *string) (int64, error) {

	cols := []string{"name", "identifier", "category", "description"}
	var vals []interface{}

	vals = append(vals, name, identifier, category, description)

	if imageURL != nil {
		cols = append(cols, "image_url")
		vals = append(vals, imageURL)
	}

	q := sq.Insert("assets").
		Columns(cols...).
		Values(vals...)

	query, args, err := q.ToSql()

	if err != nil {
		return 0, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Exec(query, args...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute sql query: %w", err)
	}

	lid, err := rows.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("failed to get the last inserted ID: %w", err)
	}

	return lid, nil
}
