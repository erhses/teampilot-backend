package dts

import (
	"teampilot/constants"
	"time"

	"gorm.io/gorm"
)

type (
	Base struct {
		ID        uint           `gorm:"primary_key;autoIncrement:true" json:"id"`
		CreatedAt time.Time      `gorm:"autoCreateTime;index" json:"created_at"`
		UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
		DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	}
	User struct {
		Base
		Name              string             `gorm:"column:name" json:"name"`
		UserType          constants.UserType `gorm:"column:user_type" json:"user_type"`
		Email             string             `gorm:"column:email" json:"email"`
		Password          string             `gorm:"column:password" json:"password"`
		IsActive          bool               `gorm:"column:is_active" json:"is_active"`
		RoleID            uint               `gorm:"column:role_id" json:"role_id"`
		Role              *Role              `gorm:"foreignKey:RoleID" json:"role,omitempty"`
		NotificationToken string             `gorm:"column:notification_token" json:"notification_token"`
		DeviceToken       string             `gorm:"column:device_token"  json:"device_token"`
		StudentProfile    *StudentProfile
	}
)
