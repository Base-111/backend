package main

import "github.com/Base-111/backend/internal/api"

func main() {
	err := api.Run()
	if err != nil {
		panic(err)
	}
}
