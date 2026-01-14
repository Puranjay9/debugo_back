package main

import (
	"debugo_back/server"
	"fmt"
)

func main() {
	api := server.NewServer(":8000")

	if err := api.Run(); err != nil {
		fmt.Println("Error starting the server:", err)
	}
}
