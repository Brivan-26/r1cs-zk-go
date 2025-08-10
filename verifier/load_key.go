package verifier 

import (
	"r1cs-zk-go/keys"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
)

func loadVerifyingKey() keys.VerifyingKey {
	var vk keys.VerifyingKey
	
	jsonData, err := ioutil.ReadFile("vk.json")
	if err != nil {
		fmt.Println("Could not read vk.json, make sure you have ran the trusted setup")
		os.Exit(1)
	}
	
	err = json.Unmarshal(jsonData, &vk)
	if err != nil {
		fmt.Println("failed to parse vk.json")
		os.Exit(1)
	}
	
	return vk
}

func loadProof() keys.Proof {
	var proof keys.Proof
	
	jsonData, err := ioutil.ReadFile("proof.json")
	if err != nil {
		fmt.Println("Could not read proof.json, make sure you have ran the prove command")
		os.Exit(1)
	}
	
	err = json.Unmarshal(jsonData, &proof)
	if err != nil {
		fmt.Println("failed to parse proof.json")
		os.Exit(1)
	}
	
	return proof
}