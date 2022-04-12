package main

import (
	"fmt"
	"github.com/spf13/viper"
)

type User struct {
	name     string
	password string
	rights   string
}

func main() {
	initConfig()
	users := getUsersFromConfig(viper.AllSettings())
	fmt.Println(users)
}

func initConfig() {
	viper.SetConfigFile("config.yml")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w\n", err))
	}
}

func getUsersFromConfig(usersMap map[string]interface{}) []User {
	users := make([]User, 0)
	for key, element := range usersMap {
		fmt.Println("KEY: ", key, " => ", "ELEMENT: ", element)
		el := element.(map[string]interface{})

		users = append(users, User{
			name:     fmt.Sprintf("%v", el["name"]),
			password: fmt.Sprintf("%v", el["password"]),
			rights:   fmt.Sprintf("%v", el["rights"]),
		})
	}
	return users
}
