package storage

import "market/internal/entities"

type Controller interface {
	GetBalance(profileId int) (float64, error)
	UpdateBalance(profileId int, amount float64) (float64, error)
	GetBriefcases(profileId int) ([]entities.Briefcase, error)
	BuyBriefcase(profileId int, ticker string, amount int) error
	SellBriefcase(profileId int, ticker string, amount int) error
}

type Repository interface {
	GetBalance(profileId int) (entities.Profile, error)
	UpdateBalance(profileId int, newBalance float64) error
	GetBriefcases(profileId int) ([]entities.Briefcase, error)
	GetBriefcaseByTicker(profileId int, ticker string) (entities.Briefcase, error)
	UpdateBriefcase(profileId int, ticker string, amount int) error
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
	profile, err := app.repo.GetBalance(profileId)
	if err != nil {
		return 0, AccountNotFoundError
	}
	return profile.Balance, nil
}

func (app application) UpdateBalance(profileId int, amount float64) (float64, error) {
	profile, err := app.repo.GetBalance(profileId)
	if err != nil {
		return 0, AccountNotFoundError
	}

	err = app.repo.UpdateBalance(profileId, profile.Balance+amount)
	if err != nil {
		return 0, InternalError
	}

	profile, err = app.repo.GetBalance(profileId)
	return profile.Balance, nil
}

func (app application) GetBriefcases(profileId int) ([]entities.Briefcase, error) {
	briefcases, err := app.repo.GetBriefcases(profileId)
	if err != nil {
		return []entities.Briefcase{}, AccountNotFoundError
	}
	return briefcases, nil
}

func (app application) BuyBriefcase(profileId int, ticker string, amount int) error {
	profile, err := app.repo.GetBalance(profileId)
	if err != nil {
		return AccountNotFoundError
	}

	tickerCost := 1000

	price := float64(tickerCost * amount)

	if profile.Balance < price {
		return NotEnoughMoneyError
	}

	err = app.repo.UpdateBriefcase(profileId, ticker, amount)
	err = app.repo.UpdateBalance(profileId, profile.Balance-price)
	return nil
}
