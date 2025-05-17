package main

import (
	"fmt"

	"github.com/mehmetcc/go-ems/internal/config"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println(config)
}
