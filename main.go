package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

var receiptsStore = make(map[string]int)

func calculatePoints(receipt Receipt) int {
	points := 0

	retailerName := receipt.Retailer
	regex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(regex.FindAllString(retailerName, -1))

	total, _ := strconv.ParseFloat(receipt.Total, 64)
	if total == float64(int(total)) {
		points += 50
	}

	if math.Mod(total, 0.25) == 0 {
		points += 25
	}

	numItems := len(receipt.Items)
	points += (numItems / 2) * 5

	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription)%3 == 0 {
			price, _ := strconv.ParseFloat(item.Price, 64)
			points += int(math.Ceil(price * 0.2))
		}
	}

	purchaseDate, _ := time.Parse("2024-09-18", receipt.PurchaseDate)
	if purchaseDate.Day()%2 != 0 {
		points += 6
	}

	purchaseTime, _ := time.Parse("09:00", receipt.PurchaseTime)
	if purchaseTime.Hour() == 14 || (purchaseTime.Hour() == 15 && purchaseTime.Minute() == 0) {
		points += 10
	}

	return points
}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiptID := uuid.New().String()
	points := calculatePoints(receipt)
	receiptsStore[receiptID] = points

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"id": receiptID})
}

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
