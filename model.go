package main

import (
	"time"

	"gorm.io/gorm"
)

type Deal struct {
	ID        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Link      string
	Price     string
	Expired   bool
}
