package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type Wallets struct {
	Addresses []string `json:"wallets"`
}

func LoadWallets(filename string) (*Wallets, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var wallets Wallets
	if err := json.Unmarshal(bytes, &wallets); err != nil {
		return nil, err
	}

	return &wallets, nil
}

func Contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
