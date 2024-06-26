package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/cmakafui/sdm-golang-app/internal/sdm"
)

var tpl = template.Must(template.ParseFiles("web/templates/index.html"))

var memory *sdm.SDM

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := tpl.Execute(w, nil); err != nil {
			log.Printf("Error executing template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		address := r.FormValue("address")
		data := r.FormValue("data")
		iterationsStr := r.FormValue("iterations")

		iterations, err := strconv.Atoi(iterationsStr)
		if err != nil || iterations <= 0 {
			log.Printf("Invalid number of iterations: %s\n", iterationsStr)
			http.Error(w, "Iterations must be a positive integer", http.StatusBadRequest)
			return
		}

		if len(address) > memory.AddressSize() || len(data)*8 > memory.AddressSize() {
			log.Printf("Invalid address or data length: address length %d, data length %d\n", len(address), len(data)*8)
			http.Error(w, fmt.Sprintf("Address must be %d characters long and data must be %d characters long", memory.AddressSize(), memory.AddressSize()/8), http.StatusBadRequest)
			return
		}

		memAddress := sdm.EncodeTextToBinary(address, memory.AddressSize())
		memData := sdm.EncodeTextToBinary(data, memory.AddressSize())

		// Store the data in memory
		memory.Write(memAddress, memData)

		// Retrieve the data from memory using parallel processing
		retrievedData := memory.ReadWithIterationsParallel(memAddress, iterations)
		retrievedText := sdm.DecodeBinaryToText(retrievedData)

		response := map[string]string{
			"stored":    data,
			"retrieved": retrievedText,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

func generateRandomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// addressSize := uint(memory.AddressSize() / 8)
	address := gofakeit.UUID()
	data := gofakeit.Sentence(10)

	response := map[string]string{
		"address": address,
		"data":    data,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func clearMemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	memory.Clear()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Memory cleared")
}

func memoryStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	stats := memory.GetStats()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func memoryHistoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	history := memory.GetHistory()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

func main() {
	addressSize := 1000
	numAddresses := 10000
	memory = sdm.NewSDM(addressSize, numAddresses)

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/generate-random", generateRandomHandler)
	http.HandleFunc("/clear-memory", clearMemoryHandler)
	http.HandleFunc("/memory-stats", memoryStatsHandler)
	http.HandleFunc("/memory-history", memoryHistoryHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	log.Println("Starting server on :5080")
	if err := http.ListenAndServe(":5080", nil); err != nil {
		log.Fatal(err)
	}
}
