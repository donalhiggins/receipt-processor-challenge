package main

import (
	"testing"
)

func TestCalculatePoints(t *testing.T) {
	tests := []struct {
		name   string
		rec    Receipt
		points int
	}{
		{
			name: "Test Receipt One",
			rec: Receipt{
				Retailer:     "Target",
				PurchaseDate: "2022-01-01",
				PurchaseTime: "13:01",
				Items: []Item{
					{
						ShortDescription: "Mountain Dew 12PK",
						Price:            6.49,
					},
					{
						ShortDescription: "Emils Cheese Pizza",
						Price:            12.25,
					},
					{
						ShortDescription: "Knorr Creamy Chicken",
						Price:            1.26,
					},
					{
						ShortDescription: "Doritos Nacho Cheese",
						Price:            3.35,
					},
					{
						ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
						Price:            12.00,
					},
				},
				Total: 35.35,
			},
			points: 28,
		},
		{
			name: "Test Receipt Two",
			rec: Receipt{
				Retailer:     "M&M Corner Market",
				PurchaseDate: "2022-03-20",
				PurchaseTime: "14:33",
				Items: []Item{
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
					{
						ShortDescription: "Gatorade",
						Price:            2.25,
					},
				},
				Total: 9.00,
			},
			points: 109,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := calculatePoints(tt.rec)
			if actual != tt.points {
				t.Errorf("Expected: %d, Actual: %d", tt.points, actual)
			}
		})
	}
}
