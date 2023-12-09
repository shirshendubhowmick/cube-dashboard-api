package db

import (
	"fmt"
	"main/configs"
	"time"

	uuid "github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MeteoriteData struct {
	ID            uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Name          string
	Mass          float32
	DiscoveryType string `gorm:"discovery_type_check,discovery_type in ('fell', 'found')"`
	Year          *uint
	Latitude      float32
	Longitude     float32
	Region        string
	Country       string
	City          string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Kolkata", configs.PostgresHost, configs.PostgresUser, configs.PostgresPassword, configs.PostgresDBName, configs.PostgresPort)
	DB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	if err := DB.AutoMigrate(&MeteoriteData{}); err != nil {
		return nil, err
	}

	return DB, nil
}
