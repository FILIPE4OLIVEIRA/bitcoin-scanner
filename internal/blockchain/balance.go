package blockchain

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func CheckBalance(address string) float64 {
	const retries = 2
	const delay = 3 * time.Second

	client := &http.Client{Timeout: 5 * time.Second}
	for attempt := 0; attempt < retries; attempt++ {
		resp, err := client.Get(fmt.Sprintf("https://blockchain.info/balance?active=%s", address))
		if err != nil {
			if attempt < retries-1 {
				log.Printf("Error Checking Balance, Retrying in %v: %v", delay, err)
				time.Sleep(delay)
				continue
			} else {
				log.Printf("Error Checking Balance: %v", err)
				return 0
			}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Non-OK HTTP Status: %s", resp.Status)
			return 0
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Printf("Error Reading Response Body: %v", err)
			return 0
		}

		var result map[string]map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			log.Printf("Error Unmarshalling JSON: %v", err)
			return 0
		}

		finalBalance := result[address]["final_balance"].(float64)
		return finalBalance / 100000000 // Convert from satoshis to bitcoins
	}
	return 0
}
