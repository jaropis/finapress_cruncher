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

func interpolate(zero []float64, storage [][]float64, nplus1 []float64) {

}
