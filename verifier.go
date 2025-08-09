package main 

import (
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"fmt"
	"math/big"
)

func verifyProof(A, C, alpha curve.G1Affine, B, beta, gamma, teta curve.G2Affine, psi []curve.G1Affine) bool {
	publicInputs, err := LoadPublicInputsFromJSON()
	if err != nil {
		panic(fmt.Sprintf("Failed to load public witness: %v", err))
	}

	X := calculateX(psi, publicInputs)

	P := make([]curve.G1Affine, 1)
	Q := make([]curve.G2Affine, 1)

	P[0] = A 
	Q[0] = B 
	leftSide, err1 := curve.Pair(P, Q)

	P[0] = C 
	Q[0] = teta
	e2, err2 := curve.Pair(P, Q)

	P[0] = alpha
	Q[0] = beta 
	e3, err3 := curve.Pair(P, Q)

	P[0] = X
	Q[0] = gamma 
	e4, err4 := curve.Pair(P, Q)

	var rightSide curve.GT
	rightSide.Mul(&e3, &e2)
	rightSide.Mul(&rightSide, &e4)

	if !(&leftSide).Equal(&rightSide) || err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return false
	}

	return true
}

func calculateX(psi []curve.G1Affine, publicInputs []int) curve.G1Affine {
	if len(psi) != len(publicInputs) {
		panic("Missmatch public witness")
	}

	var X curve.G1Affine 
	for i:=0; i < len(publicInputs); i++ {
		a_i := big.NewInt(int64(publicInputs[i]))
		var tmp curve.G1Affine
		tmp.ScalarMultiplication(&psi[i], a_i)
		X.Add(&X, &tmp)
	}

	return X
}