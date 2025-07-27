// I want to prove that I know a number "x" such that "x^3 + 5x + 5 = y" where "y" is 155
// Our Arithmetic circuit is as follows:
//      v = x^2
//      y = x*v + 5x + 5
// To make our life easier, we convert to the following R1CS
//      v = x^2
//      y - 5x - 5 = x*v
// Our witness vector is then [1, y, v, x]
// We need to to convert to R1CS:
// | 0 0 1  0| |1| = |0 0 0 1| |1| * |0 0 0 1| |1|
// |-5 1 0 -5| |y|   |0 0 0 1| |y|   |0 0 1 0| |y|
//             |v|             |v|             |v|
//             |x|             |x|             |x|
package main

import (
	"github.com/cloudflare/circl/ecc/bls12381"
	"fmt"
)

func main() {

	// =============================== PROVER ===============================
	L11, L21, O11, O21, R11, R21 := prove()
	
	// prover sends L11, L21, R11, R21, O11, O21

	// =============================== VERIFIER ===============================

	// 1. e(L11, R11) =? e(O11, G2)
	// 2. e(L21, R21) =? e(O21, G2)

	if verify(&L11, &L21, &O11, &O21, &R11, &R21) {
		fmt.Println("Proof valid!")
	}else {
		panic("Invalid Proof!")
	}

}

func prove() (bls12381.G1, bls12381.G1, bls12381.G1, bls12381.G1, bls12381.G2, bls12381.G2) {
	// witness = a = [1, y, v, x] = [1, 155, 25, 5]

	// L11 = x
	// L21 = x

	L11 := mulG1(5, false)
	L21 := mulG1(5, false)

	// R11 = x
	// R21 = v
	var s2 bls12381.Scalar
	s2.SetUint64(25)

	R11 := mulG2(5, false)
	R21 := mulG2(25, false)

	// O11 = v = 25
	// O21 = -5 + y -5x = 125
		// O21 = -5G1 + 155G1 -5(5G1)

	var five, oneFiveFive bls12381.Scalar
	five.SetUint64(5)
	oneFiveFive.SetUint64(155)
	
	O11 := mulG1(25, false)
	fiveG1 := mulG1(5, false)
	oneFiveFiveG1 := mulG1(155, false)
	minusFiveG1 := mulG1(5, true)

	var O21 bls12381.G1 
	five.Neg()

	fiveG1.ScalarMult(&five, &fiveG1)

	O21.Add(&minusFiveG1, &oneFiveFiveG1)
	O21.Add(&O21, &fiveG1)

	return L11, L21, O11, O21, R11, R21

}

func verify(L11, L21, O11, O21 *bls12381.G1, R11, R21 *bls12381.G2) bool {
	gen2 := bls12381.G2Generator()

	e1 := bls12381.Pair(L11, R11)
	checkE1 := bls12381.Pair(O11, gen2)

	e2 := bls12381.Pair(L21, R21)
	checkE2 := bls12381.Pair(O21, gen2)

	return checkE1.IsEqual(e1) && checkE2.IsEqual(e2)

}

func mulG1(num uint64, signed bool) bls12381.G1 {
	var s bls12381.Scalar 
	s.SetUint64(num)
	if signed {
		s.Neg()
	}

	gen1 := bls12381.G1Generator()

	var point bls12381.G1 
	point.ScalarMult(&s, gen1)

	return point

}

func mulG2(num uint64, signed bool) bls12381.G2 {
	var s bls12381.Scalar 
	s.SetUint64(num)
	if signed {
		s.Neg()
	}

	gen2 := bls12381.G2Generator()

	var point bls12381.G2
	point.ScalarMult(&s, gen2)

	return point
}
