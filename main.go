package main 

import (
	"r1cs-zk-go/trusted_setup"
	"r1cs-zk-go/prover"
	"r1cs-zk-go/verifier"
	"fmt"
	"os"
)

func main() {

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
