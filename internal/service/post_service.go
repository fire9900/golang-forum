package service

import (
	"forum/internal/model"
	"forum/internal/repository"
)

type PostService struct {
	repo *repository.PostRepository
}

func NewPostService(repo *repository.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) Create(post *model.Post) error {
	return s.repo.Create(post)
}

func (s *PostService) GetByID(id int64) (*model.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) GetAll(page, perPage int) ([]*model.Post, error) {
	offset := (page - 1) * perPage
	return s.repo.GetAll(perPage, offset)
}

func (s *PostService) Update(post *model.Post) error {
	return s.repo.Update(post)
}

func (s *PostService) Delete(id, authorID int64) error {
	return s.repo.Delete(id, authorID)
}
