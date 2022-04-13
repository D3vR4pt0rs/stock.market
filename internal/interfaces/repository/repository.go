package repository

import (
	"github.com/D3vR4pt0rs/logger"
	"market/internal/entities"
)

type database interface {
	GetProfileById(id int) (entities.Profile, error)
	UpdateBalance(id int, balance float64) error
	GetBriefcasesByAccountId(profileId int) ([]entities.Briefcase, error)
	GetBriefcaseByTicker(profileId int, ticker string) (entities.Briefcase, error)
	UpdateBriefcaseAmount(id int, amount int) error
	AddBriefcase(profileId int, ticker string, amount int) error
}

type driver struct {
	d database
}

func New(dbHandler database) *driver {
	return &driver{
		d: dbHandler,
	}
}

func (d driver) GetBalance(id int) (entities.Profile, error) {
	return d.d.GetProfileById(id)
}

func (d driver) GetBriefcasesByAccountId(profileId int) ([]entities.Briefcase, error) {
	return d.d.GetBriefcasesByAccountId(profileId)
}

func (d driver) GetBriefcaseByTicker(profileId int, ticker string) (entities.Briefcase, error) {
	return d.d.GetBriefcaseByTicker(profileId, ticker)
}

func (d driver) UpdateBalance(id int, balance float64) error {
	return d.d.UpdateBalance(id, balance)
}

func (d driver) UpdateBriefcaseAmount(id int, amount int) error {
	return d.d.UpdateBriefcaseAmount(id, amount)
}

func (d driver) UpdateBriefcase(profileId int, ticker string, amount int) error {
	briefcase, err := d.d.GetBriefcaseByTicker(profileId, ticker)
	if err != nil {
		err := d.d.AddBriefcase(profileId, ticker, amount)
		if err != nil {
			logger.Info.Println("Failed to add briefcase")
			return err
		}
	} else {
		err := d.d.UpdateBriefcaseAmount(briefcase.Id, briefcase.Amount+amount)
		if err != nil {
			logger.Info.Println("Failed to update briefcase")
			return err
		}
	}
	return nil
}
