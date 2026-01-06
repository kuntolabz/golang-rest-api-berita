package models

import (
	"time"

	"github.com/google/uuid"
)

type Ms_user struct {
	IdUser    uuid.UUID `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primary_key"`
	Name      string    `json:"name" binding:"required" gorm:"type:varchar(50);not null"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Username  string    `json:"username"`
	Alamat    string    `json:"alamatUser"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy string    `json:"created_by"`
	UpdatedAt time.Time `json:"updated_at"`
	UpdatedBy string    `json:"updated_by"`
	IdRole    string    `json:"id_role"`
}
