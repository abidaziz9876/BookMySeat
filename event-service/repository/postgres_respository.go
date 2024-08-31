package repository

import (
	"github.com/abidaziz9876/event-service/models"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PostgesDatabase interface {
	RawQueryWithFind(result interface{}, sql string, values ...interface{}) *gorm.DB
	CreateEvent(value *models.Event) *gorm.DB
	CreateEventOccurrences(value interface{}) *gorm.DB
	CreateSeatAllocation(value interface{}) *gorm.DB
	CreateSeats(value interface{}) *gorm.DB
	CreateBooking(value interface{}) *gorm.DB
	Delete(value interface{}) *gorm.DB
	UpdateQuery(sql string, values ...interface{}) *gorm.DB
}

type GormDBWrapper struct {
	DB *gorm.DB
}

func (r *GormDBWrapper) RawQueryWithFind(result interface{}, sql string, values ...interface{}) *gorm.DB {
	return r.DB.Raw(sql, values...).Find(result)
}
func (r *GormDBWrapper) CreateEvent(value *models.Event) *gorm.DB {
	result := r.DB.Table("ticketbooking.events").Create(value)
	if result.Error != nil {
		log.Error("Error occurred while creating events " + result.Error.Error())
	}
	return result
}

func (r *GormDBWrapper) CreateEventOccurrences(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.occurrence").Create(value)
}

func (r *GormDBWrapper) Delete(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.events").Delete(value)
}

func (r *GormDBWrapper) CreateSeatAllocation(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.seat_allocation").Create(value)
}

func (r *GormDBWrapper) CreateSeats(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.seats").Create(value)
}

func (r *GormDBWrapper) CreateBooking(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.booking").Create(value)
}

func (r *GormDBWrapper) UpdateQuery(sql string, values ...interface{}) *gorm.DB {
	return r.DB.Exec(sql, values...)
}