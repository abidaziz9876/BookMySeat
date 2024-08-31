package models

import (
	"time"

	"gorm.io/gorm"
)
type BookingRequest struct {
    EventID      uint   `json:"event_id"`
    OccurrenceID uint   `json:"occurrence_id"`
    SeatIDs      []uint `json:"seat_ids"` // List of seat IDs being booked
    UserID       uint   `json:"user_id"`
}
// Event represents a movie or show event in the booking system.
type Event struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Title       string         `gorm:"size:255;not null" json:"title"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	Language    string         `gorm:"size:100;not null" json:"language"`
	Genre       string         `gorm:"size:100;not null" json:"genre"`
	Rating      float32        `gorm:"type:float" json:"rating,omitempty"`
	CreatedAt   time.Time      `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime" json:"updatedAt"`
	IsActive    bool           `gorm:"not null;default:true" json:"isActive"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName sets the table name for the Event struct.
func (Event) TableName() string {
	return "ticketbooking.events"
}

// EventOccurrence represents the occurrence of an event.
type EventOccurrences struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EventID   uint      `gorm:"not null;index" json:"eventId"`
	Event     Event     `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"event"`
	Venue     string    `gorm:"size:255;not null" json:"venue"`
	Date      time.Time `gorm:"not null" json:"date"`
	TimeSlot  string    `gorm:"size:50;not null" json:"timeSlot"`
	Duration  int       `gorm:"not null" json:"duration"`
	Price     float64   `gorm:"not null" json:"price"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

// TableName sets the table name for the EventOccurrence struct.
func (EventOccurrences) TableName() string {
	return "ticketbooking.occurrence"
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

type Booking struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null" json:"userId"`       // Reference to the user making the booking
	EventID      uint      `gorm:"not null" json:"eventId"`      // Reference to the event being booked
	OccurrenceID uint      `gorm:"not null" json:"occurrenceId"` // Reference to the specific occurrence of the event
	TotalAmount  float64   `gorm:"not null" json:"totalAmount"`  // Total amount for the booking
	Status       string    `gorm:"size:20;not null" json:"status"` // Status of the booking (e.g., Pending, Confirmed, Canceled)
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

type SeatAllocation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	OccurrenceID uint      `gorm:"not null" json:"occurrenceId"`
	SeatID       uint      `gorm:"not null" json:"seatId"`
	Status       string    `gorm:"size:20;not null" json:"status"` // e.g., Available, Reserved, Booked
	BookedAt     time.Time `json:"bookedAt,omitempty"`             // Timestamp when the seat was booked
	ReservedAt   time.Time `json:"reservedAt,omitempty"`           // Timestamp when the seat was reserved
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UserID       uint      `gorm:"not null" json:"userId"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
}

func (SeatAllocation) TableName() string {
	return "seat_allocation"
}