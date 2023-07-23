package config

import "os"

func SetEnv() {
	os.Setenv("user", "matias")
	os.Setenv("pwd", "Demichelis")
	os.Setenv("db_name", "Inmobiliaria_BD")
}

func GetEnv(name string) (user, pwd, db_name string) {
	user = os.Getenv("user")
	pwd = os.Getenv("pwd")
	db_name = os.Getenv("db_name")
	return
}
