package service

import (
	"challenge2/models"
	"challenge2/repository"
)

type ItemService struct {
	repo *repository.ItemRepository
}

func NewItemService(repo *repository.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

func (s *ItemService) GetItem() ([]models.Item, error) {
	return s.repo.GetAll()
}

func (s *ItemService) CreateItem(item *models.Item) error {
	return s.repo.Create(item)
}

func (s *ItemService) DeleteItem(id int) error {
	return s.repo.DeleteByID(id)
}