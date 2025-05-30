package models

var migrates []interface{}

func init() {
	registerMigrates()
}

func registerMigrates() {
	migrates = append(migrates, &Account{})
}

func GetMigrates() []interface{} {
	return migrates
}
