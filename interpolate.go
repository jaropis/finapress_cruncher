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

func interpolate(column []float64) (result []float64) {
	zero := -2.
	nplus1 := -2.
	badCounter := 0
	insideBad := false
	result = make([]float64, len(column))
	for i, elem := range column {
		if elem == -1 && !insideBad {
			insideBad = true
			if i == 0 {
				zero = -2
			} else {
				zero = column[i-1]
			}
			for j := 0; j+i < len(column); j++ {
				if column[j+i] != -1 {
					nplus1 = column[j+i]
					break
					// dalej
				} else {
					badCounter++
				}
			}
		}

		if elem == -1 && insideBad {
			if i == 0 {
				zero = nplus1
			}
			if i == (len(column) - 1) {
				nplus1 = zero
			}
			// fmt.Println("interpoluje", i, zero, nplus1)
			result[i] = (zero + nplus1) / 2
			badCounter--
			if badCounter == 0 {
				insideBad = false
			}
		}

		if elem != -1 {
			result[i] = column[i]
		}
	}
	return result
}
