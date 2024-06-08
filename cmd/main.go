package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/cmakafui/sdm-golang-app/internal/sdm"
)

var tpl = template.Must(template.ParseFiles("web/templates/index.html"))

var memory *sdm.SDM

func homeHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling request")
	if r.Method == http.MethodGet {
		log.Println("GET request received")
		err := tpl.Execute(w, nil)
		if err != nil {
			log.Printf("Error executing template: %v\n", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
	} else if r.Method == http.MethodPost {
		log.Println("POST request received")
		address := r.FormValue("address")
		data := r.FormValue("data")
		iterationsStr := r.FormValue("iterations")
		iterations, err := strconv.Atoi(iterationsStr)
		if err != nil || iterations <= 0 {
			log.Printf("Invalid number of iterations: %s\n", iterationsStr)
			http.Error(w, "Iterations must be a positive integer", http.StatusBadRequest)
			return
		}

		log.Printf("Received address: %s, data: %s, iterations: %d\n", address, data, iterations)

		if len(address) != memory.AddressSize() || len(data) != memory.AddressSize() {
			log.Printf("Invalid address or data length: address length %d, data length %d\n", len(address), len(data))
			http.Error(w, fmt.Sprintf("Address and data must be %d characters long", memory.AddressSize()), http.StatusBadRequest)
			return
		}

		if address != "" && data != "" {
			memAddress := []byte(address)
			memData := []byte(data)
			log.Println("Writing data to memory")
			memory.Write(memAddress, memData)
		}

		retrievedData := memory.ReadWithIterations([]byte(address), iterations)
		log.Printf("Retrieved data: %s\n", string(retrievedData))
		w.Header().Set("Content-Type", "text/plain")
		w.Write(retrievedData)
	}
}

func generateRandomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	address := sdm.GenerateRandomBinaryVector(memory.AddressSize())
	data := sdm.GenerateRandomBinaryVector(memory.AddressSize())
	response := map[string]string{
		"address": string(address),
		"data":    string(data),
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