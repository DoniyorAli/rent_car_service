package postgres

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/brand"
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/car"
	"errors"
	"time"
)

func (psql Postgres) AddCar(id string, req *car.CreateCarRequest) error {
	Id := &brand.GetBrandByIDRequest{
		BrandId: req.BrandId,
	}

	_, err := psql.GetBrandById(Id)
	if err != nil {
		return errors.New("brand not found to AddCar")
	}

	_, err = psql.homeDB.Exec(`
	INSERT INTO "car" 
	(
		id,
		model,
		color,
		cartype,
		mileage,
		year,
		price,
		brand_id,
		created_at
	) VALUES (
		$1,
		$2,
		$3,
		$4,
		$5,
		$6,
		$7,
		$8,
		now()
	)`,
		id,
		req.Model,
		req.Color,
		req.CarType,
		req.Mileage,
		req.Year,
		req.Price,
		req.BrandId,
	)
	if err != nil {
		return err
	}
	return nil
}

// *=========================================================================
func (psql Postgres) GetCarById(id string) (*car.GetCarByIDResponse, error) {
	res := &car.GetCarByIDResponse{
		Brand: &car.GetCarByIDResponse_Brand{},
	}

	var deletedAt *time.Time
	var updatedAt, brandUpdatedAt *string
	err := psql.homeDB.QueryRow(`SELECT 
		c.id,
		c.model,
		c.color,
		c.cartype,
		c.mileage,
		c.year,
		c.price,
		c.created_at,
		c.updated_at,
		b.id,
		b.name,
		b.country,
		b.manufacturer,
		b.about_brand,
		b.created_at,
		b.updated_at
    FROM "car" AS c JOIN brand AS b ON c.brand_id = b.id WHERE c.id = $1`, id).Scan(
		&res.CarId,
		&res.Model,
		&res.Color,
		&res.CarType,
		&res.Mileage,
		&res.Year,
		&res.Price,
		&res.CreatedAt,
		&updatedAt,
		&res.Brand.BrandId,
		&res.Brand.Name,
		&res.Brand.Country,
		&res.Brand.Manufacturer,
		&res.Brand.AboutBrand,
		&res.Brand.CreatedAt,
		&brandUpdatedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if brandUpdatedAt != nil {
		res.Brand.UpdatedAt = *brandUpdatedAt
	}

	if deletedAt != nil {
		return res, errors.New("car not fount")
	}

	return res, nil
}

// *=========================================================================
func (psql Postgres) GetCarList(limit, offset int, search string) (*car.GetCarListResponse, error) {
	resp := &car.GetCarListResponse{
		Cars: make([]*car.Car, 0),
	}

	rows, err := psql.homeDB.Queryx(`SELECT 
		id,
		model,
		color,
		cartype,
		mileage,
		year,
		price,
		brand_id,
		created_at,
		updated_at
		FROM "car" WHERE deleted_at IS NULL AND (model ILIKE '%' || $1 || '%')
		LIMIT $2
		OFFSET $3
	`,
		search,
		limit,
		offset,
	)
	if err != nil {
		return resp, err
	}

	for rows.Next() {
		var c = &car.Car{}
		var updatedAt *string

		err := rows.Scan(
			&c.CarId,
			&c.Model,
			&c.Color,
			&c.CarType,
			&c.Mileage,
			&c.Year,
			&c.Price,
			&c.BrandId,
			&c.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return resp, err
		}

		if updatedAt != nil {
			c.UpdatedAt = *updatedAt
		}
		resp.Cars = append(resp.Cars, c)

	}
	return resp, err
}

// *=========================================================================
func (psql Postgres) UpdateCar(id string, box *car.UpdateCarRequest) error {

	res, err := psql.homeDB.NamedExec(`
	UPDATE "car"  
		SET model=:m, color=:c, cartype=:ct, mileage=:ml, year=:y, price=:p, brand_id=:b, updated_at=now() 
			WHERE deleted_at IS NULL AND id=:id`, 
	map[string]interface{}{
		"id": box.Id,
		"m":  box.Model,
		"c":  box.Color,
		"ml": box.Mileage,
		"y": box.Year,
		"p": box.Price,
		"b": box.BrandId,
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
	return errors.New("car not found")
}

// *=========================================================================
func (psql Postgres) DeleteCar(id string) error {
	res, err := psql.homeDB.Exec(`
	UPDATE "car" 
		SET deleted_at=now() 
			WHERE id=$1 AND deleted_at IS NULL`, id)
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
	return errors.New("car not found")
}