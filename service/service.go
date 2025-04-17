package service

import (
	"block-test/config"
	"block-test/repository"

	"github.com/inconshreveable/log15"
)

type Service struct {
	config     *config.Config
	repository *repository.Repository
	log        log15.Logger

	difficulty int64
}

func NewService(repository *repository.Repository, difficulty int64) *Service {
	s := &Service{
		repository: repository,
		log:        log15.New("module", "service"),
	}

	return s
}
