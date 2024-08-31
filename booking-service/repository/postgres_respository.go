package repository

import (
	"gorm.io/gorm"
)

type PostgesDatabase interface {
	RawQueryWithFind(result interface{}, sql string, values ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
}

type GormDBWrapper struct {
	DB *gorm.DB
}

func (r *GormDBWrapper) RawQueryWithFind(result interface{}, sql string, values ...interface{}) *gorm.DB {
	return r.DB.Raw(sql, values...).Find(result)
}
func (r *GormDBWrapper) Create(value interface{}) *gorm.DB {
	return r.DB.Table("ticketbooking.users").Create(value)
}
