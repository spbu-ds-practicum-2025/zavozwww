package main

import (
	"authServ/pkg/logger"
	"fmt"
)

func main() {
	logger, err := logger.NewLogger()
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	fmt.Println(logger.Level())
}
