package models

import (
	"time"
)

// Product model corresponds to the 'products' table in the database.
type Product struct {
	Id            uint      `gorm:"primaryKey" json:"id"`
	SKU           string    `gorm:"uniqueIndex;size:100" json:"sku"`
	Name          string    `gorm:"size:255;not null" json:"name"`
	Description   string    `json:"description"`
	Price         float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	StockQuantity int       `gorm:"not null" json:"stock_quantity"`
	CategoryID    uint      `json:"category_id"`
	ImageURL      string    `gorm:"size:255" json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
