package models

var migrates []interface{}

func init() {
	registerMigrates()
}

func registerMigrates() {
	migrates = append(migrates, &Product{})
}

func GetMigrates() []interface{} {
	return migrates
}
