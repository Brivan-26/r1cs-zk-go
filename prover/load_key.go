package prover 

import (
	"r1cs-zk-go/keys"
	"encoding/json"
	"io/ioutil"
	"os"
	"fmt"
)

func loadProvingKey() keys.ProvingKey {
	var pk keys.ProvingKey
	
	jsonData, err := ioutil.ReadFile("pk.json")
	if err != nil {
		fmt.Println("Could not read pk.json, make sure you have ran the trusted setup")
		os.Exit(1)
	}
	
	err = json.Unmarshal(jsonData, &pk)
	if err != nil {
		fmt.Println("failed to parse pk.json")
		os.Exit(1)
	}
	
	return pk
}
