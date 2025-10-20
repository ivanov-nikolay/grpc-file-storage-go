package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/grpc-file-storage-go/internal/config"

	_ "github.com/lib/pq"
)

const (
	retries int           = 3
	delay   time.Duration = 5 * time.Second
)

func NewDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	if cfg.CreateDB {
		if err := CreateDB(cfg); err != nil {
			log.Printf("warning: failed to create database: %v", err)
		}
	}

	return connectWithRetry(cfg, retries, delay)
}

func CreateDB(cfg config.DatabaseConfig) error {
	superCfg := config.DatabaseConfig{
		Host:     cfg.Host,
		Port:     cfg.Port,
		User:     cfg.SuperUser,
		Password: cfg.SuperPassword,
		DBName:   "postgres",
		SSLMode:  cfg.SSLMode,
	}

	db, err := connect(superCfg)
	if err != nil {
		return fmt.Errorf("failed to connect to database as superUser: %v", err)
	}
	defer db.Close()

	var exists bool
	checkQuery := "SELECT EXISTS (SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(checkQuery, cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %v", err)
	}

	if !exists {
		createQuery := fmt.Sprintf("CREATE DATABASE %s", cfg.DBName)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %v", err)
		}
		log.Printf("database %s created successfully", cfg.DBName)
	} else {
		log.Printf("database %s already exists", cfg.DBName)
	}

	return nil
}

func connectWithRetry(cfg config.DatabaseConfig, maxRetries int, delay time.Duration) (*sql.DB, error) {
	var db *sql.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		db, err = connect(cfg)
		if err == nil {
			if err = db.Ping(); err == nil {
				log.Printf("database %s connected successfully", cfg.DBName)
				return db, nil
			}
			db.Close()
		}
		log.Printf("connection attempt %d retries failed: %v", i+1, err)

		if i < maxRetries-1 {
			log.Printf("retrying database connection in %v...", delay)
			time.Sleep(delay)
		}
	}

	return nil, fmt.Errorf("could not connect to database after %d retries: %v", maxRetries, err)
}

func connect(cfg config.DatabaseConfig) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}
