package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	godotenv.Load("../../.env")

	// Get database configuration from environment
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "customermx")

	// Build connection string
	connString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=disable",
		host, port, user, password,
	)

	fmt.Println("CustomerMX Database Migration")
	fmt.Println("================================")
	fmt.Printf("Host: %s:%s\n", host, port)
	fmt.Printf("Database: %s\n", dbname)
	fmt.Printf("User: %s\n\n", user)

	ctx := context.Background()

	// Connect to postgres database first
	fmt.Println("📡 Connecting to PostgreSQL...")
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database: %v\n", err)
	}
	defer conn.Close(ctx)
	fmt.Println("✅ Connected to PostgreSQL\n")

	// Check if database exists, create if not
	fmt.Printf("🔍 Checking if database '%s' exists...\n", dbname)
	var exists bool
	err = conn.QueryRow(ctx,
		"SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)",
		dbname,
	).Scan(&exists)
	if err != nil {
		log.Fatalf("❌ Error checking database: %v\n", err)
	}

	if !exists {
		fmt.Printf("📦 Creating database '%s'...\n", dbname)
		_, err = conn.Exec(ctx, fmt.Sprintf("CREATE DATABASE %s", dbname))
		if err != nil {
			log.Fatalf("❌ Error creating database: %v\n", err)
		}
		fmt.Println("✅ Database created\n")
	} else {
		fmt.Println("✅ Database exists\n")
	}

	// Close connection to postgres and connect to target database
	conn.Close(ctx)

	connString = fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	conn, err = pgx.Connect(ctx, connString)
	if err != nil {
		log.Fatalf("❌ Unable to connect to database %s: %v\n", dbname, err)
	}
	defer conn.Close(ctx)

	// Get migrations directory
	migrationsDir := "../../migrations"
	if len(os.Args) > 1 {
		migrationsDir = os.Args[1]
	}

	// Read migration files
	fmt.Println("📂 Reading migration files...")
	files, err := filepath.Glob(filepath.Join(migrationsDir, "V*.sql"))
	if err != nil {
		log.Fatalf("❌ Error reading migrations: %v\n", err)
	}

	if len(files) == 0 {
		log.Fatal("❌ No migration files found")
	}

	sort.Strings(files)
	fmt.Printf("Found %d migration(s)\n\n", len(files))

	// Create migrations tracking table if not exists
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			filename TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		log.Fatalf("❌ Error creating schema_migrations table: %v\n", err)
	}

	// Load already-applied migrations
	rows, err := conn.Query(ctx, "SELECT filename FROM schema_migrations")
	if err != nil {
		log.Fatalf("❌ Error reading schema_migrations: %v\n", err)
	}
	applied := map[string]bool{}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		applied[name] = true
	}
	rows.Close()

	// Bootstrap: if schema_migrations is empty but tables already exist,
	// seed all found migration files as already applied.
	if len(applied) == 0 {
		var tableExists bool
		conn.QueryRow(ctx,
			"SELECT EXISTS(SELECT 1 FROM information_schema.tables WHERE table_name='brands')",
		).Scan(&tableExists)
		if tableExists {
			fmt.Println("⚙️  Bootstrapping schema_migrations from existing database...\n")
			for _, file := range files {
				filename := filepath.Base(file)
				_, err = conn.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1) ON CONFLICT DO NOTHING", filename)
				if err != nil {
					log.Fatalf("❌ Error seeding migration %s: %v\n", filename, err)
				}
				applied[filename] = true
				fmt.Printf("   ✅ Marked as applied: %s\n", filename)
			}
			fmt.Println()
		}
	}

	// Execute migrations
	fmt.Println("🚀 Running migrations...\n")
	for _, file := range files {
		filename := filepath.Base(file)

		if applied[filename] {
			fmt.Printf("   ⏭️  %s (already applied)\n\n", filename)
			continue
		}

		fmt.Printf("   Running: %s\n", filename)

		content, err := os.ReadFile(file)
		if err != nil {
			log.Fatalf("❌ Error reading %s: %v\n", filename, err)
		}

		_, err = conn.Exec(ctx, string(content))
		if err != nil {
			log.Fatalf("❌ Error executing %s: %v\n", filename, err)
		}

		_, err = conn.Exec(ctx, "INSERT INTO schema_migrations (filename) VALUES ($1)", filename)
		if err != nil {
			log.Fatalf("❌ Error recording migration %s: %v\n", filename, err)
		}

		fmt.Printf("   ✅ %s completed\n\n", filename)
	}

	fmt.Println("================================")
	fmt.Println("✅ All migrations completed successfully!\n")

	// Show summary
	fmt.Println("📊 Database Summary:")
	summaryRows, err := conn.Query(ctx, `
		SELECT 'Brands' as table_name, COUNT(*) as count FROM brands
		UNION ALL
		SELECT 'Vehicles', COUNT(*) FROM vehicles
		UNION ALL
		SELECT 'Users', COUNT(*) FROM users
		UNION ALL
		SELECT 'Events', COUNT(*) FROM events
	`)
	if err != nil {
		log.Printf("Warning: Could not fetch summary: %v\n", err)
		return
	}
	defer summaryRows.Close()

	fmt.Println("\nTable          | Count")
	fmt.Println("---------------|------")
	for summaryRows.Next() {
		var tableName string
		var count int
		if err := summaryRows.Scan(&tableName, &count); err != nil {
			continue
		}
		fmt.Printf("%-14s | %d\n", tableName, count)
	}

	fmt.Println("\n🚀 Ready to go!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
