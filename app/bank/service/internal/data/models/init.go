package models

var migrates []interface{}

func init() {
	registerMigrates()
}

func registerMigrates() {
}

func GetMigrates() []interface{} {
	return migrates
}
