package models

import "time"

type BookingRequest struct {
	EventID      uint   `json:"event_id"`
	OccurrenceID uint   `json:"occurrence_id"`
	SeatIDs      []uint `json:"seat_ids"` // List of seat IDs being booked
	UserID       uint   `json:"user_id"`
}

type Booking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"userId"`         // Reference to the user making the booking
	EventID      uint      `gorm:"not null" json:"eventId"`        // Reference to the event being booked
	OccurrenceID uint      `gorm:"not null" json:"occurrenceId"`   // Reference to the specific occurrence of the event
	TotalAmount  float64   `gorm:"not null" json:"totalAmount"`    // Total amount for the booking
	Status       string    `gorm:"size:20;not null" json:"status"` // Status of the booking (e.g., Pending, Confirmed, Canceled)
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type Seat struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Venue     string    `gorm:"size:255;not null" json:"venue"`
	Row       string    `gorm:"size:10;not null" json:"row"`
	Number    int       `gorm:"not null" json:"number"`
	Category  string    `gorm:"size:50;not null" json:"category"` // e.g., VIP, Regular
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}