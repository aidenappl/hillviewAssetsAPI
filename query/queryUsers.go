package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/structs"
)

func CreateNewUser(db db.Queryable, name string, email string, tag string, photoURL string) error {
	query, args, err := sq.Insert("users").
		Columns("name", "email", "identifier", "profile_image_url").
		Values(name, email, tag, photoURL).
		ToSql()
	if err != nil {
		return fmt.Errorf("error building query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil

}

func ReadUser(db db.Queryable, id *int, tag *string) (*structs.User, error) {

	if id == nil && tag == nil {
		return nil, fmt.Errorf("must specify either id or tag")
	}

	q := sq.Select(
		"users.id",
		"users.name",
		"users.email",
		"users.identifier",
		"users.profile_image_url",
		"users.inserted_at",
	).
		From("users")

	if id != nil {
		q = q.Where(sq.Eq{"users.id": *id})
	}

	if tag != nil {
		q = q.Where(sq.Eq{"users.identifier": *tag})
	}

	query, args, err := q.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build sql query: %w", err)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute sql query: %w", err)
	}

	if !rows.Next() {
		return nil, nil
	}

	defer rows.Close()

	var user structs.User

	err = rows.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Identifier,
		&user.ProfileImageURL,
		&user.InsertedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan user row: %w", err)
	}

	return &user, nil

}
