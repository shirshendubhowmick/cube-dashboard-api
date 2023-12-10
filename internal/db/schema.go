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
	Name          string    `gorm:"unique"`
	Mass          float32
	DiscoveryType string `gorm:"check:discovery_type_check,discovery_type in ('fell', 'found')"`
	Year          *uint
	Latitude      float32
	Longitude     float32
	Country       string
	Region        string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type Users struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()" json:"id"`
	UserName  string    `gorm:"unique" json:"username"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Active    bool      `json:"active"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var DB *gorm.DB

func Init() (*gorm.DB, error) {
	if DB != nil {
		return DB, nil
	}
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require TimeZone=Asia/Kolkata", configs.PostgresHost, configs.PostgresUser, configs.PostgresPassword, configs.PostgresDBName, configs.PostgresPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	DB = db

	if err := DB.AutoMigrate(&MeteoriteData{}, &Users{}); err != nil {
		return nil, err
	}

	return DB, nil
}
