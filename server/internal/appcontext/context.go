package appcontext

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
	"server/internal/repository"
	"server/internal/service"
)

type AppContext struct {
	DB             *gorm.DB
	PacksService   service.PacksService
	PackingService service.PackagingService
}

func BuildAppContext() *AppContext {
	db := createDbConnection()
	repo := repository.NewPacksRepository(db)
	packsService := service.NewPacksService(repo)
	packingService := service.NewPackagingService(packsService)
	return &AppContext{
		DB:             db,
		PacksService:   packsService,
		PackingService: packingService,
	}
}

func createDbConnection() *gorm.DB {
	dbHost := readOsEnv("DB_HOST")
	dbPort := readOsEnv("DB_PORT")
	dbUsername := readOsEnv("DB_USERNAME")
	dbPassword := readOsEnv("DB_PASSWORD")
	dbName := readOsEnv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=%v port=%v sslmode=disable",
		dbHost, dbUsername, dbPassword, dbName, dbPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error opening DB connection: %v", err)
	}

	return db
}

func readOsEnv(key string) string {
	res := os.Getenv(key)
	if res == "" {
		log.Fatalf("Required env var %s is NOT present", key)
	}

	return res
}
