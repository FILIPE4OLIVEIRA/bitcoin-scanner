package scanner

import (
	"bitcoin-scanner/internal/blockchain"
	"bitcoin-scanner/internal/config"
	"bitcoin-scanner/internal/crypto"
	"bitcoin-scanner/internal/utils"
	"fmt"
	"math/big"
	"math/rand"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type Scanner struct {
	Ranges      *config.Ranges
	Wallets     *config.Wallets
	RangeNumber int
	SearchMode  int
	MinKeyInt   *big.Int
	MaxKeyInt   *big.Int
	PrivKeyInt  *big.Int
	CurrentKey  *big.Int
}

func NewScanner(ranges *config.Ranges, wallets *config.Wallets, rangeNumber, searchMode int) *Scanner {
	minKeyHex := ranges.Ranges[rangeNumber-1].Min
	maxKeyHex := ranges.Ranges[rangeNumber-1].Max

	minKeyInt := new(big.Int)
	minKeyInt.SetString(minKeyHex[2:], 16)

	maxKeyInt := new(big.Int)
	maxKeyInt.SetString(maxKeyHex[2:], 16)

	return &Scanner{
		Ranges:      ranges,
		Wallets:     wallets,
		RangeNumber: rangeNumber,
		SearchMode:  searchMode,
		MinKeyInt:   minKeyInt,
		MaxKeyInt:   maxKeyInt,
		PrivKeyInt:  new(big.Int).Set(minKeyInt),
		CurrentKey:  new(big.Int).Set(minKeyInt),
	}
}

func (s *Scanner) Start() (*big.Int, int, time.Time) {
	fileName := filepath.Join("data", fmt.Sprintf("tested_keys_%d.txt", s.RangeNumber))
	totalKeysChecked := utils.GetInitialKeysChecked(fileName)
	keysChecked := 0
	startTime := time.Now()

	privKeyChan := make(chan *big.Int, 10000)
	resultChan := make(chan *big.Int, 1) // Adicionar buffer de 1
	var wg sync.WaitGroup

	// Start workers
	numCPU := runtime.NumCPU() / 2
	for i := 0; i < numCPU*2; i++ {
		wg.Add(1)
		go s.cpuworker(privKeyChan, resultChan, &wg)
	}

	// Preparar dados para estatísticas
	statsData := &utils.StatsData{
		CurrentKey:       s.CurrentKey,
		KeysChecked:      &keysChecked,
		TotalKeysChecked: totalKeysChecked,
		StartTime:        startTime,
		SearchMode:       s.SearchMode,
		Direction:        "",
	}

	// Start stats reporter
	done := make(chan bool)
	go utils.ReportStats(statsData, done)

	// Start key generator com controle de parada
	stopGenerator := make(chan bool)
	go s.keyGenerator(privKeyChan, &keysChecked, statsData, stopGenerator)

	foundAddress := <-resultChan

	// Parar o generator e estatísticas
	stopGenerator <- true
	close(done)

	// Fechar canal de workers
	close(privKeyChan)

	// Esperar workers terminarem
	wg.Wait()

	return foundAddress, keysChecked, startTime
}

func (s *Scanner) keyGenerator(privKeyChan chan<- *big.Int, keysChecked *int, statsData *utils.StatsData, stopChan <-chan bool) {
	sequentialCount := 0
	var direction int = 1
	var baseKey *big.Int
	var directionStr string
	const maxDistance = 10000000

	for {
		select {
		case <-stopChan:
			return // Para imediatamente quando receber sinal de parada
		default:
			var currentKey *big.Int

			switch s.SearchMode {
			case 1: // Sequencial
				currentKey = new(big.Int).Set(s.PrivKeyInt)
				s.PrivKeyInt.Add(s.PrivKeyInt, big.NewInt(1))
				if *keysChecked%2500000 == 0 {
					utils.SaveTestedKeys(currentKey, s.RangeNumber)
				}

			case 2: // Aleatório
				currentKey = s.getRandomKeyInRange()
				if *keysChecked%2500000 == 0 {
					utils.SaveTestedKeys(currentKey, s.RangeNumber)
				}

			case 3: // Aleatório + Sequencial
				if sequentialCount == 0 || sequentialCount >= 25000000 {
					s.PrivKeyInt = s.getRandomKeyInRange()
					sequentialCount = 0
				}
				currentKey = new(big.Int).Set(s.PrivKeyInt)
				s.PrivKeyInt.Add(s.PrivKeyInt, big.NewInt(1))
				sequentialCount++

				if *keysChecked%2500000 == 0 {
					utils.SaveTestedKeys(currentKey, s.RangeNumber)
				}

			case 4: // Aleatório + Bidirecional
				if sequentialCount == 0 {
					baseKey = s.getRandomKeyInRange()
					currentKey = new(big.Int).Set(baseKey)
					direction = 1
					sequentialCount = 1
					directionStr = "Forward"
				} else {
					if direction == 1 {
						// Continua indo para frente
						currentKey = new(big.Int).Add(baseKey, big.NewInt(int64(sequentialCount)))
						sequentialCount++
						directionStr = "Forward"

						// Se atingiu o limite máximo (10 milhões), inverte a direção
						if sequentialCount > maxDistance {
							direction = -1
							sequentialCount = 1
						}
					} else {
						// Agora vai para trás
						currentKey = new(big.Int).Sub(baseKey, big.NewInt(int64(sequentialCount)))
						sequentialCount++
						directionStr = "Backward"

						// Se atingiu o limite máximo (10 milhões), reinicia com nova base
						if sequentialCount > maxDistance {
							sequentialCount = 0
							continue
						}
					}
				}

				statsData.Direction = directionStr

				if *keysChecked%2500000 == 0 {
					utils.SaveTestedKeys(currentKey, s.RangeNumber)
				}
			}

			// Verifica se a chave está dentro do range válido
			if currentKey.Cmp(s.MinKeyInt) < 0 || currentKey.Cmp(s.MaxKeyInt) > 0 {
				if s.SearchMode == 4 {
					sequentialCount = 0
				}
				continue
			}

			s.CurrentKey.Set(currentKey)

			// Envia para o canal com verificação de fechamento
			select {
			case privKeyChan <- currentKey:
				*keysChecked++
			case <-stopChan:
				return
			}
		}
	}
}

func (s *Scanner) getRandomKeyInRange() *big.Int {
	rangeInt := new(big.Int).Sub(s.MaxKeyInt, s.MinKeyInt)
	randInt := new(big.Int).Rand(rand.New(rand.NewSource(time.Now().UnixNano())), rangeInt)
	randInt.Add(randInt, s.MinKeyInt)
	return randInt
}

func (s *Scanner) CreatePublicAddress(privKeyInt *big.Int) string {
	return crypto.CreatePublicAddress(privKeyInt)
}

func (s *Scanner) PrivateKeyToWIF(privKeyInt *big.Int) string {
	return crypto.PrivateKeyToWIF(privKeyInt)
}

func (s *Scanner) CheckBalance(address string) float64 {
	return blockchain.CheckBalance(address)
}
