package main 

import (
	"github.com/cloudflare/circl/ecc/bls12381"
)

func verifyProof(A, C bls12381.G1, B bls12381.G2) bool {

	gen2 := bls12381.G2Generator();

	e1 := bls12381.Pair(&A, &B)
	e2 := bls12381.Pair(&C, gen2)
	
	if !e1.IsEqual(e2) {
		return false
	}

	return true
}