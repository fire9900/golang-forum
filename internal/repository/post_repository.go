package repository

import (
	"database/sql"
	"forum/internal/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(post *model.Post) error {
	query := `
		INSERT INTO posts (title, content, author_id, created_at, updated_at)
		VALUES (?, ?, ?, NOW(), NOW())
	`
	result, err := r.db.Exec(query, post.Title, post.Content, post.AuthorID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	post.ID = id
	return nil
}

func (r *PostRepository) GetByID(id int64) (*model.Post, error) {
	post := &model.Post{}
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts WHERE id = ?
	`
	err := r.db.QueryRow(query, id).Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.AuthorID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (r *PostRepository) GetAll(limit, offset int) ([]*model.Post, error) {
	query := `
		SELECT id, title, content, author_id, created_at, updated_at
		FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?
	`
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*model.Post
	for rows.Next() {
		post := &model.Post{}
		err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.AuthorID,
			&post.CreatedAt,
			&post.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func (r *PostRepository) Update(post *model.Post) error {
	query := `
		UPDATE posts 
		SET title = ?, content = ?, updated_at = NOW()
		WHERE id = ? AND author_id = ?
	`
	result, err := r.db.Exec(query, post.Title, post.Content, post.ID, post.AuthorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *PostRepository) Delete(id, authorID int64) error {
	query := "DELETE FROM posts WHERE id = ? AND author_id = ?"
	result, err := r.db.Exec(query, id, authorID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
