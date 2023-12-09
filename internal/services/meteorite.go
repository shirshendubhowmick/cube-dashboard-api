package services

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

var nasaMeteoriteDataUrl = "https://data.nasa.gov/resource/gh4g-9sfh.json"

// name": "Aachen",
//     "id": "1",
//     "nametype": "Valid",
//     "recclass": "L5",
//     "mass": "21",
//     "fall": "Fell",
//     "year": "1880-01-01T00:00:00.000",
//     "reclat": "50.775000",
//     "reclong": "6.083330",
//     "geolocation": {
//       "latitude": "50.775",
//       "longitude": "6.08333"
//     }

type MeteoriteData struct {
	Name        string    `json:"name"`
	ID          string    `json:"id"`
	NameType    string    `json:"nametype"`
	RecClass    string    `json:"recclass"`
	Mass        string    `json:"mass"`
	Fall        string    `json:"fall"`
	Year        time.Time `json:"year"`
	RecLat      float32   `json:"reclat"`
	RecLong     float32   `json:"reclong"`
	GeoLocation struct {
		Latitude  float32 `json:"latitude"`
		Longitude float32 `json:"longitude"`
	} `json:"geolocation"`
}

func DownloadMeteoriteData() (*[]MeteoriteData, error) {
	response, err := resty.New().R().SetResult(&[]MeteoriteData{}).Get(nasaMeteoriteDataUrl)

	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	result, ok := response.Result().(*[]MeteoriteData)

	if !ok {
		fmt.Print()
		return nil, err
	}

	return result, err
}

func LoadMeteoriteData() error {
	meteoriteData, err := DownloadMeteoriteData()

	if err != nil {
		return err
	}

	fmt.Println(len(*meteoriteData))

	return nil
}
