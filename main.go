package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func processCSVFiles(directory string, extension string) {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".a00") {
			processFile(filepath.Join(directory, file.Name()), extension)
		}
	}
}

func processFile(filepath string, extension string) {
	matrix := [][]float64{}
	insideBad := false
	zero := []float64{}
	storage := [][]float64{}
	nplus1 := []float64{}
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
		if i > 6 {
			dummy, _ := convertToFloats(row)
			if !anyBad(dummy) && !insideBad {
				matrix = append(matrix, dummy)
			}

			if anyBad(dummy) && insideBad {
				storage = append(storage, dummy)
			}
			if anyBad(dummy) && (i == 7) {
				zero = genrateBads(dummy)
				insideBad = true
			}

			// you are not insideBad, but a bad row happens - zero is the previous row
			if anyBad(dummy) && !insideBad {
				insideBad = true
				//zero = matrix[i-1]
			}

			// you are inside bad region, but it ends
			if !anyBad(dummy) && insideBad {
				if i == len(data) {
					nplus1 = genrateBads(dummy)
				}
				nplus1 = dummy
				insideBad = false
				for _, storageRow := range storage {
					interpolate(zero, storage, nplus1)
					matrix = append(matrix, storageRow)
				}

				if i < len(data) {
					matrix = append(matrix, nplus1)
				}
				storage = [][]float64{} // reset
			}
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

	newFilepath := strings.TrimSuffix(filepath, path.Ext(filepath)) + "." + extension
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

func addNewLineAfterLine(directory string, lineNumber int, extension string) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("Error reading directory %s: %v\n", directory, err)
		return
	}

	for _, file := range files {
		// Skip if not a file or doesn't have the right extension
		if file.IsDir() || filepath.Ext(file.Name()) != "."+extension {
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
	processCSVFiles(".", "dum")
	addNewLineAfterLine(".", 3, "dum")
}
