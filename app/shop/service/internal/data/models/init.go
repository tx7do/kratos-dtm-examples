package models

var migrates []interface{}

func init() {
	registerMigrates()
}

func registerMigrates() {
	migrates = append(migrates, &Product{})
	migrates = append(migrates, &Order{})
	migrates = append(migrates, &User{})
	migrates = append(migrates, &Stock{})
	migrates = append(migrates, &StockDeductionLog{})
}

func GetMigrates() []interface{} {
	return migrates
}
