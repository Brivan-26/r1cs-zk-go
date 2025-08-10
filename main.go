package main 

import (
	"r1cs-zk-go/trusted_setup"
	"r1cs-zk-go/prover"
	"r1cs-zk-go/verifier"
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
	
	alpha, beta, gamma, teta, SRS1, SRS2, SRS3, psi, publicInputsSize := trusted_setup.GenerateSRS()

	proverPsi := psi[publicInputsSize:]
	verifierPsi := psi[:publicInputsSize]

	A, B, C := prover.Prove(SRS1, SRS3, SRS2, alpha, beta, proverPsi)

	if !verifier.VerifyProof(A, C, alpha, B, beta, gamma, teta, verifierPsi) {
		panic("Invalid proof!")
	}

	fmt.Println("Valid Proof!")
}

