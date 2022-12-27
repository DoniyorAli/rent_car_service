package postgres

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/brand"
	"database/sql"
	"errors"
)

// *=========================================================================
func (psql Postgres) AddNewBrand(id string, req *brand.CreateBrandRequest) (res *brand.Brand, err error) {
	_, err = psql.homeDB.Exec(`INSERT INTO "brand" 
	(
		id,
		name,
		country,
		manufacturer,
		about_brand,
		created_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		now()
	)`,
		id,
		req.Name,
		req.Country,
		req.Manufacturer,
		req.AboutBrand,
	)
	if err != nil {
		return nil, errors.New("error in AddNewBrand")
	}

	Id := &brand.GetBrandByIDRequest{
		BrandId: id,
	}
	if err != nil {
		return nil, err
	}

	res, err = psql.GetBrandById(Id)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// *=========================================================================
func (psql Postgres) GetBrandById(req *brand.GetBrandByIDRequest) (*brand.Brand, error) {
	box := &brand.Brand{}

	var updated_at sql.NullString
	var deleted_at sql.NullString

	rows := psql.homeDB.QueryRow(`SELECT 
		id,
		name,
		country,
		manufacturer,
		about_brand,
		created_at,
		updated_at,
		deleted_at
    FROM "brand" WHERE id = $1`, req.BrandId)
	err := rows.Scan(
		&box.BrandId,
		&box.Name,
		&box.Country,
		&box.Manufacturer,
		&box.AboutBrand,
		&box.CreatedAt,
		&updated_at,
		&deleted_at,
	)

	if updated_at.Valid {
		box.UpdatedAt = updated_at.String
	}

	if deleted_at.Valid {
		return nil, errors.New("brand not found")
	}

	if err != nil {
		return nil, err
	}
	return box, err
}

// *=========================================================================
func (psql Postgres) GetBrandList(req *brand.GetBrandListRequest) (*brand.GetBrandListResponse, error) {
	resp := &brand.GetBrandListResponse{
		Brands: []*brand.Brand{},
	}
	rows, err := psql.homeDB.Queryx(`SELECT
		id,
		name,
		country,
		manufacturer,
		about_brand,
		created_at,
		updated_at
	FROM "brand" WHERE deleted_at IS NULL AND (name || about_brand ILIKE '%' || $1 || '%')
		LIMIT $2
		OFFSET $3
	`, req.Search, int(req.Limit), int(req.Offset))

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var b brand.Brand
		var updated_at sql.NullString

		err := rows.Scan(
			&b.BrandId,
			&b.Name,
			&b.Country,
			&b.Manufacturer,
			&b.AboutBrand,
			&b.CreatedAt,
			&updated_at,
		)
		if err != nil {
			return nil, err
		}

		if updated_at.Valid {
			b.UpdatedAt = updated_at.String
		}

		resp.Brands = append(resp.Brands, &b)
	}
	return resp, nil
}

// *=========================================================================
func (psql Postgres) UpdateBrand(id string, req *brand.UpdateBrandRequest) error {
	res, err := psql.homeDB.NamedExec(`
	UPDATE "brand"  
		SET name=:n, country=:c, manufacturer=:m, about_brand=:a, updated_at=now() 
			WHERE deleted_at IS NULL AND id=:id`, map[string]interface{}{
		"id": id,
		"n": req.Name,
		"c": req.Country,
		"m": req.Manufacturer,
		"a": req.AboutBrand,

	})
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect > 0 {
		return nil
	}
	return errors.New("brand not found")
}

// *=========================================================================
func (psql Postgres) DeleteBrand(req *brand.DeleteBrandRequest) error {
	res, err := psql.homeDB.Exec(`
	UPDATE "brand" 
		SET deleted_at=now() 
			WHERE id=$1 AND deleted_at IS NULL`, req.BrandId)
	if err != nil {
		return err
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if affect > 0 {
		return nil
	}
	return errors.New("brand had been already deleted")
}