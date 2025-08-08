package main 

import (
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"gitlab.com/oelmekki/matrix"
	// "math/rand"
    // "time"
	"math/big"
)
func generateSRS(n1, n2, n int, L, R, O matrix.Matrix) (curve.G1Affine, curve.G2Affine, []curve.G1Affine, []curve.G2Affine, []curve.G1Affine, []curve.G1Affine){
	_, _, g1Gen, g2Gen := curve.Generators()
	// ABSOLUTELY NOT SAFE RANDOM GENERATION
	// rand.Seed(time.Now().UnixNano())
	// tau := big.NewInt(int64(rand.Intn(1000000000000)))
	tau := big.NewInt(int64(10))


	omega := make([]curve.G1Affine, n1)
	theta := make([]curve.G2Affine, n1)
	upsilon := make([]curve.G1Affine, n2)
	
	// Generate Omega
	omega[0] = g1Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		omega[i] = newPoint
	}

	// Generate Theta
	theta[0] = g2Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G2Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		theta[i] = newPoint
	}

	t_x := buildTx(n)
	var element fr.Element
	element.SetBigInt(tau)
	t_tau := t_x.Eval(&element)

	var t_tau_bigInt big.Int 
	t_tau.BigInt(&t_tau_bigInt)

	var point curve.G1Affine
	point.ScalarMultiplicationBase(&t_tau_bigInt)
	upsilon[0] = point
	for i:=1; i < n2; i++ {
		tau_exp := new(big.Int).Exp(tau, big.NewInt(int64(i)), nil)
		tmp := new(big.Int).Mul(tau_exp, &t_tau_bigInt)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tmp)
		upsilon[i] = newPoint
	} 
	
	
	// ABSOLUTELY NOT SAFE RANDOM GENERATION
	// a := big.NewInt(int64(rand.Intn(10000))) 
	// b := big.NewInt(int64(rand.Intn(10000)))

	a := big.NewInt(int64(20)) 
	b := big.NewInt(int64(21))

	// generate alpha
	var alpha curve.G1Affine
	alpha.ScalarMultiplicationBase(a)
	
	// generate beta
	var beta curve.G2Affine
	beta.ScalarMultiplicationBase(b)
	
	// Generate Psi
	u_s := interpolateFromMatrixCols(L)
	v_s := interpolateFromMatrixCols(R)
	w_s := interpolateFromMatrixCols(O)

	psi := make([]curve.G1Affine, L.Cols())
	for i:=0; i < len(psi); i++ {
		u_i := u_s[i]
		v_i := v_s[i]
		w_i := w_s[i]

		v_tau := v_i.Eval(&element)
		u_tau := u_i.Eval(&element)
		w_tau := w_i.Eval(&element)

		var alpha_fr, beta_fr, mul1, mul2, tmp_sum, sum fr.Element
		alpha_fr.SetBigInt(a)
		beta_fr.SetBigInt(b)

		mul1.Mul(&alpha_fr, &v_tau)
		mul2.Mul(&beta_fr, &u_tau)

		tmp_sum.Add(&mul1, &mul2)
		sum.Add(&tmp_sum, &w_tau)

		scalar := frElementToBigInt(sum)
		var point curve.G1Affine
		point.ScalarMultiplicationBase(&scalar)
		
		psi[i] = point
	}

	return alpha, beta, omega, theta, upsilon, psi

}