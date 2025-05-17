package main

import (
	"log"
	"sort"

	"github.com/mehmetcc/go-ems/internal/config"
	"github.com/mehmetcc/go-ems/internal/reader"
	"github.com/mehmetcc/go-ems/internal/writer"
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

	fileWriter, err := writer.New("output")
	if err != nil {
		panic(err)
	}

	count := 0

	for {
		batch, err := fileReader.Read(config.Application.BatchSize)
		if err != nil {
			log.Fatalf("Error reading batch: %v", err)
		}
		if batch == nil {
			break
		}

		sort.Ints(batch)
		if err := fileWriter.Write(count, batch); err != nil {
			log.Fatalf("Error writing batch %d: %v", count, err)
		}
		count++
	}
}
