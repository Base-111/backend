package main

import (
	_ "github.com/Base-111/backend/docs"
	"github.com/Base-111/backend/internal/api"
)

// @title           Base-111 Backend API
// @version         0.0.1
// @description     API хакатона
// @termsOfService  http://swagger.io/terms/

// @contact.name   Team Base-111
// @contact.email  koreshkov200@mail.ru

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8081
// @BasePath  /

func main() {
	err := api.Run()
	if err != nil {
		panic(err)
	}
}
