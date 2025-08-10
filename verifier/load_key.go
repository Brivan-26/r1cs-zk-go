package verifier 

import (
	"r1cs-zk-go/keys"
	"encoding/json"
	"io/ioutil"
)

func loadVerifyingKey() keys.VerifyingKey {
	var vk keys.VerifyingKey
	
	jsonData, err := ioutil.ReadFile("vk.json")
	if err != nil {
		panic("failed to read vk.json")
	}
	
	err = json.Unmarshal(jsonData, &vk)
	if err != nil {
		panic("failed to parse vk.json")
	}
	
	return vk
}

func loadProof() keys.Proof {
	var proof keys.Proof
	
	jsonData, err := ioutil.ReadFile("proof.json")
	if err != nil {
		panic("failed to read proof.json")
	}
	
	err = json.Unmarshal(jsonData, &proof)
	if err != nil {
		panic("failed to parse proof.json")
	}
	
	return proof
}