package repository

import "market/internal/entities"

type database interface {
	GetProfileById(id int) (entities.Profile, error)
	UpdateBalance(id int, balance float64) error
}

type driver struct {
	d database
}

func New(dbHandler database) *driver {
	return &driver{
		d: dbHandler,
	}
}

func (d driver) GetProfile(id int) (entities.Profile, error) {
	return d.d.GetProfileById(id)
}

func (d driver) UpdateBalance(id int, balance float64) error {
	return d.d.UpdateBalance(id, balance)
}
