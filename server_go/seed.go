package main

import (
	"log"
	"server/models"
	"time"

	"gorm.io/gorm"
)

func SeedCategories(db *gorm.DB) error {
	categories := []models.Category{
		{
			Name:        "Electronics",
			Description: "Devices and gadgets powered by electricity",
		},
		{
			Name:        "Clothing",
			Description: "Apparel and fashion items",
		},
		{
			Name:        "Books",
			Description: "Printed and digital reading materials",
		},
		{
			Name:        "Home & Kitchen",
			Description: "Appliances and utensils for home use",
		},
		{
			Name:        "Toys & Games",
			Description: "Entertainment items for kids and adults",
		},
	}

	// Check before save to avoid duplicate
	for _, category := range categories {
		var existingCategory models.Category

		if err := db.Where("name=?", category.Name).First(&existingCategory).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				category.CreatedAt = time.Now()
				category.UpdatedAt = time.Now()

				if err := db.Create(&category).Error; err != nil {
					log.Printf("Failed to seed category %s: %v", category.Name, err)
					return err
				}

				log.Printf("Seeded category: %s", category.Name)
			} else {
				log.Printf("Error checking category %s: %v", category.Name, err)
				return err
			}
		} else {
			log.Printf("Category %s already exists, skipping", category.Name)
		}

	}

	return nil
}
