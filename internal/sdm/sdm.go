package sdm

import (
	"log"
	"math/rand"
)

// SDM struct represents the sparse distributed memory
type SDM struct {
	addressSize  int
	numAddresses int
	addresses    [][]byte
	counters     [][]int
	history      []string
}

// NewSDM initializes a new sparse distributed memory
func NewSDM(addressSize, numAddresses int) *SDM {
	addresses := make([][]byte, numAddresses)
	counters := make([][]int, numAddresses)
	history := []string{}

	for i := 0; i < numAddresses; i++ {
		addresses[i] = GenerateRandomBinaryVector(addressSize)
		counters[i] = make([]int, addressSize)
	}

	return &SDM{
		addressSize:  addressSize,
		numAddresses: numAddresses,
		addresses:    addresses,
		counters:     counters,
		history:      history,
	}
}

// GenerateRandomBinaryVector generates a random binary vector of a given size
func GenerateRandomBinaryVector(size int) []byte {
	vector := make([]byte, size)
	for i := 0; i < size; i++ {
		vector[i] = byte(rand.Intn(2) + '0')
	}
	return vector
}

// Write stores data in the SDM at the given address
func (s *SDM) Write(address, data []byte) {
	log.Printf("Writing to address: %s, data: %s\n", string(address), string(data))
	for i, addr := range s.addresses {
		if hammingDistance(address, addr) < s.addressSize/2 {
			for j := 0; j < s.addressSize; j++ {
				if data[j] == '1' {
					s.counters[i][j]++
				} else {
					s.counters[i][j]--
				}
			}
		}
	}
	s.history = append(s.history, string(address))
	log.Println("Write operation completed.")
}

// Read retrieves data from the SDM at the given address using convergence
func (s *SDM) Read(address []byte) []byte {
	return s.ReadWithIterations(address, 10)
}

// ReadWithIterations retrieves data from the SDM at the given address using a specified number of convergence iterations
func (s *SDM) ReadWithIterations(address []byte, iterations int) []byte {
	retrieved := make([]byte, s.addressSize)

	for iteration := 0; iteration < iterations; iteration++ {
		votes := make([]int, s.addressSize)
		for i, addr := range s.addresses {
			if hammingDistance(address, addr) < s.addressSize/2 {
				for j := 0; j < s.addressSize; j++ {
					votes[j] += s.counters[i][j]
				}
			}
		}
		for j := 0; j < s.addressSize; j++ {
			if votes[j] > 0 {
				retrieved[j] = '1'
			} else {
				retrieved[j] = '0'
			}
		}
		log.Printf("Iteration %d: retrieved data: %s\n", iteration, string(retrieved))
		address = retrieved
	}

	log.Printf("Final retrieved data: %s\n", string(retrieved))
	return retrieved
}

// hammingDistance calculates the Hamming distance between two binary vectors
func hammingDistance(a, b []byte) int {
	distance := 0
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			distance++
		}
	}
	return distance
}

// AddressSize returns the size of the address in the SDM
func (s *SDM) AddressSize() int {
	return s.addressSize
}

// Clear clears the memory
func (s *SDM) Clear() {
	s.addresses = make([][]byte, s.numAddresses)
	s.counters = make([][]int, s.numAddresses)
	s.history = []string{}
	for i := 0; i < s.numAddresses; i++ {
		s.addresses[i] = GenerateRandomBinaryVector(s.addressSize)
		s.counters[i] = make([]int, s.addressSize)
	}
	log.Println("Memory cleared.")
}

// GetStats returns memory stats and history of stored addresses
func (s *SDM) GetStats() map[string]interface{} {
	return map[string]interface{}{
		"totalAddresses": s.numAddresses,
		"history":        s.history,
	}
}

// GetHistory returns the history of stored addresses
func (s *SDM) GetHistory() []string {
	return s.history
}