package postgres

import (
	"strconv"

	"github.com/D3vR4pt0rs/logger"
	"market/internal/entities"

	"github.com/jackc/pgx"
)

type Config struct {
	Username string
	Password string
	Ip       string
	Port     string
	Database string
}

type dbClient struct {
	client *pgx.Conn
}

func New(cnfg Config) *dbClient {
	port, _ := strconv.Atoi(cnfg.Port)
	postgressConfig := pgx.ConnConfig{Host: cnfg.Ip, Port: uint16(port), User: cnfg.Username, Password: cnfg.Password, Database: cnfg.Database}
	conn, err := pgx.Connect(postgressConfig)
	if err != nil {
		logger.Error.Println(err.Error())
	}
	return &dbClient{
		client: conn,
	}
}

func (postgres *dbClient) GetProfileById(id int) (entities.Profile, error) {
	var profile Profile
	err := postgres.client.QueryRow("select id, balance from profiles where id=$1", id).Scan(&profile.ID, &profile.Balance)
	if err != nil {
		logger.Error.Println(err.Error())
		return entities.Profile{}, err
	}

	return entities.Profile{
		Id:      int(profile.ID),
		Balance: float64(profile.Balance),
	}, nil
}

func (postgres *dbClient) UpdateBalance(id int, balance float64) error {
	_, err := postgres.client.Exec("update profiles set balance=$1 where id=$2", balance, id)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}
