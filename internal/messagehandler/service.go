package messagehandler

import (
	"github.com/Restyx/golang-reviews-service/internal/model"
	"github.com/Restyx/golang-reviews-service/internal/store"
)

type ServiceI interface {
	Create(*model.Review) error
	Update(*model.Review) error
	Delete(uint) error
	ReadOne(uint) (*model.Review, error)
	ReadAll() ([]model.Review, error)
}

type Service struct {
	store store.StoreI
}

func NewService(store store.StoreI) ServiceI {
	return &Service{
		store: store,
	}
}

func (h *Service) Create(data *model.Review) error {
	return h.store.Review().Create(data)
}

func (h *Service) Update(data *model.Review) error {
	return h.store.Review().Update(data)
}

func (h *Service) Delete(id uint) error {
	return h.store.Review().Delete(id)
}

func (h *Service) ReadOne(id uint) (*model.Review, error) {
	return h.store.Review().FindOne(id)
}

func (h *Service) ReadAll() ([]model.Review, error) {
	return h.store.Review().FindAll()
}
