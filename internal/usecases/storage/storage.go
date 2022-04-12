package storage

import "market/internal/entities"

type Controller interface {
	GetBalance(profileId int) (float64, error)
	UpdateBalance(profileId int, amount float64) (float64, error)
}

type Repository interface {
	GetProfile(profileId int) (entities.Profile, error)
	UpdateBalance(profileId int, newBalance float64) error
}

type application struct {
	repo Repository
}

func New(repo Repository) *application {
	return &application{
		repo: repo,
	}
}

func (app application) GetBalance(profileId int) (float64, error) {
	profile, err := app.repo.GetProfile(profileId)
	if err != nil {
		return 0, AccountNotFoundError
	}
	return profile.Balance, nil
}

func (app application) UpdateBalance(profileId int, amount float64) (float64, error) {
	profile, err := app.repo.GetProfile(profileId)
	if err != nil {
		return 0, AccountNotFoundError
	}

	err = app.repo.UpdateBalance(profileId, profile.Balance+amount)
	if err != nil {
		return 0, InternalError
	}

	profile, err = app.repo.GetProfile(profileId)
	return profile.Balance, nil
}
