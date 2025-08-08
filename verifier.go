package main 

import (
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func verifyProof(A, C, alpha curve.G1Affine, B, beta curve.G2Affine) bool {
	_,_,_,gen2 := curve.Generators()
	P := make([]curve.G1Affine, 1)
	Q := make([]curve.G2Affine, 1)

	P[0] = A 
	Q[0] = B 
	leftSide, err1 := curve.Pair(P, Q)

	P[0] = C 
	Q[0] = gen2
	e2, err2 := curve.Pair(P, Q)

	P[0] = alpha
	Q[0] = beta 
	e3, err3 := curve.Pair(P, Q)

	var rightSide curve.GT
	rightSide.Mul(&e3, &e2)

	if !(&leftSide).Equal(&rightSide) || err1 != nil || err2 != nil || err3 != nil {
		return false
	}

	return true
}