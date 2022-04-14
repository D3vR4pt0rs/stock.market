package postgres

type Profile struct {
	ID      int32
	Balance float32
}

type Briefcase struct {
	Id     int32
	Ticker string
	Amount int32
}
