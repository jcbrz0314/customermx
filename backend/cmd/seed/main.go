package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/customermx/backend/internal/config"
	"github.com/customermx/backend/internal/domain/user"
	"github.com/customermx/backend/internal/infra/db"
	"github.com/customermx/backend/internal/infra/security"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Try loading from backend directory
		godotenv.Load("../../.env")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	dbConn, err := db.New(cfg.Database.GetDSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	fmt.Println("🌱 CustomerMX Seed Data")
	fmt.Println("========================\n")

	// Initialize services
	jwtService := security.NewJWTService(
		cfg.JWT.AccessSecret,
		cfg.JWT.RefreshSecret,
		cfg.JWT.AccessTTL,
		cfg.JWT.RefreshTTL,
	)
	passwordService := security.NewPasswordService()
	userRepo := user.NewRepository(dbConn.Pool)
	userService := user.NewService(userRepo, jwtService, passwordService)

	ctx := context.Background()

	// Check if admin already exists
	adminEmail := getEnv("ADMIN_EMAIL", "admin@customermx.com")
	_, err = userService.GetUserByEmail(ctx, adminEmail)
	if err == nil {
		fmt.Printf("⚠️  Admin user already exists: %s\n", adminEmail)
		fmt.Println("✅ Seed data complete")
		os.Exit(0)
	}

	// Create admin user
	adminPassword := getEnv("ADMIN_PASSWORD", "admin123")
	adminName := getEnv("ADMIN_NAME", "Administrator")

	fmt.Println("📝 Creating admin user...")
	fmt.Printf("   Email: %s\n", adminEmail)
	fmt.Printf("   Name: %s\n", adminName)
	fmt.Printf("   Password: %s\n\n", maskPassword(adminPassword))

	createReq := &user.CreateUserRequest{
		Name:     adminName,
		Email:    adminEmail,
		Password: adminPassword,
		Role:     user.RoleAdmin,
		BrandID:  nil,
	}

	adminUser, err := userService.CreateUser(ctx, createReq)
	if err != nil {
		log.Fatalf("❌ Failed to create admin user: %v", err)
	}

	fmt.Printf("✅ Admin user created successfully!\n")
	fmt.Printf("   ID: %s\n", adminUser.ID)
	fmt.Printf("   Email: %s\n", adminUser.Email)
	fmt.Printf("   Role: %s\n\n", adminUser.Role)

	fmt.Println("🎉 Seed data complete!")
	fmt.Println("\nYou can now login with:")
	fmt.Printf("  Email: %s\n", adminEmail)
	fmt.Printf("  Password: %s\n", adminPassword)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskPassword(password string) string {
	if len(password) <= 4 {
		return "****"
	}
	return password[:2] + "****" + password[len(password)-2:]
}
