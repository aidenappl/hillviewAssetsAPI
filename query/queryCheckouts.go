package query

import (
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/structs"
)

type ReadAssetCheckoutsRequest struct {
	AssetID int
	Limit   int
}

func ReadAssetCheckoutHistory(db db.Queryable, req ReadAssetCheckoutsRequest) ([]*structs.AssetCheckout, error) {
	query, args, err := sq.Select(
		"asset_checkouts.id",
		"asset_checkouts.asset_id",
		"asset_checkouts.associated_user",
		"asset_checkouts.checkout_notes",
		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",

		"checkout_statuses.id",
		"checkout_statuses.name",
		"checkout_statuses.short_name",

		"users.id",
		"users.name",
		"users.email",
		"users.identifier",
		"users.profile_image_url",
		"users.inserted_at",
	).
		From("asset_checkouts").
		Where(sq.Eq{"asset_checkouts.asset_id": req.AssetID}).
		Join("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
		Join("users ON asset_checkouts.associated_user = users.id").
		Limit(uint64(req.Limit)).
		ToSql()
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

	var checkouts []*structs.AssetCheckout

	for rows.Next() {
		var checkout structs.AssetCheckout
		var checkout_status structs.GeneralNSN
		var user structs.User

		err = rows.Scan(
			&checkout.ID,
			&checkout.AssetID,
			&checkout.AssociatedUser,
			&checkout.CheckoutNotes,
			&checkout.TimeOut,
			&checkout.TimeIn,
			&checkout.ExpectedIn,

			&checkout_status.ID,
			&checkout_status.Name,
			&checkout_status.ShortName,

			&user.ID,
			&user.Name,
			&user.Email,
			&user.Identifier,
			&user.ProfileImageURL,
			&user.InsertedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		checkout.CheckoutStatus = &checkout_status
		checkout.User = &user

		checkouts = append(checkouts, &checkout)
	}

	return checkouts, nil

}

func CheckinAsset(db db.Queryable, checkoutID int, notes *string) error {
	query, args, err := sq.Update("asset_checkouts").
		Set("checkout_status", 2).
		Set("time_in", time.Now()).
		Set("checkout_notes", notes).
		Where(sq.Eq{"id": checkoutID}).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}

func CheckoutAsset(db db.Queryable, userID int, assetID int, duration time.Time, offsite bool, notes *string) error {

	query, args, err := sq.Insert("asset_checkouts").
		Columns(
			"asset_id",
			"associated_user",
			"checkout_status",
			"checkout_notes",
			"expected_in",
			"offsite",
		).
		Values(assetID, userID, 1, notes, duration, offsite).
		ToSql()
	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}

func ReadActiveCheckouts(db db.Queryable, id int) (*structs.AssetCheckout, error) {

	query, args, err := sq.Select(
		"asset_checkouts.id",
		"asset_checkouts.asset_id",
		"asset_checkouts.associated_user",
		"asset_checkouts.checkout_notes",
		"asset_checkouts.time_out",
		"asset_checkouts.time_in",
		"asset_checkouts.expected_in",

		"checkout_statuses.id",
		"checkout_statuses.name",
		"checkout_statuses.short_name",
	).
		From("asset_checkouts").
		Where(sq.Eq{"asset_checkouts.asset_id": id}).
		Where(sq.Eq{"asset_checkouts.checkout_status": 1}).
		Join("checkout_statuses ON asset_checkouts.checkout_status = checkout_statuses.id").
		ToSql()
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

	var checkout structs.AssetCheckout
	var checkout_status structs.GeneralNSN

	err = rows.Scan(
		&checkout.ID,
		&checkout.AssetID,
		&checkout.AssociatedUser,
		&checkout.CheckoutNotes,
		&checkout.TimeOut,
		&checkout.TimeIn,
		&checkout.ExpectedIn,

		&checkout_status.ID,
		&checkout_status.Name,
		&checkout_status.ShortName,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	checkout.CheckoutStatus = &checkout_status

	return &checkout, nil
}
