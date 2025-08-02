package main 

import (
	"gitlab.com/oelmekki/matrix"
	"github.com/cloudflare/circl/ecc/bls12381"
)

func prove(L, R, O matrix.Matrix) ([]bls12381.G1, []bls12381.G2, []bls12381.G1) {
	// witness for our specific problem in `main.go`:  w = [1, y, v, x] = [1, 155, 25, 5]
	W, _ := matrix.Build(
		matrix.Builder{
			matrix.Row{1},
			matrix.Row{155},
			matrix.Row{25},
			matrix.Row{5},
		},
	)

	// sanity checks
	if !matricesSanityChecks(L, R, O, W) {
		panic("Invalid Matrices!")
	}

	Lw, _ := L.DotProduct(W)
	Rw, _ := R.DotProduct(W)
	Ow, _ := O.DotProduct(W)

	// encrypt the values using ECC points
		// Lw will be encrypted in G1
		// Rw will be encrypted in G2
		// Ow will be encrypted in G1

	encrypted_Lw := make([]bls12381.G1, Lw.Rows())
	encrypted_Rw := make([]bls12381.G2, Rw.Rows())
	encrypted_Ow := make([]bls12381.G1, Ow.Rows())

	for i:=0; i < Lw.Rows(); i++ {
		val := uint64(Lw.At(i, 0))
		if val > 0 {
			encrypted_Lw[i] = mulG1(val, false)
		}else {
			encrypted_Lw[i] = mulG1(val, true)
		}
	}

	for i:=0; i < Rw.Rows(); i++ {
		val := uint64(Rw.At(i, 0))
		if val > 0 {
			encrypted_Rw[i] = mulG2(val, false)
		}else {
			encrypted_Rw[i] = mulG2(val, true)
		}
	}

	for i:=0; i < Ow.Rows(); i++ {
		val := uint64(Ow.At(i, 0))
		if val > 0 {
			encrypted_Ow[i] = mulG1(val, false)
		}else {
			encrypted_Ow[i] = mulG1(val, true)
		}
	}

	return encrypted_Lw, encrypted_Rw, encrypted_Ow
}

func matricesSanityChecks(L, R, O, W matrix.Matrix) bool {
	return (L.Rows() == R.Rows() && L.Rows() == O.Rows() && L.Cols() == R.Cols() && L.Cols() == O.Cols() && W.Cols() == 1 && W.Rows() == L.Cols())
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