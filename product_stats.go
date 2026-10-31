package main

import (
	"encoding/json"
	"net/http"
)

type ProductStats struct {
	UserID     string `json:"user_id"`
	ProductID  string `json:"product_id"`
	ClickCount int    `json:"click_count"`
}

func countProductClicks(items []Event, userID, productID string) int {
	count := 0
	for _, event := range items {
		if event.UserID == userID && event.ProductID == productID && event.EventType == EventClick {
			count++
		}
	}
	return count
}

func getProductStatsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed!", http.StatusMethodNotAllowed)
		return
	}

	userID := r.URL.Query().Get("user_id")
	productID := r.URL.Query().Get("product_id")
	if userID == "" || productID == "" {
		http.Error(w, "user_id and product_id are required!", http.StatusBadRequest)
		return
	}

	result := ProductStats{
		UserID:     userID,
		ProductID:  productID,
		ClickCount: countProductClicks(events, userID, productID),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
