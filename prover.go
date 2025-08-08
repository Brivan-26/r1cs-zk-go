package main 

import (
	"gitlab.com/oelmekki/matrix"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"math/big"
)

func prove(L, R, O matrix.Matrix, SRS1, SRS3 []curve.G1Affine, SRS2 []curve.G2Affine, alpha curve.G1Affine, beta curve.G2Affine, psi []curve.G1Affine) (curve.G1Affine, curve.G2Affine, curve.G1Affine) {
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
	u_x, v_x, _, _, h_x := R1CSToQAP(L, R, O, W)

	
	// TODO add sanity checks on SRSs, that they power was generated successfully...
	
	A := EvalLAtSRS1(u_x, SRS1, alpha)
	B := EvalRAtSRS2(v_x, SRS2, beta)
	C := EvalOutputAtSRS13(psi, h_x, SRS3, W)	

	return A, B, C
}

func matricesSanityChecks(L, R, O, W matrix.Matrix) bool {
	return (L.Rows() == R.Rows() && L.Rows() == O.Rows() && L.Cols() == R.Cols() && L.Cols() == O.Cols() && W.Cols() == 1 && W.Rows() == L.Cols())
}

func EvalLAtSRS1(u_x polynomial.Polynomial, srs []curve.G1Affine, alpha curve.G1Affine) curve.G1Affine {
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

	A.Add(&A, &alpha)

	return A
}

func EvalRAtSRS2(v_x polynomial.Polynomial, srs []curve.G2Affine, beta curve.G2Affine) curve.G2Affine {
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

	B.Add(&B, &beta)

	return B
}

func EvalOutputAtSRS13(psi []curve.G1Affine, h_x polynomial.Polynomial, srs3 []curve.G1Affine, w matrix.Matrix) curve.G1Affine {
	if len(psi) != w.Rows() {
		panic("Incorrect psi!")
	}
	var C curve.G1Affine 
	a_i := big.NewInt(int64(w.At(0,0)))
	C.ScalarMultiplication(&psi[0], a_i)
	
	for i:= 1; i < len(psi); i++ {
		a_i = big.NewInt(int64(w.At(i,0)))
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&psi[i], a_i)
		C.Add(&C, &tmp)
	}

	coeffs := h_x
	for i:=0; i < len(coeffs); i++ {
		coeff := frElementToBigInt(coeffs[i])
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
