package utils

import (
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// ensureDataDir cria o diretório data se não existir
func ensureDataDir() error {
	err := os.MkdirAll("data", 0755)
	if err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}
	return nil
}

func SaveTestedKeys(privKeyInt *big.Int, walletNumber int) {
	// Garante que o diretório data existe
	if err := ensureDataDir(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	filename := filepath.Join("data", fmt.Sprintf("tested_keys_%d.txt", walletNumber))

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to Open File: %v\n", err)
		return
	}
	defer file.Close()

	privKeyHex := fmt.Sprintf("0x%s", privKeyInt.Text(16))
	_, err = file.WriteString(privKeyHex + "\n")
	if err != nil {
		fmt.Printf("Failed to Write to File: %v\n", err)
	}
}

func GetInitialKeysChecked(fileName string) int {
	// Verificar se o arquivo existe primeiro
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return 0
	}

	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return 0
	}

	lines := strings.Split(string(content), "\n")
	lineCount := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			lineCount++
		}
	}

	// TODOS os modos agora salvam a cada 2.500.000 chaves
	return lineCount * 2500000
}

func SaveTargetWallet(address, privateKey string, balance float64) {
	// Garante que o diretório data existe
	if err := ensureDataDir(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	filename := filepath.Join("data", "target_wallet.txt")

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to Open File: %v\n", err)
		return
	}
	defer file.Close()

	data := fmt.Sprintf("Wallet Address: %s\nPrivate Key: %s\nBalance: %.12f BTC\n\n", address, privateKey, balance)
	_, err = file.WriteString(data)
	if err != nil {
		fmt.Printf("Failed to Write to File: %v\n", err)
	}
}
