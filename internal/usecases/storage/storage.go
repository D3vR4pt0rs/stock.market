package storage

import (
	"github.com/D3vR4pt0rs/logger"
	"market/internal/entities"
)

type Controller interface {
	GetBalance(profileId int) (float64, error)
	UpdateBalance(profileId int, amount float64) (float64, error)
	GetBriefcase(profileId int) ([]entities.Briefcase, error)
	BuyStocks(profileId int, ticker string, amount int) error
	SellStocks(profileId int, ticker string, amount int) error
}

type Repository interface {
	GetBalance(profileId int) (entities.Profile, error)
	GetTickerPrice(ticker string) (float64, error)
	UpdateBalance(profileId int, newBalance float64) error
	GetBriefcase(profileId int) ([]entities.Briefcase, error)
	GetStockInBriefcaseByTicker(profileId int, ticker string) (entities.Briefcase, error)
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

func (app application) GetBriefcase(profileId int) ([]entities.Briefcase, error) {
	briefcases, err := app.repo.GetBriefcase(profileId)
	if err != nil {
		return []entities.Briefcase{}, AccountNotFoundError
	}
	return briefcases, nil
}

func (app application) BuyStocks(profileId int, ticker string, amount int) error {
	profile, err := app.repo.GetBalance(profileId)
	if err != nil {
		return AccountNotFoundError
	}

	tickerCost, err := app.repo.GetTickerPrice(ticker)
	if err != nil {
		return InternalError
	}

	price := tickerCost * float64(amount)
	logger.Info.Println("Price for this operation: ", price)

	if profile.Balance < price {
		return NotEnoughMoneyError
	}

	logger.Info.Printf("Change balance on account")
	err = app.repo.UpdateBalance(profileId, profile.Balance-price)
	logger.Info.Printf("Add stocks to account")
	err = app.repo.UpdateBriefcase(profileId, ticker, amount)

	return nil
}

func (app application) SellStocks(profileId int, ticker string, amount int) error {
	briefcase, err := app.repo.GetStockInBriefcaseByTicker(profileId, ticker)
	if err != nil {
		return InternalError
	}

	if briefcase.Amount < amount {
		return NotEnoughStocksError
	}

	logger.Info.Printf("Move stocks from account")
	tickerCost, err := app.repo.GetTickerPrice(ticker)
	if err != nil {
		return InternalError
	}

	price := tickerCost * float64(amount)

	err = app.repo.UpdateBriefcase(profileId, ticker, briefcase.Amount-amount)

	profile, err := app.repo.GetBalance(profileId)
	if err != nil {
		return InternalError
	}

	err = app.repo.UpdateBalance(profileId, profile.Balance+price)
	if err != nil {
		return InternalError
	}
	return nil
}
