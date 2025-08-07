package main 

import (
	"gitlab.com/oelmekki/matrix"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"math/big"
)

func prove(L, R, O matrix.Matrix) (curve.G1Affine, curve.G2Affine, curve.G1Affine) {
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

func EvalLAtSRS1(u_x polynomial.Polynomial, srs []curve.G1Affine) curve.G1Affine {
	if len(u_x) != len(srs) {
		panic("Incorrect SRS")
	}
	var A curve.G1Affine
	coeffs := u_x
	coeff := frElementToBigInt(coeffs[0])
	A.ScalarMultiplication(&srs[0], &coeff)
	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToBigInt(coeffs[i])
		var tmp curve.G1Affine 
		tmp.ScalarMultiplication(&srs[i], &coeff)
		A.Add(&A, &tmp)
	}

	return A
}

func EvalRAtSRS2(v_x polynomial.Polynomial, srs []curve.G2Affine) curve.G2Affine {
	if len(v_x) != len(srs) {
		panic("Incorrect SRS")
	}
	var B curve.G2Affine
	coeffs := v_x
	coeff := frElementToBigInt(coeffs[0])
	B.ScalarMultiplication(&srs[0], &coeff)
	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToBigInt(coeffs[i])
		var tmp curve.G2Affine 
		tmp.ScalarMultiplication(&srs[i], &coeff)
		B.Add(&B, &tmp)
	}

	return B
}

func EvalOutputAtSRS13(w_x, h_x polynomial.Polynomial, srs1 , srs3 []curve.G1Affine) curve.G1Affine {
	if len(w_x) != len(srs1) || len(h_x) != len(srs3) {
		panic("Incorrect SRS")
	}
	var C curve.G1Affine 
	coeffs := w_x
	coeff := frElementToBigInt(coeffs[0])
	C.ScalarMultiplication(&srs1[0], &coeff)

	for i:=1; i < len(coeffs); i++ {
		coeff := frElementToBigInt(coeffs[i])
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&srs1[i], &coeff)
		C.Add(&C, &tmp)
	}

	coeffs = h_x
	for i:=0; i < len(coeffs); i++ {
		coeff = frElementToBigInt(coeffs[i])
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&srs3[i], &coeff)
		C.Add(&C, &tmp)
	}

	return C
}

func frElementToBigInt(e fr.Element) big.Int {
    var ret big.Int 
	e.BigInt(&ret)

	return ret
}
