package main 

import (
	"r1cs-zk-go/trusted_setup"
	"r1cs-zk-go/prover"
	"r1cs-zk-go/verifier"
	"fmt"
	"os"
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

	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "setup":
		trusted_setup.GenerateSRS()
	case "prove":
		prover.Prove()
	case "verify":
		if !verifier.VerifyProof() {
			fmt.Println("Invalid Proof!")
		}else {
			fmt.Println("Valid Proof!")
		}
	default:
		fmt.Println("Unkown command")
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go build && ./r1cs-zk-go <command>")
	fmt.Println("")
	fmt.Println("Commands:")
	fmt.Println("  setup    Run trusted setup and generate proving/verifying keys")
	fmt.Println("  prove    Generate a Groth16 zk proof using 'pk.json' and save it to 'proof.json' file")
	fmt.Println("  verify   Verify a Groth16 zk proof from 'proof.json' and 'vk.json' file ")
	fmt.Println("")
	fmt.Println("Description:")
	fmt.Println("  This program implements a Groth16 zero-knowledge proof system")
	fmt.Println("")
}
