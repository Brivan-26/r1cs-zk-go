package main 

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"gitlab.com/oelmekki/matrix"
)

func R1CSToQAP(L, R, O, W matrix.Matrix) (polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial, polynomial.Polynomial) {
	// safety check 
	matricesSanityChecks(L, R, O, W)
	
	u_x := interpolateFromMatrix(L, W)
	v_x := interpolateFromMatrix(R, W)
	w_x := interpolateFromMatrix(O, W)
	t_x := buildTx(L.Rows())
	h_x := buildHx(u_x, v_x, w_x, t_x)

	return u_x, v_x, w_x, t_x, h_x
}

func interpolateFromMatrix(matrix, w matrix.Matrix) polynomial.Polynomial {
	witness := toFr(w)
	
	var zeroElement fr.Element 
	zeroElement.SetZero()
	p_x := polynomial.Polynomial{zeroElement}

	num_cols := matrix.Cols()
	for i:=0; i < num_cols; i++ {
		col := getColInFr(matrix, i)
		p_i := polynomial.InterpolateOnRange(col)
		p_i.ScaleInPlace(&witness[i])

		p_x.Add(p_x, p_i)
	}

	return p_x
}

func getColInFr(matrix matrix.Matrix, col int) []fr.Element {
	rows := matrix.Rows()
	output := make([]fr.Element, rows)

	for row:=0; row < rows; row++ {
		val := matrix.At(row, col)
		var e fr.Element
		e.SetInt64(int64(val))
		output[row] = e
	}

	return output
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

// MultiplyPolys multiplies two polynomials.
func MultiplyPolys(p1, p2 polynomial.Polynomial) polynomial.Polynomial {
    product := make(polynomial.Polynomial, len(p1)+len(p2)-1)
    // Initialize the polynomial to zero
    for i := range product {
        product[i].SetZero()
    }
    for i1, n1 := range p1 {
        for i2, n2 := range p2 {
            var mul fr.Element
            mul.Mul(&n1, &n2)
            product[i1+i2].Add(&product[i1+i2], &mul)
        }
    }
    return product
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

func buildTx(n int) polynomial.Polynomial {
	var oneElement, zeroElement fr.Element 
	oneElement.SetUint64(1)
	zeroElement.SetZero()

	t_x := polynomial.Polynomial{zeroElement, oneElement}
	for i:=1; i < n; i++ {
		var element fr.Element
		element.SetInt64(int64(-i))
		t_x = MultiplyPolys(t_x, polynomial.Polynomial{element, oneElement})
	}

	return t_x
}

func buildHx(u_x, v_x, w_x, t_x polynomial.Polynomial) polynomial.Polynomial {
	uv_x := MultiplyPolys(u_x, v_x)
	var uvMinusw_x polynomial.Polynomial
	uvMinusw_x.Sub(uv_x, w_x)
	// reminder should always be zero for valid witness, // TODO include check for this
	h_x, _ := DividePolys(uvMinusw_x, t_x)

	return h_x
}
