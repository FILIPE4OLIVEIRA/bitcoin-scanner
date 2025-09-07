package scanner

import (
	"bitcoin-scanner/internal/config"
	"math/big"
	"sync"
)

func (s *Scanner) cpuworker(privKeyChan <-chan *big.Int, resultChan chan<- *big.Int, wg *sync.WaitGroup) {
	defer wg.Done()

	for privKeyInt := range privKeyChan {
		address := s.CreatePublicAddress(privKeyInt)
		if config.Contains(s.Wallets.Addresses, address) {
			select {
			case resultChan <- privKeyInt:
				return
			default:
				return
			}
		}
	}
}
