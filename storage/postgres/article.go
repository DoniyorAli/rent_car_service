package postgres

import (
	"MyProjects/RentCar_gRPC/rent_car_service/protogen/blogpost"
	"errors"
	"time"
)

// *=========================================================================
func (stg Postgres) AddNewArticle(id string, box *blogpost.CreateArticleRequest) error {
	var err error
	_, err = stg.GetAuthorById(box.AuthorId)
	if err != nil {
		return err
	}

	if box.Content == nil {
		box.Content = &blogpost.Content{}
	}

	_, err = stg.homeDB.Exec(`INSERT INTO article 
	(
		id,
		title,
		body,
		author_id
	) VALUES (
		$1,
		$2,
		$3,
		$4
	)`,
		id,
		box.Content.Title,
		box.Content.Body,
		box.AuthorId,
	)
	if err != nil {
		return err
	}
	return nil
}

// *=========================================================================
func (stg Postgres) GetArticleById(id string) (*blogpost.GetArticleByIDResponse, error) {
	res := &blogpost.GetArticleByIDResponse{
		Content: &blogpost.Content{},
		Author:  &blogpost.GetArticleByIDResponse_Author{},
	}
	var deletedAt *time.Time

	var updatedAt, authorUpdatedAt *string

	var tempMiddlename *string

	err := stg.homeDB.QueryRow(`SELECT 
		ar.id,
		ar.title,
		ar.body,
		ar.created_at,
		ar.updated_at,
		au.id,
		au.fullname,
		au.middlename,
		au.created_at,
		au.updated_at
    FROM article AS ar JOIN author AS au ON ar.author_id = au.id WHERE ar.id = $1`, id).Scan(
		&res.Id,
		&res.Content.Title,
		&res.Content.Body,
		&res.CreatedAt,
		&updatedAt,
		&res.Author.Id,
		&res.Author.Fullname,
		&tempMiddlename,
		&res.Author.CreatedAt,
		&authorUpdatedAt,
	)
	if err != nil {
		return res, err
	}

	if updatedAt != nil {
		res.UpdatedAt = *updatedAt
	}

	if authorUpdatedAt != nil {
		res.Author.UpdatedAt = *authorUpdatedAt
	}

	if tempMiddlename != nil {
		res.Author.Middlename = *tempMiddlename
	}

	if deletedAt != nil {
		return res, errors.New("article not found")
	}

	return res, err
}

// *=========================================================================
func (stg Postgres) GetArticleList(offset, limit int, search string) (*blogpost.GetArticleListResponse, error) {
	res := &blogpost.GetArticleListResponse{
		Articles: make([]*blogpost.Article, 0),
	}
	rows, err := stg.homeDB.Queryx(`SELECT
	id,
	title,
	body,
	author_id,
	created_at,
	updated_at
	FROM article WHERE deleted_at IS NULL AND ((title ILIKE '%' || $1 || '%') OR (body ILIKE '%' || $1 || '%'))
	LIMIT $2
	OFFSET $3
	`, search, limit, offset)

	if err != nil {
		return res, err
	}

	for rows.Next() {
		a := &blogpost.Article{
			Content: &blogpost.Content{},
		}

		var updatedAt *string

		err := rows.Scan(
			&a.Id,
			&a.Content.Title,
			&a.Content.Body,
			&a.AuthorId,
			&a.CreatedAt,
			&updatedAt,
		)
		if err != nil {
			return res, err
		}

		if updatedAt != nil {
			a.UpdatedAt = *updatedAt
		}

		res.Articles = append(res.Articles, a)
	}
	return res, err
}

// *=========================================================================
func (stg Postgres) UpdateArticle(box *blogpost.UpdateArticleRequest) error {
	if box.Content == nil {
		box.Content = &blogpost.Content{}
	}
	res, err := stg.homeDB.NamedExec("UPDATE article  SET title=:t, body=:b, updated_at=now() WHERE deleted_at IS NULL AND id=:id", map[string]interface{}{
		"id": box.Id,
		"t":  box.Content.Title,
		"b":  box.Content.Body,
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
	return errors.New("article not found")
}

// *=========================================================================
func (stg Postgres) DeleteArticle(id string) error {
	res, err := stg.homeDB.Exec("UPDATE article  SET deleted_at=now() WHERE id=$1 AND deleted_at IS NULL", id)
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
	return errors.New("article not found")
}
