package repository

import (
	"github.com/abidaziz9876/user-service/models"
	"gorm.io/gorm"
)

type PostgesDatabase interface {
	RawQueryWithFind(result interface{}, sql string, values ...interface{}) *gorm.DB
	Create(value interface{}) *gorm.DB
	FindUserByEmailOrPhone(email string, phone string) (*models.UserModel, error)
	FindUserByID(id int64) (*models.UserModel, error)
	UpdateUser(user *models.UserModel) *gorm.DB
	DeleteUserByID(id int64) *gorm.DB
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

func (r *GormDBWrapper) FindUserByEmailOrPhone(email string, phone string) (*models.UserModel, error) {
	var user models.UserModel
	result := r.DB.Table("ticketbooking.users").Where("email = ? OR phone = ?", email, phone).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormDBWrapper) FindUserByID(id int64) (*models.UserModel, error) {
	var user models.UserModel
	result := r.DB.Table("ticketbooking.users").Where("id = ?", id).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func (r *GormDBWrapper) UpdateUser(user *models.UserModel) *gorm.DB {
	return r.DB.Table("ticketbooking.users").Save(user)
}

func (r *GormDBWrapper) DeleteUserByID(id int64) *gorm.DB {
	return r.DB.Table("ticketbooking.users").Delete(&models.UserModel{}, id)
}