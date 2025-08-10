package prover 

import (
	"r1cs-zk-go/keys"
	"encoding/json"
	"io/ioutil"
)

func loadProvingKey() keys.ProvingKey {
	var pk keys.ProvingKey
	
	jsonData, err := ioutil.ReadFile("pk.json")
	if err != nil {
		panic("failed to read pk.json")
	}
	
	err = json.Unmarshal(jsonData, &pk)
	if err != nil {
		panic("failed to parse pk.json")
	}
	
	return pk
}