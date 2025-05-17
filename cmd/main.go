package main

import (
	"fmt"
	"log"

	"github.com/mehmetcc/go-ems/internal/config"
	"github.com/mehmetcc/go-ems/internal/reader"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fileReader, err := reader.New("input.txt")
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}
	defer fileReader.Close()

	count := 0

	for {
		batch, err := fileReader.Read(config.Application.BatchSize)
		if err != nil {
			log.Fatalf("Error reading batch: %v", err)
		}
		if batch == nil {
			break
		}

		count++
	}

	fmt.Printf("Read %d batches\n", count)
}
