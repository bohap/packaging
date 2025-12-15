package itest

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	dbContainer, err := initializePgSQLContainer(ctx)
	if err != nil {
		log.Printf("Failed to start PostgreSQL container")
		os.Exit(1)
	}

	result := m.Run()

	if err := dbContainer.Terminate(ctx); err != nil {
		log.Printf("Failed to terminate PostgreSQL container")
		os.Exit(1)
	}

	os.Exit(result)
}

func initializePgSQLContainer(ctx context.Context) (testcontainers.Container, error) {
	dbName := "gotest"
	dbUser := "gotest"
	dbPassword := "gotest"

	postgresContainer, err := postgres.Run(
		ctx,
		"public.ecr.aws/docker/library/postgres:16.2-alpine",
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)

	if err != nil {
		return nil, err
	}

	// Get the port mapped to 5432 and set as ENV
	port, _ := postgresContainer.MappedPort(ctx, "5432")

	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", port.Port())
	os.Setenv("DB_USERNAME", dbUser)
	os.Setenv("DB_PASSWORD", dbPassword)
	os.Setenv("DB_NAME", dbName)

	// initialize the database
	connStr, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgresDriver.Open(connStr), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := initializeSchema(db); err != nil {
		return nil, err
	}

	return postgresContainer, nil
}

func initializeSchema(db *gorm.DB) error {
	initSQL := `CREATE TABLE packs (size BIGINT PRIMARY KEY);`
	if err := db.Exec(initSQL).Error; err != nil {
		return err
	}

	return nil
}

func cleanupDb(db *gorm.DB) error {
	sql := `DELETE FROM packs WHERE 1=1;`
	if err := db.Exec(sql).Error; err != nil {
		return err
	}

	return nil
}
