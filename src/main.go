package main

import (
	"fmt"
	"pkg/logger"
)

func main() {
	logger, err := logger.NewLogger("config/logger.yaml")
	if err != nil {
		fmt.Println("Error initializing logger:", err)
		return
	}
	fmt.Println(logger.Level())
}
