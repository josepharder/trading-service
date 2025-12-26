package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	KeyName   string
	KeySecret string
)

func main() {
	orders, err := FetchOrders()
	if err != nil {
		log.Fatalf("Failed to fetch orders: %v", err)
	}

	log.Printf("Fetched %d orders", len(orders))

	calculator := NewFIFOCalculator()
	trades, err := calculator.ProcessOrders(orders)
	if err != nil {
		log.Fatalf("Failed to process orders: %v", err)
	}

	log.Printf("Processed %d trades", len(trades))

	report, err := GenerateReport(trades)
	if err != nil {
		log.Fatalf("Failed to generate report: %v", err)
	}

	if err := WriteReport(report, "pnl_report.json"); err != nil {
		log.Fatalf("Failed to write report: %v", err)
	}

	log.Printf("Report written to pnl_report.json")
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

	KeyName = MustGetenv("KEY_NAME")
	KeySecret = MustGetenv("KEY_SECRET")
}

func MustGetenv(k string) string {
	v, ok := os.LookupEnv(k)
	if !ok {
		log.Panicf("%s environment variable not set.", k)
	}
	if v == "" {
		log.Printf("%s environment variable is empty\n", k)
	}
	return v
}
