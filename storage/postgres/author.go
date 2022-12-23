package postgres

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/blogpost"
	"database/sql"
	"errors"
)

func (stg Postgres) AddAuthor(id string, box *blogpost.CreateAuthorRequest) error {
	_, err := stg.homeDB.Exec(`INSERT INTO author 
	(
		id,
		firstname,
		lastname,
		fullname,
		middlename
	) VALUES (
		$1,
		'',
		'',
		$2,
		$3
	)`,
		id,
		box.Fullname,
		box.Middlename,
	)
	if err != nil {
		return err
	}
	return nil
}

// *=========================================================================
func (stg Postgres) GetAuthorById(id string) (*blogpost.Author, error) {
	author := &blogpost.Author{}
	var updatedAt *string
	var tempMiddlename *string
	var deletedAt sql.NullString
	err := stg.homeDB.QueryRow(`SELECT 
		id,
		fullname,
		middlename,
		created_at,
		updated_at,
		deleted_at
    FROM author WHERE id = $1 AND deleted_at IS NULL`, id).Scan(
		&author.Id,
		&author.Fullname,
		&tempMiddlename,
		&author.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return author, err
	}

	if updatedAt != nil {
		author.UpdatedAt = *updatedAt
	}

	if tempMiddlename != nil {
		author.Middlename = *tempMiddlename
	}

	author.DeletedAt = deletedAt.String

	return author, nil
}

// *=========================================================================
func (stg Postgres) GetAuthorList(offset, limit int, search string) (*blogpost.GetAuthorListResponse, error) {
	res := &blogpost.GetAuthorListResponse{
		Authors: make([]*blogpost.Author, 0),
	}
	var tempMiddlename *string
	var updatedAt *string
	var deletedAt *string
	rows, err := stg.homeDB.Queryx(`SELECT 
		id,
		fullname,
		middlename,
		created_at,
		updated_at,
		deleted_at 
		FROM author
		WHERE ((fullname ILIKE '%' || $1 || '%') or (middlename ILIKE '%' || $1 || '%')) AND deleted_at IS NULL
		LIMIT $2
		OFFSET $3
	`,
		search,
		limit,
		offset,
	)
	if err != nil {
		return res, err
	}

	for rows.Next() {
		var author blogpost.Author
		err := rows.Scan(
			&author.Id,
			&author.Fullname,
			&tempMiddlename,
			&author.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return res, err
		}

		if updatedAt != nil {
			author.UpdatedAt = *updatedAt
		}
		if deletedAt != nil {
			author.DeletedAt = *deletedAt
		}

		if tempMiddlename != nil {
			author.Middlename = *tempMiddlename
		}
		res.Authors = append(res.Authors, &author)

	}
	return res, err
}

// *=========================================================================
func (stg Postgres) UpdateAuthor(box *blogpost.UpdateAuthorRequest) error {

	res, err := stg.homeDB.NamedExec("UPDATE author  SET fullname=:f, middlename=:m, updated_at=now() WHERE deleted_at IS NULL AND id=:id", map[string]interface{}{
		"id": box.Id,
		"f":  box.Fullname,
		"m":  box.Middlename,
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
	return errors.New("author not found")
}

// *=========================================================================
func (stg Postgres) DeleteAuthor(id string) error {
	res, err := stg.homeDB.Exec("UPDATE author SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
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
	return errors.New("author not found")
}
