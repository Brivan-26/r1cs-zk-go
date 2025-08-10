package r1cs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"gitlab.com/oelmekki/matrix"
)

type R1CSData struct {
	L [][]int `json:"L"`
	R [][]int `json:"R"`
	O [][]int `json:"O"`
}

// LoadR1CSFromJSON reads and parses the R1CS JSON file
func LoadR1CSFromJSON() (matrix.Matrix, matrix.Matrix, matrix.Matrix, error) {
	jsonData, err := ioutil.ReadFile("r1cs.json")
	if err != nil {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("failed to read R1CS file: %v", err)
	}

	var r1csData R1CSData
	err = json.Unmarshal(jsonData, &r1csData)
	if err != nil {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("failed to parse R1CS JSON: %v", err)
	}

	if len(r1csData.L) != len(r1csData.R) || len(r1csData.R) != len(r1csData.O) {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("R1CS matrices must have the same number of rows")
	}

	// sanity checks that all matrices' rows have the same column number
	if len(r1csData.L) > 0 {
		lCols := len(r1csData.L[0])
		for i, row := range r1csData.L {
			if len(row) != lCols {
				return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("L matrix row %d has inconsistent column count", i)
			}
		}
		for i, row := range r1csData.R {
			if len(row) != lCols {
				return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("R matrix row %d has inconsistent column count", i)
			}
		}
		for i, row := range r1csData.O {
			if len(row) != lCols {
				return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("O matrix row %d has inconsistent column count", i)
			}
		}
	}

	// Build L matrix
	L, err := buildMatrixFromIntArray(r1csData.L)
	if err != nil {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("failed to build L matrix: %v", err)
	}

	// Build R matrix
	R, err := buildMatrixFromIntArray(r1csData.R)
	if err != nil {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("failed to build R matrix: %v", err)
	}

	// Build O matrix
	O, err := buildMatrixFromIntArray(r1csData.O)
	if err != nil {
		return matrix.Matrix{}, matrix.Matrix{}, matrix.Matrix{}, fmt.Errorf("failed to build O matrix: %v", err)
	}

	return L, R, O, nil
}

// buildMatrixFromIntArray converts a 2D int array to a matrix.Matrix
func buildMatrixFromIntArray(data [][]int) (matrix.Matrix, error) {
	if len(data) == 0 {
		return matrix.Matrix{}, fmt.Errorf("cannot build matrix from empty data")
	}

	builder := make(matrix.Builder, len(data))
	for i, row := range data {
		builder[i] = make(matrix.Row, len(row))
		for j, val := range row {
			builder[i][j] = float64(val)
		}
	}

	return matrix.Build(builder)
}