package main

import (
	"encoding/csv"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func processCSVFiles(directory string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".a00") {
			processFile(filepath.Join(directory, file.Name()))
		}
	}
}

func processFile(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	data, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	//interpolation
	interpolate1()
	interpolate2()
	outputFile, err := os.Create("processed_" + filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = '\t'

	err = writer.WriteAll(data)
	if err != nil {
		log.Fatal(err)
	}

	writer.Flush()
}

func main() {
	processCSVFiles(".")
}
