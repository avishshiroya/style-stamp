package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct{
	ID uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"idf"`
	Username string `gorm:"type:string" json:"username"`
	Password string `gorm:"type:string" json:"password"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func(u *User) BeforeCreate(tx *gorm.DB)(err error) {
	if u.ID == uuid.Nil{
		u.ID =  uuid.New()
	}
	return nil
}