package main 

import (
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"math/rand"
    "time"
	"math/big"
)
func generateSRS(n1, n2, n int) ([]curve.G1Affine, []curve.G2Affine, []curve.G1Affine){
	_, _, g1Gen, g2Gen := curve.Generators()
	// ABSOLUTELY NOT SAFE RANDOM GENERATION
	rand.Seed(time.Now().UnixNano())
	tau := big.NewInt(int64(rand.Intn(1000000000000)))


	ret1 := make([]curve.G1Affine, n1)
	ret2 := make([]curve.G2Affine, n1)
	ret3 := make([]curve.G1Affine, n2)
	
	// Generate Omega
	ret1[0] = g1Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		ret1[i] = newPoint
	}

	// Generate Theta
	ret2[0] = g2Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G2Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		ret2[i] = newPoint
	}

	t_x := buildTx(n)
	var element fr.Element
	element.SetBigInt(tau)
	t_tau := t_x.Eval(&element)

	var t_tau_bigInt big.Int 
	t_tau.BigInt(&t_tau_bigInt)

	var point curve.G1Affine
	point.ScalarMultiplicationBase(&t_tau_bigInt)
	ret3[0] = point
	for i:=1; i < n2; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		tmp := new(big.Int).Mul(tau_exp, &t_tau_bigInt)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tmp)
		ret3[i] = newPoint
	} 

	return ret1, ret2, ret3

}