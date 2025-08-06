package main 

import (
	"github.com/cloudflare/circl/ecc/bls12381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"math/rand"
    "time"
)
func generateSRS(n1, n2, n int) ([]bls12381.G1, []bls12381.G2, []bls12381.G1){

	// ABSOLUTELY NOT SAFE RANDOM GENERATION
	rand.Seed(time.Now().UnixNano())
    r := rand.Intn(1000000000000) 
	
	var tau bls12381.Scalar 
	tau.SetUint64(uint64(r))

	ret1 := make([]bls12381.G1, n1)
	ret2 := make([]bls12381.G2, n1)
	ret3 := make([]bls12381.G1, n2)
	
	// Generate Omega
	point1 := bls12381.G1Generator()
	ret1[0] = *point1
	for i:=1; i < n1; i++ {
		var newPoint bls12381.G1
		newPoint.ScalarMult(&tau, point1)
		ret1[i] = newPoint
		point1 = &newPoint
	}

	// Generate Theta
	point2 := bls12381.G2Generator()
	ret2[0] = *point2
	for i:=1; i < n1; i++ {
		var newPoint bls12381.G2
		newPoint.ScalarMult(&tau, point2)
		ret2[i] = newPoint
		point2 = &newPoint
	}

	t_x := buildTx(n)
	var element fr.Element
	element.SetInt64(int64(r))
	t_tau := t_x.Eval(&element)

	var s bls12381.Scalar 
	s.SetUint64(t_tau.Uint64())

	// Generate Upsilon
	point1 = bls12381.G1Generator()
	var p bls12381.G1
	p.ScalarMult(&s, point1)
	ret3[0] = p 
	for i:=1; i < n2; i++ {
		var newPoint bls12381.G1 
		newPoint.ScalarMult(&tau, &p)
		ret3[i] = newPoint
		p = newPoint
	} 

	return ret1, ret2, ret3

}