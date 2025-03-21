package models

import (
	"time"

)

// User model
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Email        string    `gorm:"unique;not null;type:varchar(255)"`
	PasswordHash string    `gorm:"not null;type:varchar(255)"`
	FirstName    string    `gorm:"type:varchar(50)"`
	LastName     string    `gorm:"type:varchar(50)"`
	CreatedAt    time.Time `gorm:"default:now()"`
	UpdatedAt    time.Time `gorm:"default:now()"`

	CartItems []Cart   `gorm:"foreignKey:UserID"`
	Orders    []Order  `gorm:"foreignKey:UserID"`
	Reviews   []Review `gorm:"foreignKey:UserID"`
}

// Category model
type Category struct {
	ID          uint   `gorm:"primaryKey"`
	Name        string `gorm:"unique;not null;type:varchar(50)"`
	Description string `gorm:"type:text"`

	Products []Product `gorm:"foreignKey:CategoryID"`
}

// Product model
type Product struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"not null;type:varchar(100)"`
	Price       float64   `gorm:"not null;type:decimal(10,2)"`
	Description string    `gorm:"not null;type:text"`
	CategoryID  uint      `gorm:"index"`
	Stock       int       `gorm:"not null;default:0"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`

	Category   Category       `gorm:"foreignKey:CategoryID"`
	Images     []ProductImage `gorm:"foreignKey:ProductID"`
	Specs      []ProductSpec  `gorm:"foreignKey:ProductID"`
	CartItems  []Cart         `gorm:"foreignKey:ProductID"`
	OrderItems []OrderItem    `gorm:"foreignKey:ProductID"`
	Reviews    []Review       `gorm:"foreignKey:ProductID"`
}

// ProductImage model
type ProductImage struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index"`
	URL       string `gorm:"not null;type:varchar(255)"`
	IsPrimary bool   `gorm:"default:false"`
	Order     int    `gorm:"default:0"`

	Product Product `gorm:"foreignKey:ProductID"`
}

// ProductSpec model
type ProductSpec struct {
	ID        uint   `gorm:"primaryKey"`
	ProductID uint   `gorm:"index"`
	Key       string `gorm:"not null;type:varchar(50)"`
	Value     string `gorm:"not null;type:varchar(255)"`

	Product Product `gorm:"foreignKey:ProductID"`
}

// Cart model
type Cart struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	ProductID uint      `gorm:"index"`
	Quantity  int       `gorm:"not null;check:quantity > 0"`
	AddedAt   time.Time `gorm:"default:now()"`

	User    User    `gorm:"foreignKey:UserID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

// Order model
type Order struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	Total     float64   `gorm:"not null;type:decimal(10,2)"`
	Status    string    `gorm:"not null;default:'pending';type:varchar(20)"`
	CreatedAt time.Time `gorm:"default:now()"`
	UpdatedAt time.Time `gorm:"default:now()"`

	User  User        `gorm:"foreignKey:UserID"`
	Items []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem model
type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"index"`
	ProductID uint    `gorm:"index"`
	Quantity  int     `gorm:"not null;check:quantity > 0"`
	Price     float64 `gorm:"not null;type:decimal(10,2)"`

	Order   Order   `gorm:"foreignKey:OrderID"`
	Product Product `gorm:"foreignKey:ProductID"`
}

// Review model
type Review struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"index"`
	ProductID uint      `gorm:"index"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5"`
	Comment   string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"default:now()"`

	User    User    `gorm:"foreignKey:UserID"`
	Product Product `gorm:"foreignKey:ProductID"`
}
