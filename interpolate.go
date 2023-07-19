package main

import (
	"strconv"
)

// import packages needed

func convertToFloats(strings []string) ([]float64, error) {
	floats := make([]float64, len(strings))
	for i, s := range strings {
		if s == "" {
			s = "-1"
		}
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		floats[i] = f
	}
	return floats, nil
}

func anyBad(vector []float64) (isbad bool) {
	for _, elem := range vector {
		if elem < 0 {
			return true
		}
	}
	return false
}

func genrateBads(vector []float64) (badsy []float64) {
	badsy = make([]float64, len(vector))
	for i := range badsy {
		badsy[i] = -2.0
	}
	return badsy
}

func interpolate(zero []float64, storage [][]float64, nplus1 []float64) (result [][]float64) {
	if zero[0] == -2 {
		//fmt.Println(zero)
		//fmt.Println(storage)
		//fmt.Println("pierwsza!")
		zero = nplus1
	}

	if nplus1[0] == -2 {
		nplus1 = zero
	}

	means := make([]float64, len(zero))
	for i, _ := range means {
		means[i] = (zero[i] + nplus1[i]) / 2
	}

	for i, _ := range storage {
		for j, _ := range storage[0] {
			if storage[i][j] < 0 {
				storage[i][j] = means[j]
			}
		}
	}
	result = storage
	return result
}

func interpolate_dummy(zero []float64, storage [][]float64, nplus1 []float64) {
}
