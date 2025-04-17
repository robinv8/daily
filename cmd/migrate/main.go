package main

import (
	"daily/internal/db"
	"daily/internal/entity"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Starting database schema migration...")
	
	// Get database connection
	engine := db.NewDB()
	defer engine.Close()
	
	// Sync the updated schema
	err := engine.Sync2(new(entity.SiteInfo))
	if err != nil {
		log.Fatalf("Failed to sync database schema: %v", err)
	}
	
	fmt.Println("Database schema migration completed successfully!")
	fmt.Println("The following columns in site_info table have been updated:")
	fmt.Println("- title: VARCHAR(100) -> VARCHAR(255)")
	fmt.Println("- keywords: VARCHAR(100) -> VARCHAR(255)")
	fmt.Println("- origin_url: VARCHAR(100) -> VARCHAR(512)")
	fmt.Println("- image_url: VARCHAR(100) -> VARCHAR(512)")
}
