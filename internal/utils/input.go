package utils

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func PromptRangeNumber(totalRanges int) int {
	reader := bufio.NewReader(os.Stdin)
	charReadline := '\n'

	if runtime.GOOS == "windows" {
		charReadline = '\r'
	}

	for {
		fmt.Printf("Escolha a Wallet (1 a %d): ", totalRanges)
		input, _ := reader.ReadString(byte(charReadline))
		input = strings.TrimSpace(input)
		rangeNumber, err := strconv.Atoi(input)
		if err == nil && rangeNumber >= 1 && rangeNumber <= totalRanges {
			return rangeNumber
		}
		fmt.Println("Numero Inválido.")
	}
}

func PromptSearchMode() int {
	reader := bufio.NewReader(os.Stdin)
	charReadline := '\n'

	if runtime.GOOS == "windows" {
		charReadline = '\r'
	}

	for {
		fmt.Print("Escolha o Modo de Busca (1 = Sequencial, 2 = Aleatório, 3 = Aleatório + Sequencial, 4 = Aleatório + Bidirecional): ")
		input, _ := reader.ReadString(byte(charReadline))
		input = strings.TrimSpace(input)
		searchMode, err := strconv.Atoi(input)
		if err == nil && (searchMode == 1 || searchMode == 2 || searchMode == 3 || searchMode == 4) {
			return searchMode
		}
		fmt.Println("Modo de Busca Inválido.")
	}
}
