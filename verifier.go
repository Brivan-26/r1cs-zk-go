package main 

import (
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
)

func verifyProof(A, C curve.G1Affine, B curve.G2Affine) bool {
	_,_,_,gen2 := curve.Generators()
	P := make([]curve.G1Affine, 1)
	Q := make([]curve.G2Affine, 1)

	P[0] = A 
	Q[0] = B 
	e1, err1 := curve.Pair(P, Q)

	P[0] = C 
	Q[0] = gen2
	e2, err2 := curve.Pair(P, Q)
	
	if !e1.Equal(&e2) || err1 != nil || err2 != nil {
		return false
	}

	return true
}