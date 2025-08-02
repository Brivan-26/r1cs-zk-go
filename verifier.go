package main 

import (
	"github.com/cloudflare/circl/ecc/bls12381"
)

func verifyProof(L, O []bls12381.G1, R []bls12381.G2) bool {
	if len(L) != len(R) && len(L) != len(O) {
		return false
	}
	
	gen2 := bls12381.G2Generator();

	for i:=0; i < len(L); i++ {
		e1 := bls12381.Pair(&L[i], &R[i])
		e2 := bls12381.Pair(&O[i], gen2)
		if !e1.IsEqual(e2) {
			return false
		}
	}

	return true
}