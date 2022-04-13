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

func (postgres *dbClient) GetBriefcasesByAccountId(accountId int) ([]entities.Briefcase, error) {
	rows, err := postgres.client.Query("select b.id, b.ticker, b.amount from  briefcase b, profiles p where b.account_id = p.id and b.account_id=$1", accountId)
	if err != nil {
		logger.Error.Println(err.Error())
		return []entities.Briefcase{}, err
	}

	defer rows.Close()

	var briefcases []entities.Briefcase
	for rows.Next() {
		var briefcase Briefcase
		err := rows.Scan(&briefcase)
		if err != nil {
			logger.Error.Println(err.Error())
			return []entities.Briefcase{}, err
		}
		briefcases = append(briefcases, entities.Briefcase{Id: int(briefcase.Id), Ticker: briefcase.Ticker, Amount: int(briefcase.Amount)})
	}

	if rows.Err() != nil {
		return []entities.Briefcase{}, err
	}
	return briefcases, nil
}

func (postgres *dbClient) GetBriefcaseByTicker(accountId int, ticker string) (entities.Briefcase, error) {
	var briefcase Briefcase
	err := postgres.client.QueryRow("select id,ticker,amount from briefcase where account_id=$1 and ticker=$2", accountId, ticker).Scan(&briefcase.Id, &briefcase.Ticker, &briefcase.Amount)
	if err != nil {
		logger.Error.Println(err.Error())
		return entities.Briefcase{}, err
	}

	return entities.Briefcase{
		Id:     int(briefcase.Id),
		Ticker: briefcase.Ticker,
		Amount: int(briefcase.Amount),
	}, nil
}

func (postgres *dbClient) UpdateBriefcaseAmount(id int, amount int) error {
	_, err := postgres.client.Exec("update briefcase set amount=$1 where id=$2", amount, id)

	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func (postgres *dbClient) DeleteBriefcase(id int) error {
	_, err := postgres.client.Exec("delete from briefcase where id=$1", id)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}

func (postgres *dbClient) AddBriefcase(ticker string, amount int) error {
	_, err := postgres.client.Exec("insert into briefcase (ticker,amount) values ($1,$2)", ticker, amount)
	if err != nil {
		logger.Error.Println(err.Error())
		return err
	}
	return nil
}
