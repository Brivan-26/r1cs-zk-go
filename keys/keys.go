package keys

import (
	"fmt"
	"encoding/json"
	"io/ioutil"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

type ProvingKey struct {
	SRS1      []G1AffineJSON `json:"srs1"`
	SRS2      []G2AffineJSON `json:"srs2"`
	SRS3      []G1AffineJSON `json:"srs3"`
	Alpha     G1AffineJSON   `json:"alpha"`
	Beta      G2AffineJSON   `json:"beta"`
	ProverPsi []G1AffineJSON `json:"proverPsi"`
}

type VerifyingKey struct {
	Alpha       G1AffineJSON   `json:"alpha"`
	Beta        G2AffineJSON   `json:"beta"`
	Gamma       G2AffineJSON   `json:"gamma"`
	Teta        G2AffineJSON   `json:"teta"`
	VerifierPsi []G1AffineJSON `json:"verifierPsi"`
}

type Proof struct {
	A G1AffineJSON `json:"A"`
	B G2AffineJSON `json:"B"`
	C G1AffineJSON `json:"C"`
}

type G1AffineJSON struct {
	X string `json:"x"`
	Y string `json:"y"`
}

type G2AffineJSON struct {
	X0 string `json:"x0"`
	X1 string `json:"x1"`
	Y0 string `json:"y0"`
	Y1 string `json:"y1"`
}

func g1AffineToJSON(point curve.G1Affine) G1AffineJSON {
	x := point.X.String()
	y := point.Y.String()
	return G1AffineJSON{
		X: x,
		Y: y,
	}
}

func g2AffineToJSON(point curve.G2Affine) G2AffineJSON {
	return G2AffineJSON{
		X0: point.X.A0.String(),
		X1: point.X.A1.String(),
		Y0: point.Y.A0.String(),
		Y1: point.Y.A1.String(),
	}
}

func g1SliceToJSON(points []curve.G1Affine) []G1AffineJSON {
	result := make([]G1AffineJSON, len(points))
	for i, point := range points {
		result[i] = g1AffineToJSON(point)
	}
	return result
}

func g2SliceToJSON(points []curve.G2Affine) []G2AffineJSON {
	result := make([]G2AffineJSON, len(points))
	for i, point := range points {
		result[i] = g2AffineToJSON(point)
	}
	return result
}

func saveToJSONFile(filename string, data interface{}) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s: %v", filename, err)
	}
	
	err = ioutil.WriteFile(filename, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write %s: %v", filename, err)
	}
	
	return nil
}

func JsonToG1Affine(jsonPoint G1AffineJSON) curve.G1Affine {
	var point curve.G1Affine
	
	point.X.SetString(jsonPoint.X)
	point.Y.SetString(jsonPoint.Y)
	
	return point
}

func JsonToG2Affine(jsonPoint G2AffineJSON) curve.G2Affine {
	var point curve.G2Affine
	
	
	point.X.A0.SetString(jsonPoint.X0)
	point.X.A1.SetString(jsonPoint.X1)
	point.Y.A0.SetString(jsonPoint.Y0)
	point.Y.A1.SetString(jsonPoint.Y1)
	
	return point
}

func JsonToG1AffineSlice(jsonPoints []G1AffineJSON) []curve.G1Affine {
	points := make([]curve.G1Affine, len(jsonPoints))
	
	for i, jsonPoint := range jsonPoints {
		point := JsonToG1Affine(jsonPoint)
		points[i] = point
	}
	
	return points
}

func JsonToG2AffineSlice(jsonPoints []G2AffineJSON) []curve.G2Affine {
	points := make([]curve.G2Affine, len(jsonPoints))
	
	for i, jsonPoint := range jsonPoints {
		point := JsonToG2Affine(jsonPoint)
		points[i] = point
	}
	
	return points
}

func BuildPk(srs1, srs3, proverPsi []curve.G1Affine, srs2 []curve.G2Affine, alpha curve.G1Affine, beta curve.G2Affine) {
	pk := ProvingKey{
		SRS1:      g1SliceToJSON(srs1),
		SRS2:      g2SliceToJSON(srs2),
		SRS3:      g1SliceToJSON(srs3),
		Alpha:     g1AffineToJSON(alpha),
		Beta:      g2AffineToJSON(beta),
		ProverPsi: g1SliceToJSON(proverPsi),
	}

	err := saveToJSONFile("pk.json", pk)
	if err != nil {
		panic(fmt.Sprintf("Failed to save proving key: %v", err))
	}

	fmt.Println("Proving key saved to pk.json")
}

func BuildVk(alpha curve.G1Affine, verifierPsi []curve.G1Affine, beta, gamma, teta curve.G2Affine) {
	vk := VerifyingKey{
		Alpha:       g1AffineToJSON(alpha),
		Beta:        g2AffineToJSON(beta),
		Gamma:       g2AffineToJSON(gamma),
		Teta:        g2AffineToJSON(teta),
		VerifierPsi: g1SliceToJSON(verifierPsi),
	}
	
	err := saveToJSONFile("vk.json", vk)
	if err != nil {
		panic(fmt.Sprintf("Failed to save verifying key: %v", err))
	}
	fmt.Println("Verifying key saved to vk.json")
}

func BuildProof(a, c curve.G1Affine, b curve.G2Affine) {
	proof := Proof{
		A: g1AffineToJSON(a),
		B: g2AffineToJSON(b),
		C: g1AffineToJSON(c),
	}
	
	err := saveToJSONFile("proof.json", proof)
	if err != nil {
		panic(fmt.Sprintf("Failed to save proof: %v", err))
	}
	fmt.Println("Proof saved to proof.json")
}