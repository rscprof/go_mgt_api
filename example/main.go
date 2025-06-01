package main

import (
	"fmt"

	api "github.com/rscprof/go_mgt_api/api"
)

func main() {
	client := api.NewClient()
	data, err := client.GetStopData("9d7f733a-d532-4fca-a922-4c978b79681c")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Остановка: %s\n", data.Name)
}
