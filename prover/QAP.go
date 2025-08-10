package prover 

import (
	"r1cs-zk-go/utils"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"gitlab.com/oelmekki/matrix"
)

func R1CSToQAP(L, R, O, W matrix.Matrix) (polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial) {
	// safety check 
	matricesSanityChecks(L, R, O, W)
	
	u_s := utils.InterpolateFromMatrixCols(L)
	v_s := utils.InterpolateFromMatrixCols(R)
	w_s := utils.InterpolateFromMatrixCols(O)

	u_x := buildPolyFromInterpolationAndWitness(u_s, W)
	v_x := buildPolyFromInterpolationAndWitness(v_s, W)
	w_x := buildPolyFromInterpolationAndWitness(w_s, W)
	t_x := utils.BuildTx(L.Rows())
	h_x := buildHx(u_x, v_x, w_x, t_x)

	return u_x, v_x, w_x, t_x, h_x
}

func buildPolyFromInterpolationAndWitness(polys []polynomial.Polynomial, w matrix.Matrix) polynomial.Polynomial {
	witness := toFr(w)

	if len(polys) != len(witness) {
		panic("Malformed Polynomials")
	}

	var zeroElement fr.Element 
	zeroElement.SetZero()
	p_x := polynomial.Polynomial{zeroElement}
	for i:=0; i < len(polys); i++ {
		p_i := polys[i]
		p_i.ScaleInPlace(&witness[i])
		p_x.Add(p_x, p_i)
	}

	return p_x
}

func toFr(w matrix.Matrix) []fr.Element {
	if w.Cols() > 1 {
		// we shouldn't enter here due to previous checks, but safety checks
		panic("malformed witness!")
	}
	n := w.Rows()
	output := make([]fr.Element, n)
	for r:=0; r < n; r++ {
		val := w.At(r, 0)
		var e fr.Element 
		e.SetInt64(int64(val))
		output[r] = e
	}

	return output
}



// DividePolys returns quotient and remainder for (numerator)/(denominator)
func DividePolys(nom, den polynomial.Polynomial) (quot, rem polynomial.Polynomial) {
	// Check for zero denominator
	allZero := true
	for _, c := range den {
		if !c.IsZero() {
			allZero = false
			break
		}
	}
	if allZero {
		panic("Division by Zero Polynomial")
	}

	// Remove leading zeros from denominator
	den = normalizePoly(den)
	if len(nom) < len(den) {
		var zero fr.Element
		zero.SetZero()
		return polynomial.Polynomial{zero}, nom
	}

	// Copy numerator as remainder
	rem = make(polynomial.Polynomial, len(nom))
	copy(rem, nom)

	degreeDiff := len(nom) - len(den)
	quot = make(polynomial.Polynomial, degreeDiff+1)

	// Inverse of highest-order coeff of denominator
	var dInv fr.Element
	dInv.Inverse(&den[len(den)-1])

	for i := len(quot) - 1; i >= 0; i-- {
		var q fr.Element
		q.Mul(&rem[i+len(den)-1], &dInv)
		quot[i].Set(&q)

		// rem -= q * den shifted
		for j, n := range den {
			var prod fr.Element
			prod.Mul(&q, &n)
			rem[i+j].Sub(&rem[i+j], &prod)
		}
	}

	rem = normalizePoly(rem)
	return quot, rem
}

// normalizePoly removes redundant trailing zeros (highest-degree) from a polynomial
func normalizePoly(p polynomial.Polynomial) polynomial.Polynomial {
	i := len(p) - 1
	for ; i > 0; i-- {
		if !p[i].IsZero() {
			break
		}
	}
	return p[:i+1]
}

func buildHx(u_x, v_x, w_x, t_x polynomial.Polynomial) polynomial.Polynomial {
	uv_x := utils.MultiplyPolys(u_x, v_x)
	uvMinusw_x := utils.SubtractPolys(uv_x, w_x)
	// reminder should always be zero for valid witness, // TODO include check for this
	h_x, _ := DividePolys(uvMinusw_x, t_x)
	return h_x
}
