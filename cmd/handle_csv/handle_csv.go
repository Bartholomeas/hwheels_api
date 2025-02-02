package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	startTime := time.Now()

	filePath := flag.String("file", "", "Path to the CSV file")
	flag.Parse()

	if *filePath == "" {
		panic("Please provide a file path using -file flag")
	}

	fmt.Println("Processing file:", *filePath)
	if cwd, err := os.Getwd(); err == nil {
		fmt.Println("Current working directory:", cwd)
	}

	if err := handleCSV(*filePath); err != nil {
		panic(err)
	}
	duration := time.Since(startTime)
	fmt.Printf("Execution time: %v\n", duration)

}

func handleCSV(filePath string) error {
	file, err := os.Open(filePath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	reader.LazyQuotes = true
	reader.TrimLeadingSpace = true

	data, err := reader.ReadAll()

	if err != nil {
		panic(err)
	}

	for _, row := range data {
		for _, col := range row {
			fmt.Printf("%s,", col)
		}
		fmt.Println()
	}

	return nil
}
