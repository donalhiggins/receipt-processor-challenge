package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total,string"`
}

type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price,string"`
}

type ReceiptID struct {
	ID      string
	Receipt Receipt
}

var receipts = make(map[string]Receipt)
var receiptsLock sync.RWMutex

func main() {
	http.HandleFunc("/receipts/process", processReceipt)
	http.HandleFunc("/receipts/", getPoints)
	http.ListenAndServe(":8080", nil)

}

func processReceipt(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var rec Receipt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rec)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Generate UUID.
	id := uuid.New()

	receiptsLock.Lock()
	receipts[id.String()] = rec
	receiptsLock.Unlock()

	response := map[string]string{
		"id": id.String(),
	}
	json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) != 4 {
		http.Error(w, "Invalid path", http.StatusBadRequest)
		return
	}
	id := parts[2]

	receiptsLock.RLock()
	rec, ok := receipts[id]
	receiptsLock.RUnlock()

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	totalPoints := calculatePoints(rec)

	response := map[string]int{
		"points": totalPoints,
	}
	json.NewEncoder(w).Encode(response)
}

func calculatePoints(rec Receipt) int {
	totalPoints := 0

	// Add one point for every alphanumeric character in the retailer name.
	var alphanumeric = regexp.MustCompile("^[a-zA-Z0-9_]*$")

	for _, char := range rec.Retailer {
		if alphanumeric.MatchString(string(char)) {
			totalPoints += 1
		}
	}

	// Check if total is a round amount.
	if rec.Total-float64(int(rec.Total)) == 0 {
		totalPoints += 50
	}

	// Check if total is a multiple of 0.25.
	if (rec.Total*4)-float64(int(rec.Total*4)) == 0 {
		totalPoints += 25
	}

	// Add 5 points for each two items on receipt.
	totalPoints += (len(rec.Items) / 2) * 5

	// If trimmed length of item description is mutiple of three, multiply
	// item price by 0.2 and round up to nearest integer then add to totalPoints.
	for i := 0; i < len(rec.Items); i++ {

		if len(strings.TrimSpace(rec.Items[i].ShortDescription))%3 == 0 {
			price := rec.Items[i].Price
			price *= 0.2
			totalPoints += int(price) + 1
		}
	}

	// Add 6 points if the day in the date of purchase is odd.
	day := rec.PurchaseDate[len(rec.PurchaseDate)-2:]
	dayInt, _ := strconv.Atoi(day)

	if dayInt%2 == 1 {
		totalPoints += 6
	}

	// Add 10 points if time of purchase is after 2:00 pm and before 4:00 pm.
	t, err := time.Parse("15:04", rec.PurchaseTime)

	if err != nil {
		return -1
	}

	twoPM, _ := time.Parse("15:04", "14:00")
	fourPM, _ := time.Parse("15:04", "16:00")

	if t.After(twoPM) && t.Before(fourPM) {
		totalPoints += 10
	}

	return totalPoints
}
