package query

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/hillview.tv/assetsAPI/db"
	"github.com/hillview.tv/assetsAPI/structs"
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

func CreateNewAssetInfo(db db.Queryable, id int64, serialNumber *string, manufacturer *string, model *string, notes *string) error {

	cols := []string{"asset_id"}
	var vals []interface{}

	vals = append(vals, id)

	if serialNumber != nil {
		cols = append(cols, "serial_number")
		vals = append(vals, serialNumber)
	}

	if manufacturer != nil {
		cols = append(cols, "manufacturer")
		vals = append(vals, manufacturer)
	}

	if model != nil {
		cols = append(cols, "model")
		vals = append(vals, model)
	}

	if notes != nil {
		cols = append(cols, "notes")
		vals = append(vals, notes)
	}

	q := sq.Insert("asset_metadata").
		Columns(cols...).
		Values(vals...)

	query, args, err := q.ToSql()

	if err != nil {
		return fmt.Errorf("failed to build sql query: %w", err)
	}

	_, err = db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to execute sql query: %w", err)
	}

	return nil
}

func ReadAsset(db db.Queryable, id *int, tag *string) (*structs.Asset, error) {

	if id == nil && tag == nil {
		return nil, fmt.Errorf("must specify either id or tag")
	}

	q := sq.Select(
		"assets.id",
		"assets.name",
		"assets.identifier",
		"assets.image_url",
		"assets.description",

		"asset_statuses.id",
		"asset_statuses.name",
		"asset_statuses.short_name",

		"asset_categories.id",
		"asset_categories.name",
		"asset_categories.short_name",

		"asset_metadata.serial_number",
		"asset_metadata.manufacturer",
		"asset_metadata.model",
		"asset_metadata.notes",
	).
		From("assets").
		Join("asset_statuses ON assets.status = asset_statuses.id").
		Join("asset_categories ON assets.category = asset_categories.id").
		LeftJoin("asset_metadata ON assets.id = asset_metadata.asset_id")

	if id != nil {
		q = q.Where(sq.Eq{"assets.id": id})
	}

	if tag != nil {
		q = q.Where(sq.Eq{"assets.identifier": tag})
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

	asset := structs.Asset{}
	assetStatus := structs.GeneralNSN{}
	assetCategory := structs.GeneralNSN{}
	assetMetadata := structs.AssetMetadata{}

	err = rows.Scan(
		&asset.ID,
		&asset.Name,
		&asset.Identifier,
		&asset.ImageURL,
		&asset.Description,

		&assetStatus.ID,
		&assetStatus.Name,
		&assetStatus.ShortName,

		&assetCategory.ID,
		&assetCategory.Name,
		&assetCategory.ShortName,

		&assetMetadata.SerialNumber,
		&assetMetadata.Manufacturer,
		&assetMetadata.Model,
		&assetMetadata.Notes,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	asset.Status = &assetStatus
	asset.Category = &assetCategory
	asset.Metadata = &assetMetadata

	checkout, err := ReadActiveCheckouts(db, int(asset.ID))
	if err != nil {
		return nil, fmt.Errorf("failed to read active checkouts: %w", err)
	}

	if checkout != nil {
		user, err := ReadUser(db, checkout.AssociatedUser, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to read active checkouts: %w", err)
		}
		checkout.User = user
	}

	asset.ActiveTab = checkout

	return &asset, nil
}
