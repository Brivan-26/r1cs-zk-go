package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"gitlab.com/oelmekki/matrix"
)

type WitnessData struct {
	PublicInputs  []int `json:"publicInputs"`
	PrivateInputs []int `json:"privateInputs"`
}

type PublicWitnessData struct {
	PublicInputs []int `json:"publicInputs"`
}

type PrivateWitnessData struct {
	PrivateInputs []int `json:"privateInputs"`
}

func LoadWitnessFromJSON() (matrix.Matrix, int, error) {
	// Read the JSON file
	jsonData, err := ioutil.ReadFile("witness.json")
	if err != nil {
		return matrix.Matrix{}, 0, fmt.Errorf("failed to read witness file: %v", err)
	}

	var witnessData WitnessData
	err = json.Unmarshal(jsonData, &witnessData)
	if err != nil {
		return matrix.Matrix{}, 0, fmt.Errorf("failed to parse witness JSON: %v", err)
	}
	publicInputsSize := len(witnessData.PublicInputs)
	// Combine public and private inputs: [publicInputs..., privateInputs...]
	combined := append(witnessData.PublicInputs, witnessData.PrivateInputs...)
	
	builder := make(matrix.Builder, len(combined))
	for i, val := range combined {
		builder[i] = matrix.Row{float64(val)}
	}

	witness, err := matrix.Build(builder)
	if err != nil {
		return matrix.Matrix{}, 0, fmt.Errorf("failed to build witness matrix: %v", err)
	}

	return witness, publicInputsSize, nil
}

func LoadPublicInputsFromJSON() ([]int, error) {
	jsonData, err := ioutil.ReadFile("witness.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read witness file: %v", err)
	}

	var publicWitnessData PublicWitnessData
	err = json.Unmarshal(jsonData, &publicWitnessData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse witness JSON: %v", err)
	}

	return publicWitnessData.PublicInputs, nil
}
