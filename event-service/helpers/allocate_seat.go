package helpers

import (
	"github.com/abidaziz9876/event-service/models"
	"github.com/abidaziz9876/event-service/repository"
)

func SeedSeats(db repository.PostgesDatabase) error {
	seats := []models.Seat{
		{Venue: "Cinema Hall 1", Row: "A", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "A", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "A", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "A", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "B", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "B", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "B", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "B", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 1", Row: "C", Number: 1, Category: "VIP"},
		{Venue: "Cinema Hall 1", Row: "C", Number: 2, Category: "VIP"},
		{Venue: "Cinema Hall 2", Row: "A", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "A", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "A", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "A", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "B", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "B", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "B", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "B", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 2", Row: "C", Number: 1, Category: "VIP"},
		{Venue: "Cinema Hall 2", Row: "C", Number: 2, Category: "VIP"},
		{Venue: "Cinema Hall 3", Row: "A", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "A", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "A", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "A", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "B", Number: 1, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "B", Number: 2, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "B", Number: 3, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "B", Number: 4, Category: "Regular"},
		{Venue: "Cinema Hall 3", Row: "C", Number: 1, Category: "VIP"},
		{Venue: "Cinema Hall 3", Row: "c", Number: 2, Category: "VIP"},
		// Add more seats as needed
	}

	for _, seat := range seats {
		if err := db.CreateSeats(&seat).Error; err != nil {
			return err
		}
	}
	return nil
}
