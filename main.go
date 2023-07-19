package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
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
	matrix := [][]float64{}
	// insideBad := false
	// zero := []float64{}
	// storage := [][]float64{}
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.FieldsPerRecord = -1 // Allow variable number of fields

	data, err := reader.ReadAll()
	data_copy := data
	if err != nil {
		log.Fatal(err)
	}

	//interpolation
	for i, row := range data {

		// Here, row is a slice of strings (a row in the CSV file)
		if i >= 7 {
			dummy, _ := convertToFloats(row)
			matrix = append(matrix, dummy)
		}
	}

	// go through all the columns
	for i := 0; i < len(matrix[0]); i++ {
		// read a column
		column := make([]float64, len(matrix))
		for j := 0; j < len(matrix); j++ {
			column[j] = matrix[j][i]
		}
		column = interpolate(column)
		// copy interpolated column back to matrix
		for j := 0; j < len(matrix); j++ {
			matrix[j][i] = column[j]
		}
	}

	// copy interpolated matrix to str data
	for i, dummy := range matrix {
		dummyStr := make([]string, len(dummy))
		for i, v := range dummy {
			dummyStr[i] = strconv.FormatFloat(v, 'f', -1, 64)
		}
		if (i + 7) < len(data) {
			data_copy[i+7] = dummyStr
		}
	}

	newFilepath := "processed_" + filepath
	outputFile, err := os.Create(newFilepath)

	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	writer := csv.NewWriter(outputFile)
	writer.Comma = '\t'

	err = writer.WriteAll(data_copy)
	if err != nil {
		log.Fatal(err)
	}

	writer.Flush()
}

func addNewLineAfterLine(directory string, lineNumber int) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", directory, err)
		return
	}

	for _, file := range files {
		// Skip if not a file or doesn't have the right extension
		if file.IsDir() || !strings.HasPrefix(file.Name(), "processed_") {
			continue
		}

		// Open the original file
		originalPath := filepath.Join(directory, file.Name())
		original, err := os.Open(originalPath)
		if err != nil {
			fmt.Printf("Error opening file %s: %v\n", originalPath, err)
			continue
		}

		// Create a buffer to hold the file contents
		var buffer strings.Builder

		// Create a scanner to read the file line by line
		scanner := bufio.NewScanner(original)
		currentLine := 1
		for scanner.Scan() {
			// Write the line to the buffer
			buffer.WriteString(scanner.Text() + "\n")
			// If it's the desired line, write an additional newline
			if currentLine == lineNumber {
				buffer.WriteString("\n")
			}
			currentLine++
		}

		if err := scanner.Err(); err != nil {
			fmt.Printf("Error reading file %s: %v\n", originalPath, err)
		}

		original.Close()

		// Open the original file again, this time for writing
		original, err = os.OpenFile(originalPath, os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Printf("Error opening file %s for writing: %v\n", originalPath, err)
			continue
		}

		// Write the buffer contents to the file
		original.WriteString(buffer.String())

		original.Close()
	}
}

func main() {
	processCSVFiles(".")
	addNewLineAfterLine(".", 3)
}
