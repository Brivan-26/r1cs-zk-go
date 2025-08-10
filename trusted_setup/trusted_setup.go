package trusted_setup 

import (
	"r1cs-zk-go/r1cs"
	"r1cs-zk-go/witness"
	"r1cs-zk-go/utils"
	curve "github.com/consensys/gnark-crypto/ecc/bls12-381"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"math/big"
	"fmt"
)
func GenerateSRS() (curve.G1Affine, curve.G2Affine, curve.G2Affine, curve.G2Affine, []curve.G1Affine, []curve.G2Affine, []curve.G1Affine, []curve.G1Affine, int){

	L, R, O, err := r1cs.LoadR1CSFromJSON()
	if err != nil {
		panic(fmt.Sprintf("Failed to load R1CS: %v", err))
	}

	// TODO Double check n1, n2, n inputs
	n1 := L.Rows()
	n2 := L.Rows() - 1
	n := L.Rows()
	
	_, _, g1Gen, g2Gen := curve.Generators()
	
	var element fr.Element
	element.MustSetRandom()
	tau := utils.FrElementToBigInt(element)


	omega := make([]curve.G1Affine, n1)
	theta := make([]curve.G2Affine, n1)
	upsilon := make([]curve.G1Affine, n2)
	
	// Generate Omega
	omega[0] = g1Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(&tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		omega[i] = newPoint
	}

	// Generate Theta
	theta[0] = g2Gen
	for i:=1; i < n1; i++ {
		tau_exp := new(big.Int).Exp(&tau, big.NewInt(int64(i)), nil)
		var newPoint curve.G2Affine
		newPoint.ScalarMultiplicationBase(tau_exp)
		theta[i] = newPoint
	}

	t_x := utils.BuildTx(n)
	
	t_tau := t_x.Eval(&element)

	var t_tau_bigInt big.Int 
	t_tau.BigInt(&t_tau_bigInt)

	var point curve.G1Affine
	point.ScalarMultiplicationBase(&t_tau_bigInt)
	upsilon[0] = point
	for i:=1; i < n2; i++ {
		tau_exp := new(big.Int).Exp(&tau, big.NewInt(int64(i)), nil)
		tmp := new(big.Int).Mul(tau_exp, &t_tau_bigInt)
		var newPoint curve.G1Affine
		newPoint.ScalarMultiplicationBase(tmp)
		upsilon[i] = newPoint
	} 

	var alpha_fr, beta_fr fr.Element 
	alpha_fr.MustSetRandom()
	beta_fr.MustSetRandom()

	a := utils.FrElementToBigInt(alpha_fr)
	b := utils.FrElementToBigInt(beta_fr)
	// generate alpha
	var alpha curve.G1Affine
	alpha.ScalarMultiplicationBase(&a)
	
	// generate beta
	var beta curve.G2Affine
	beta.ScalarMultiplicationBase(&b)

	// ABSOLUTELY NOT SAFE RANDOM GENERATION
	gamma := big.NewInt(int64(750)) 
	teta := big.NewInt(int64(2000))

	var gammaG curve.G2Affine
	gammaG.ScalarMultiplicationBase(gamma)

	var tetaG curve.G2Affine
	tetaG.ScalarMultiplicationBase(teta)

	publicInputs, err := witness.LoadPublicInputsFromJSON()
	if err != nil {
		panic(fmt.Sprintf("Could not load public input: %v", err))
	}

	// Generate Psi
	u_s := utils.InterpolateFromMatrixCols(L)
	v_s := utils.InterpolateFromMatrixCols(R)
	w_s := utils.InterpolateFromMatrixCols(O)

	psi := make([]curve.G1Affine, L.Cols())
	for i:=0; i < len(psi); i++ {
		u_i := u_s[i]
		v_i := v_s[i]
		w_i := w_s[i]

		v_tau := v_i.Eval(&element)
		u_tau := u_i.Eval(&element)
		w_tau := w_i.Eval(&element)

		var mul1, mul2, tmp_sum, sum, gamma_fr, teta_fr fr.Element
		gamma_fr.SetBigInt(gamma)
		teta_fr.SetBigInt(teta)

		mul1.Mul(&alpha_fr, &v_tau)
		mul2.Mul(&beta_fr, &u_tau)

		tmp_sum.Add(&mul1, &mul2)
		sum.Add(&tmp_sum, &w_tau)
		if i < len(publicInputs) {
			sum.Div(&sum, &gamma_fr)
		}else {
			sum.Div(&sum, &teta_fr)
		}

		scalar := utils.FrElementToBigInt(sum)
		var point curve.G1Affine
		point.ScalarMultiplicationBase(&scalar)
		
		psi[i] = point
	}

	return alpha, beta, gammaG, tetaG, omega, theta, upsilon, psi, len(publicInputs)

}