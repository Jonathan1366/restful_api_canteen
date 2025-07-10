package db_drivers

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func InitDbPool(app *fiber.App) *pgxpool.Pool {
	// Middleware untuk CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: os.Getenv("API_ALLOW_ORIGINS"), // Mengambil dari environment variable
		AllowHeaders: os.Getenv("API_ALLOW_HEADERS"), // Mengambil dari environment variable
		AllowMethods: os.Getenv("API_ALLOW_METHODS"), // Mengambil dari environment variable
	}))	

	// Ambil DATABASE_URL dari environment variable
	urldb := os.Getenv("DATABASE_URL")
	if urldb == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}
	
	// Parse database URL
	config, err := pgxpool.ParseConfig(urldb)
	if err != nil {	
		log.Fatalf("Unable to parse database configuration: %v\n", err)
	}

	// Mengatur runtime parameter PostgreSQL
	config.ConnConfig.RuntimeParams["prefer_simple_protocol"] = "true" 

	// Mengambil nilai MaxConns dari environment variable
	maxConns, err := strconv.Atoi(os.Getenv("MAX_CONNS"))
	if err != nil {
		maxConns = 100 // Default nilai jika MAX_CONNS tidak valid
	}
	config.MaxConns = int32(maxConns)

	// Mengambil nilai MinConns dari environment variable
	minConns, err := strconv.Atoi(os.Getenv("MIN_CONNS"))
	if err != nil {
		minConns = 2 // Default nilai jika MIN_CONNS tidak valid
	}
	config.MinConns = int32(minConns)

	// Mengambil nilai MaxConnIdleTime dari environment variable
	maxConnIdleTime := os.Getenv("MAX_CONN_IDLE_TIME")
	if maxConnIdleTime == ""{
		maxConnIdleTime = "10m" // Default nilai jika MAX_CONN_IDLE_TIME tidak ada
	}
	duration, err := time.ParseDuration(maxConnIdleTime)
	if err != nil {
		duration = 10 * time.Minute // Default jika parsing gagal
	}
	config.MaxConnIdleTime = duration

	//Jalankan hook ini setiap kali koneksi baru dibuat
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		_, err := conn.Exec(ctx, "SET search_path = auth_service, public")
		return err
	}
	// Membuat pool koneksi
	dbPool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v\n", err)
	}

	// Register fungsi untuk menutup pool saat aplikasi shutdown
	app.Hooks().OnShutdown(func() error {
		CloseDbPool()
		return nil
	})
	return dbPool
}

// Menutup koneksi pool saat aplikasi shutdown

func CloseDbPool() {
	if dbPool != nil {
		dbPool.Close()
	}
}

func GetDbPool() *pgxpool.Pool {
	return dbPool
}
