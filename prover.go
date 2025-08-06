package main 

import (
	"gitlab.com/oelmekki/matrix"
	"github.com/cloudflare/circl/ecc/bls12381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"fmt"
)

func prove(L, R, O matrix.Matrix) (bls12381.G1, bls12381.G2, bls12381.G1) {
	// witness for our specific problem in `main.go`:  w = [1, y, v, x] = [1, 155, 25, 5]
	W, _ := matrix.Build(
		matrix.Builder{
			matrix.Row{1},
			matrix.Row{155},
			matrix.Row{25},
			matrix.Row{5},
		},
	)

	// sanity checks
	if !matricesSanityChecks(L, R, O, W) {
		panic("Invalid Matrices!")
	}
	u_x, v_x, w_x, _, h_x := R1CSToQAP(L, R, O, W)

	SRS1, SRS2, SRS3 := generateSRS(int(len(u_x)), int(len(h_x)), L.Rows())
	// TODO add sanity checks on SRSs, that they power was generated successfully...
	
	A := EvalLAtSRS1(u_x, SRS1)
	B := EvalRAtSRS2(v_x, SRS2)
	C := EvalOutputAtSRS13(w_x, h_x, SRS1, SRS3)	

	return A, B, C
}

func matricesSanityChecks(L, R, O, W matrix.Matrix) bool {
	return (L.Rows() == R.Rows() && L.Rows() == O.Rows() && L.Cols() == R.Cols() && L.Cols() == O.Cols() && W.Cols() == 1 && W.Rows() == L.Cols())
}

func EvalLAtSRS1(u_x polynomial.Polynomial, srs []bls12381.G1) bls12381.G1 {
	if len(u_x) != len(srs) {
		panic("Incorrect SRS")
	}
	var A bls12381.G1 
	coeffs := u_x
	coeff := frElementToScalar(coeffs[0])
	A.ScalarMult(&coeff, &srs[0])
	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToScalar(coeffs[i])
		var tmp bls12381.G1 
		tmp.ScalarMult(&coeff, &srs[i])
		A.Add(&A, &tmp)
	}

	return A
}

func EvalRAtSRS2(v_x polynomial.Polynomial, srs []bls12381.G2) bls12381.G2 {
	if len(v_x) != len(srs) {
		panic("Incorrect SRS")
	}
	var B bls12381.G2 
	coeffs := v_x
	coeff := frElementToScalar(coeffs[0])
	B.ScalarMult(&coeff, &srs[0])

	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToScalar(coeffs[i])
		var tmp bls12381.G2 
		tmp.ScalarMult(&coeff, &srs[i])
		B.Add(&B, &tmp)
	}

	return B
}

func EvalOutputAtSRS13(w_x, h_x polynomial.Polynomial, srs1 , srs3 []bls12381.G1) bls12381.G1 {
	if len(w_x) != len(srs1) || len(h_x) != len(srs3) {
		panic("Incorrect SRS")
	}
	var C bls12381.G1 
	coeffs := w_x
	coeff := frElementToScalar(coeffs[0])
	C.ScalarMult(&coeff, &srs1[0])

	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToScalar(coeffs[i])
		var tmp bls12381.G1 
		tmp.ScalarMult(&coeff, &srs1[i])
		C.Add(&C, &tmp)
	}

	coeffs = h_x
	for i:=0; i < len(coeffs); i++ {
		var coeff bls12381.Scalar
		coeff = frElementToScalar(coeffs[i])
		var tmp bls12381.G1 
		tmp.ScalarMult(&coeff, &srs3[i])
		C.Add(&C, &tmp)
	}

	return C
}

func frElementToScalar(e fr.Element) bls12381.Scalar {
    b := e.Bytes()
    var s bls12381.Scalar
    s.SetBytes(b[:])
    return s
}
