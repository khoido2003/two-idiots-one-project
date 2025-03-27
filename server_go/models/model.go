// server/models/models.go
package models

import (
	"time"
)

// User model
type User struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null;type:varchar(255)" json:"email"`
	PasswordHash string    `gorm:"not null;type:varchar(255)" json:"passwordHash"`
	FirstName    string    `gorm:"type:varchar(50)" json:"firstName"`
	LastName     string    `gorm:"type:varchar(50)" json:"lastName"`
	CreatedAt    time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"default:now()" json:"updatedAt"`

	CartItems []Cart   `gorm:"foreignKey:UserID" json:"cartItems"`
	Orders    []Order  `gorm:"foreignKey:UserID" json:"orders"`
	Reviews   []Review `gorm:"foreignKey:UserID" json:"reviews"`
}

// Category model
type Category struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"unique;not null;type:varchar(50)" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updatedAt"`

	Products []Product `gorm:"foreignKey:CategoryID" json:"products"`
}

// Product model
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;type:varchar(100)" json:"name"`
	Price       float64   `gorm:"not null;type:decimal(10,2)" json:"price"`
	Description string    `gorm:"not null;type:text" json:"description"`
	CategoryID  uint      `gorm:"index" json:"categoryId"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updatedAt"`

	Category   Category       `gorm:"foreignKey:CategoryID" json:"category"`
	Images     []ProductImage `gorm:"foreignKey:ProductID" json:"images"`
	Specs      []ProductSpec  `gorm:"foreignKey:ProductID" json:"specs"`
	CartItems  []Cart         `gorm:"foreignKey:ProductID" json:"cartItems"`
	OrderItems []OrderItem    `gorm:"foreignKey:ProductID" json:"orderItems"`
	Reviews    []Review       `gorm:"foreignKey:ProductID" json:"reviews"`
}

// ProductImage model
type ProductImage struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `gorm:"index" json:"productId"`
	URL       string `gorm:"not null;type:varchar(255)" json:"url"`
	IsPrimary bool   `gorm:"default:false" json:"isPrimary"`
	Order     int    `gorm:"default:0" json:"order"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// ProductSpec model
type ProductSpec struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	ProductID uint   `gorm:"index" json:"productId"`
	Key       string `gorm:"not null;type:varchar(50)" json:"key"`
	Value     string `gorm:"not null;type:varchar(255)" json:"value"`

	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Cart model
type Cart struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	ProductID uint      `gorm:"index" json:"productId"`
	Quantity  int       `gorm:"not null;check:quantity > 0" json:"quantity"`
	AddedAt   time.Time `gorm:"default:now()" json:"addedAt"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Order model
type Order struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	Total     float64   `gorm:"not null;type:decimal(10,2)" json:"total"`
	Status    string    `gorm:"not null;default:'pending';type:varchar(20)" json:"status"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`
	UpdatedAt time.Time `gorm:"default:now()" json:"updatedAt"`

	User  User        `gorm:"foreignKey:UserID" json:"user"`
	Items []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
}

// OrderItem model
type OrderItem struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	OrderID   uint    `gorm:"index" json:"orderId"`
	ProductID uint    `gorm:"index" json:"productId"`
	Quantity  int     `gorm:"not null;check:quantity > 0" json:"quantity"`
	Price     float64 `gorm:"not null;type:decimal(10,2)" json:"price"`

	Order   Order   `gorm:"foreignKey:OrderID" json:"order"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}

// Review model
type Review struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"userId"`
	ProductID uint      `gorm:"index" json:"productId"`
	Rating    int       `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"default:now()" json:"createdAt"`

	User    User    `gorm:"foreignKey:UserID" json:"user"`
	Product Product `gorm:"foreignKey:ProductID" json:"product"`
}
