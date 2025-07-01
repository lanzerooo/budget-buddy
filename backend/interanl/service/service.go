package service

import (
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/model"
	"github.com/lanzerooo/budget-buddy.git/budjet-buddy/interanl/repository"
)

type Service interface {
	CreateTransaction(tx model.Transaction) error
	GetAllTransactions() ([]model.Transaction, error)
	GetBalance() (float64, error)
}

type service struct {
	repo repository.Repository
}

func NewService(r repository.Repository) Service {
	return &service{repo: r}
}

func (s *service) CreateTransaction(tx model.Transaction) error {
	return s.repo.CreateTransaction(tx)
}

func (s *service) GetAllTransactions() ([]model.Transaction, error) {
	return s.repo.GetAllTransactions()
}

func (s *service) GetBalance() (float64, error) {
	return s.repo.GetBalance()
}
