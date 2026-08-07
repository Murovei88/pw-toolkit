package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/murovei88/pw-toolkit/internal/calc"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Usage: go run cmd/data-parser/main.go <path_to_xml>")
	}

	xmlPath := os.Args[1]
	data, err := os.ReadFile(xmlPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	stats, err := calc.CalculateStats(data)
	if err != nil {
		log.Fatalf("Calculation error: %v", err)
	}

	fmt.Println("✅ Successfully calculated stats!")
	
	// Красивый вывод
	jsonData, _ := json.MarshalIndent(stats, "", "  ")
	fmt.Println(string(jsonData))
}
