package main 

import (
	"gitlab.com/oelmekki/matrix"
	"fmt"
)

func main() {

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

	// O =  | 0 0 1  0| 
	//      |-5 1 0 -5|

	// L =  |0 0 0 1| 
	//      |0 0 0 1|   
	
	// R = |0 0 0 1|
	//     |0 0 1 0|
	
	L, _ := matrix.Build(
		matrix.Builder{
			matrix.Row{0, 0, 0, 1},
			matrix.Row{0, 0, 0, 1},
		},
	)
	R, _ := matrix.Build(
		matrix.Builder{
			matrix.Row{0, 0, 0, 1},
			matrix.Row{0, 0, 1, 0},
		},
	)
	O, _ := matrix.Build(
		matrix.Builder{
			matrix.Row{0, 0, 1, 0},
			matrix.Row{-5, 1, 0, -5},
		},
	)

	alpha, beta, SRS1, SRS2, SRS3, psi := generateSRS(L.Rows(), L.Rows()-1, L.Rows(), L, R, O) // TODO Double check n1, n2, n3 inputs

	A, B, C := prove(L, R, O, SRS1, SRS3, SRS2, alpha, beta, psi)

	if !verifyProof(A, C, alpha, B, beta) {
		panic("Invalid proof!")
	}

	fmt.Println("Valid Proof!")
}

