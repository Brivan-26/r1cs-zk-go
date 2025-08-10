package utils 

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr"
	"github.com/consensys/gnark-crypto/ecc/bls12-381/fr/polynomial"
	"gitlab.com/oelmekki/matrix"
	"math/big"
)

func FrElementToBigInt(e fr.Element) big.Int {
    var ret big.Int 
	e.BigInt(&ret)

	return ret
}

func BuildTx(n int) polynomial.Polynomial {
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

func InterpolateFromMatrixCols(matrix matrix.Matrix) []polynomial.Polynomial {
	num_cols := matrix.Cols()
	ret := make([]polynomial.Polynomial, num_cols)

	for i:=0; i < num_cols; i++ {
		col := getColInFr(matrix, i)
		p_i := polynomial.InterpolateOnRange(col)
		ret[i] = p_i		
	}

	return ret
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
