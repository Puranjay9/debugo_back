package main

import (
	"debugo_back/server"
	"fmt"
)

func main() {
	api := server.NewServer(":3000")

	if err := api.Run(); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
