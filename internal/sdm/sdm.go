package sdm

import (
	"bytes"
	"log"
	"math/rand"
	"sync"
)

// SDM struct represents the sparse distributed memory
type SDM struct {
	addressSize  int
	numAddresses int
	addresses    [][]byte
	counters     [][]int32
	history      []string
	historyIndex int
	mu           sync.RWMutex
	randSource   *rand.Rand
}

const maxHistorySize = 1000 // or any other reasonable limit

// NewSDM initializes a new sparse distributed memory
func NewSDM(addressSize, numAddresses int) *SDM {
	addresses := make([][]byte, numAddresses)
	counters := make([][]int32, numAddresses)
	history := make([]string, maxHistorySize)
	randSource := rand.New(rand.NewSource(rand.Int63()))

	for i := 0; i < numAddresses; i++ {
		addresses[i] = GenerateRandomBinaryVector(addressSize, randSource)
		counters[i] = make([]int32, addressSize)
	}

	return &SDM{
		addressSize:  addressSize,
		numAddresses: numAddresses,
		addresses:    addresses,
		counters:     counters,
		history:      history,
		randSource:   randSource,
	}
}

// GenerateRandomBinaryVector generates a random binary vector of a given size
func GenerateRandomBinaryVector(size int, randSource *rand.Rand) []byte {
	vector := make([]byte, size)
	for i := 0; i < size; i++ {
		vector[i] = byte(randSource.Intn(2) + '0')
	}
	return vector
}

// EncodeTextToBinary converts a string to a binary slice
func EncodeTextToBinary(text string, size int) []byte {
	data := []byte(text)
	binary := make([]byte, size)
	index := 0

	for i := 0; i < len(data); i++ {
		for j := 7; j >= 0; j-- {
			if index < size {
				if (data[i] & (1 << j)) > 0 {
					binary[index] = '1'
				} else {
					binary[index] = '0'
				}
				index++
			}
		}
	}

	for index < size {
		binary[index] = '0'
		index++
	}

	return binary
}

// DecodeBinaryToText converts a binary slice to a string
func DecodeBinaryToText(data []byte) string {
	buffer := bytes.NewBufferString("")
	for i := 0; i < len(data); i += 8 {
		var char byte
		for j := 0; j < 8; j++ {
			if data[i+j] == '1' {
				char += 1 << (7 - j)
			}
		}
		if char != 0 {
			buffer.WriteByte(char)
		}
	}
	return buffer.String()
}

// GenerateRandomASCIIString generates a random ASCII string of a given size
func GenerateRandomASCIIString(size int) string {
	b := make([]byte, size)
	for i := range b {
		b[i] = byte(rand.Intn(94) + 33) // ASCII range from '!' to '~'
	}
	return string(b)
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
	s.mu.Lock()
	s.history[s.historyIndex] = string(address)
	s.historyIndex = (s.historyIndex + 1) % maxHistorySize
	s.mu.Unlock()
	log.Println("Write operation completed.")
}

// ReadWithIterationsParallel retrieves data from the SDM at the given address using a specified number of convergence iterations with parallel processing
func (s *SDM) ReadWithIterationsParallel(address []byte, iterations int) []byte {
	retrieved := make([]byte, s.addressSize)
	previous := make([]byte, s.addressSize)
	copy(previous, retrieved)

	var pool = sync.Pool{
		New: func() interface{} {
			votes := make([]int32, s.addressSize)
			return &votes
		},
	}

	for iteration := 0; iteration < iterations; iteration++ {
		votes := make([]int32, s.addressSize)
		var wg sync.WaitGroup

		// Define the number of goroutines
		numGoroutines := 10
		addressesPerGoroutine := (len(s.addresses) + numGoroutines - 1) / numGoroutines

		// Function to be run by each goroutine
		voteWorker := func(start, end int) {
			defer wg.Done()
			localVotesPtr := pool.Get().(*[]int32)
			localVotes := *localVotesPtr
			for i := range localVotes {
				localVotes[i] = 0
			}

			for i := start; i < end; i++ {
				addr := s.addresses[i]
				if hammingDistance(address, addr) < s.addressSize/2 {
					for j := 0; j < s.addressSize; j++ {
						localVotes[j] += s.counters[i][j]
					}
				}
			}

			// Safely aggregate local votes to the global votes
			s.mu.Lock()
			for j := 0; j < s.addressSize; j++ {
				votes[j] += localVotes[j]
			}
			s.mu.Unlock()
			pool.Put(localVotesPtr)
		}

		// Launch goroutines
		for g := 0; g < numGoroutines; g++ {
			start := g * addressesPerGoroutine
			end := start + addressesPerGoroutine
			if end > len(s.addresses) {
				end = len(s.addresses)
			}
			wg.Add(1)
			go voteWorker(start, end)
		}

		// Wait for all goroutines to finish
		wg.Wait()

		// Determine the retrieved bits based on the votes
		for j := 0; j < s.addressSize; j++ {
			if votes[j] > 0 {
				retrieved[j] = '1'
			} else {
				retrieved[j] = '0'
			}
		}

		if bytes.Equal(previous, retrieved) {
			break
		}
		copy(previous, retrieved)
		address = retrieved
	}

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
	randomBits := make([]byte, s.numAddresses*s.addressSize)
	s.randSource.Read(randomBits)

	for i := 0; i < s.numAddresses; i++ {
		for j := 0; j < s.addressSize; j++ {
			s.addresses[i][j] = randomBits[i*s.addressSize+j]%2 + '0'
			s.counters[i][j] = 0
		}
	}

	s.mu.Lock()
	s.historyIndex = 0
	s.mu.Unlock()

	log.Println("Memory cleared.")
}

// GetStats returns memory stats and history of stored addresses
func (s *SDM) GetStats() map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := make([]string, maxHistorySize)
	copy(history, s.history)

	return map[string]interface{}{
		"totalAddresses": s.numAddresses,
		"history":        history,
	}
}

// GetHistory returns the history of stored addresses
func (s *SDM) GetHistory() []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := make([]string, maxHistorySize)
	copy(history, s.history)

	return history
}
