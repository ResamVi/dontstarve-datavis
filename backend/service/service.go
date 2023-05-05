package service

import (
	"github.com/ResamVi/dontstarve-datavis/storage"
)

// Service Layer
type Service struct {
	store storage.Store
}

func (s Service) Started() string {
	return s.store.GetAge().Format("2006-01-02")
}

func New(store storage.Store) Service {
	store.Start()

	return Service{
		store: store,
	}
}
