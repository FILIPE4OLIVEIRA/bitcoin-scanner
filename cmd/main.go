package main

import (
	"bitcoin-scanner/internal/config"
	"bitcoin-scanner/internal/scanner"
	"bitcoin-scanner/internal/utils"
	"fmt"
	"log"
	"runtime"

	"github.com/fatih/color"
)

func main() {
	white := color.New(color.FgWhite).SprintFunc()

	// Load configuration
	ranges, err := config.LoadRanges("ranges.json")
	if err != nil {
		log.Fatalf("Failed to Load Ranges: %v", err)
	}

	wallets, err := config.LoadWallets("wallets.json")
	if err != nil {
		log.Fatalf("Failed to Load Wallets: %v", err)
	}

	// Get user input
	rangeNumber := utils.PromptRangeNumber(len(ranges.Ranges))
	searchMode := utils.PromptSearchMode()

	// Initialize scanner
	scnr := scanner.NewScanner(ranges, wallets, rangeNumber, searchMode)

	// Configurar número de workers CPU (metade como no original)
	numCPU := runtime.NumCPU() / 2
	runtime.GOMAXPROCS(numCPU * 2)
	fmt.Printf("CPUs: %s\n", white(numCPU))

	fmt.Printf("Key Ranges Initial [0x%s] -- Final [0x%s]\n",
		scnr.MinKeyInt.Text(16), scnr.MaxKeyInt.Text(16))

	// Start scanning
	foundAddress, keysChecked, startTime := scnr.Start()

	// Process results
	if foundAddress != nil {
		walletAddress := scnr.CreatePublicAddress(foundAddress)
		privateKey := scnr.PrivateKeyToWIF(foundAddress)
		balance := scnr.CheckBalance(walletAddress)

		utils.PrintResults(foundAddress, walletAddress, privateKey, balance, keysChecked, startTime)
		utils.SaveTargetWallet(walletAddress, privateKey, balance)
	} else {
		fmt.Println("Algoritmo Finalizado")
	}
}
