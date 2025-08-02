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

	l, r, o := prove(L, R, O)

	if !verifyProof(l, o, r) {
		panic("Invalid proof!")
	}

	fmt.Println("Valid Proof!")
}

