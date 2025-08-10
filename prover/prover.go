package prover 

import (
	"r1cs-zk-go/keys"
	"r1cs-zk-go/r1cs"
	"r1cs-zk-go/witness"
	"r1cs-zk-go/utils"
	"gitlab.com/oelmekki/matrix"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"math/big"
	"fmt"
)

func Prove() {
	pkJSON := loadProvingKey()
	SRS1 := keys.JsonToG1AffineSlice(pkJSON.SRS1)
	SRS3 := keys.JsonToG1AffineSlice(pkJSON.SRS3)
	SRS2 := keys.JsonToG2AffineSlice(pkJSON.SRS2)
	alpha := keys.JsonToG1Affine(pkJSON.Alpha)
	beta := keys.JsonToG2Affine(pkJSON.Beta)
	psi := keys.JsonToG1AffineSlice(pkJSON.ProverPsi)
	
	L, R, O, err := r1cs.LoadR1CSFromJSON()
	if err != nil {
		panic(fmt.Sprintf("Failed to load R1CS: %v", err))
	}
	
	W, publicInputsSize, err := witness.LoadWitnessFromJSON()
	if err != nil {
		panic(fmt.Sprintf("Failed to load witness: %v", err))
	}

	// sanity checks
	if !matricesSanityChecks(L, R, O, W) {
		panic("Invalid Matrices!")
	}
	u_x, v_x, _, _, h_x := R1CSToQAP(L, R, O, W)

	
	// TODO add sanity checks on SRSs, that they power was generated successfully...
	
	A := EvalLAtSRS1(u_x, SRS1, alpha)
	B := EvalRAtSRS2(v_x, SRS2, beta)
	C := EvalOutputAtSRS13(psi, h_x, SRS3, W, publicInputsSize)	

	keys.BuildProof(A, C, B)
}

func matricesSanityChecks(L, R, O, W matrix.Matrix) bool {
	return (L.Rows() == R.Rows() && L.Rows() == O.Rows() && L.Cols() == R.Cols() && L.Cols() == O.Cols() && W.Cols() == 1 && W.Rows() == L.Cols())
}

func EvalLAtSRS1(u_x polynomial.Polynomial, srs []curve.G1Affine, alpha curve.G1Affine) curve.G1Affine {
	if len(u_x) != len(srs) {
		panic("Incorrect SRS")
	}

	coeffs := u_x
	var A curve.G1Affine
	for i:=0; i < len(coeffs); i++ {
		coeff := utils.FrElementToBigInt(coeffs[i])
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

	coeffs := v_x
	var B curve.G2Affine
	for i:=0; i < len(coeffs); i++ {
		coeff := utils.FrElementToBigInt(coeffs[i])
		var tmp curve.G2Affine 
		tmp.ScalarMultiplication(&srs[i], &coeff)
		B.Add(&B, &tmp)
	}

	B.Add(&B, &beta)

	return B
}

func EvalOutputAtSRS13(psi []curve.G1Affine, h_x polynomial.Polynomial, srs3 []curve.G1Affine, w matrix.Matrix, publicInputsSize int) curve.G1Affine {
	if len(psi) != (w.Rows() - publicInputsSize) {
		panic("Incorrect psi!")
	}
	var C curve.G1Affine 
	for i:=0; i < len(psi); i++ {
		a_i := big.NewInt(int64(w.At(publicInputsSize + i,0)))
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&psi[i], a_i)
		C.Add(&C, &tmp)
	}

	coeffs := h_x
	for i:=0; i < len(coeffs); i++ {
		coeff := utils.FrElementToBigInt(coeffs[i])
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&srs3[i], &coeff)
		C.Add(&C, &tmp)
	}

	return C
}


