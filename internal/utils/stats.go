package utils

import (
	"fmt"
	"math/big"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

// StatsData estrutura para compartilhar dados de estatísticas
type StatsData struct {
	CurrentKey       *big.Int
	KeysChecked      *int
	TotalKeysChecked int
	StartTime        time.Time
	SearchMode       int
	Direction        string
}

func ReportStats(data *StatsData, done chan bool) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			elapsedTime := time.Since(data.StartTime).Seconds()
			keysPerSecond := float64(*data.KeysChecked) / elapsedTime
			currentTotal := data.TotalKeysChecked + *data.KeysChecked

			// Formatação base
			baseInfo := fmt.Sprintf("[Current Key: 0x%s] || [Keys: %s] || [Keys/Seg: %s]",
				data.CurrentKey.Text(16),
				humanize.Comma(int64(currentTotal)),
				humanize.Comma(int64(keysPerSecond)))

			// Adiciona informação de direção apenas para o modo 4
			if data.SearchMode == 4 && data.Direction != "" {
				baseInfo += fmt.Sprintf(" || [Direction: %s]", data.Direction)
			}

			// Adiciona informação do modo de busca
			switch data.SearchMode {
			case 1:
				baseInfo += " || [Mode: Sequential]"
			case 2:
				baseInfo += " || [Mode: Random]"
			case 3:
				baseInfo += " || [Mode: Random+Sequential]"
			case 4:
				baseInfo += " || [Mode: Random+Bidirectional]"
			}

			fmt.Println(baseInfo)
		case <-done:
			return
		}
	}
}
func PrintResults(foundAddress *big.Int, walletAddress, privateKey string, balance float64, keysChecked int, startTime time.Time) {
	elapsedTime := time.Since(startTime).Seconds()
	keysPerSecond := float64(keysChecked) / elapsedTime
	color.White("=========================================================================\n")
	color.White("We Found the Target Wallet!!\n")
	color.Yellow("Private Key Hex: 0x%s\n", foundAddress.Text(16))
	color.Blue("Wallet Address: %s\n", walletAddress)
	color.Red("Private Key Wif: %s\n", privateKey)
	color.Green("Balance: %.12f BTC\n", balance)
	color.White("Keys: %s\n", humanize.Comma(int64(keysChecked)))
	color.White("Time: %.2f seconds\n", elapsedTime)
	color.White("Keys/Seg: %s\n", humanize.Comma(int64(keysPerSecond)))
	color.White("=========================================================================\n")
}
