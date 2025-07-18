package models

var migrates []interface{}

func init() {
	registerMigrates()
}

func registerMigrates() {
	migrates = append(migrates, &Product{})
	migrates = append(migrates, &Order{})
	migrates = append(migrates, &User{})
}

func GetMigrates() []interface{} {
	return migrates
}
