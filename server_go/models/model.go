// server/models/models.go
package models

import (
	"time"
)

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
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
	Role         UserRole  `gorm:"not null;type:varchar(20);default:'user'" json:"role"`
	Phone        string    `gorm:"unique;not null;type:varchar(255)" json:"phone"`

	CartItems []Cart   `gorm:"foreignKey:UserID" json:"cartItems"`
	Orders    []Order  `gorm:"foreignKey:UserID" json:"orders"`
	Reviews   []Review `gorm:"foreignKey:UserID" json:"reviews"`

	AddressLine1 string `gorm:"type:varchar(255)" json:"addressLine1"`
	AddressLine2 string `gorm:"type:varchar(255)" json:"addressLine2,omitempty"`
	City         string `gorm:"type:varchar(100)" json:"city"`
	State        string `gorm:"type:varchar(100)" json:"state"`
	PostalCode   string `gorm:"type:varchar(20)" json:"postalCode"`
	Country      string `gorm:"type:varchar(100)" json:"country"`
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
	ID           uint        `gorm:"primaryKey" json:"id"`
	UserID       uint        `gorm:"not null" json:"userId"`
	Total        int64       `gorm:"not null" json:"total"` // In cents
	CreatedAt    time.Time   `json:"createdAt"`
	AddressLine1 string      `gorm:"type:varchar(255)" json:"line1"`
	AddressLine2 string      `gorm:"type:varchar(255)" json:"line2,omitempty"`
	City         string      `gorm:"type:varchar(100)" json:"city"`
	State        string      `gorm:"type:varchar(100)" json:"state"`
	PostalCode   string      `gorm:"type:varchar(20)" json:"postalCode"`
	Country      string      `gorm:"type:varchar(100)" json:"country"`
	Items        []OrderItem `gorm:"foreignKey:OrderID" json:"items"`
	User         User        `gorm:"foreignKey:UserID" json:"user"`
    Status       string      `gorm:"type:varchar(50);default:'pending'" json:"status"`
}

type OrderItem struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	OrderID   uint   `gorm:"not null" json:"orderId"`
	ProductID uint   `gorm:"not null" json:"productId"`
	Quantity  int    `gorm:"not null" json:"quantity"`
	Price     int64  `gorm:"not null" json:"price"`         // Price at purchase time, in cents
	Name      string `gorm:"type:varchar(255)" json:"name"` // Product name at purchase time
}

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

type ShippingAddress struct {
    ID         uint   `gorm:"primaryKey" json:"id"`
    OrderID    uint   `gorm:"column:order_id;not null" json:"orderId"`
    Line1      string `gorm:"type:varchar(255);not null" json:"line1"`
    Line2      string `gorm:"type:varchar(255)" json:"line2"`
    City       string `gorm:"type:varchar(100);not null" json:"city"`
    State      string `gorm:"type:varchar(100);not null" json:"state"`
    PostalCode string `gorm:"type:varchar(20);not null" json:"postalCode"`
    Country    string `gorm:"type:varchar(100);not null" json:"country"`
}
