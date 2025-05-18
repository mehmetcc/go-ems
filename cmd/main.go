package main

import (
	"log"
	"sort"

	"github.com/mehmetcc/go-ems/internal/config"
	"github.com/mehmetcc/go-ems/internal/merger"
	"github.com/mehmetcc/go-ems/internal/reader"
	"github.com/mehmetcc/go-ems/internal/writer"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	fileReader, err := reader.New(config.Application.InputPath)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}
	defer fileReader.Close()

	fileWriter, err := writer.New(config.Application.TempFileDirectory)
	if err != nil {
		panic(err)
	}

	count := 0

	log.Println("Phase 1: Splitting and sorting chunks...")
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

	log.Println("Phase 2: Merging sorted chunks...")
	merger := merger.New(
		config.Application.TempFileDirectory,
		config.Application.OutputPath,
		count,
	)

	if err := merger.Merge(); err != nil {
		log.Fatalf("Error during merge: %v", err)
	}

	log.Println("Cleaning up temporary files...")
	if err := merger.Cleanup(); err != nil {
		log.Printf("Warning: error cleaning up temporary files: %v", err)
	}

	log.Println("External merge sort completed successfully!")
}
