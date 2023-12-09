package services

import (
	"errors"
	"fmt"
	"main/configs"
	"main/internal/db"
	"main/internal/services/logger"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"gorm.io/gorm/clause"
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
	Name        string  `json:"name"`
	ID          string  `json:"id"`
	NameType    string  `json:"nametype"`
	RecClass    string  `json:"recclass"`
	Mass        float32 `json:"mass,string"`
	Fall        string  `json:"fall"`
	Year        string  `json:"year"`
	RecLat      float32 `json:"reclat,string"`
	RecLong     float32 `json:"reclong,string"`
	GeoLocation struct {
		Latitude  float32 `json:"latitude,string"`
		Longitude float32 `json:"longitude,string"`
	} `json:"geolocation"`
}

func downloadMeteoriteData() (*[]MeteoriteData, error) {
	response, err := resty.New().R().SetResult(&[]MeteoriteData{}).Get(nasaMeteoriteDataUrl)

	if err != nil {
		logger.AppServiceLog.Errorw("NASA meteorite data fetch error", "error", err)
		return nil, err
	}

	result, ok := response.Result().(*[]MeteoriteData)

	if !ok {
		logger.AppServiceLog.Errorw("Meteorite data type assertion error", "error", err)
		return nil, errors.New("meteorite data type assertion error")
	}

	return result, err
}

type googleReverseGeoCodeInfo struct {
	PlusCode struct {
		GlobalCode string `json:"global_code"`
	} `json:"plus_code"`
	Results []struct {
		AddressComponents []struct {
			LongName  string   `json:"long_name"`
			ShortName string   `json:"short_name"`
			Types     []string `json:"types"`
		} `json:"address_components"`
		FormattedAddress string   `json:"formatted_address"`
		PlaceID          string   `json:"place_id"`
		Types            []string `json:"types"`
	} `json:"results"`
	Status string `json:"status"`
}

func ProcessMeteoriteData() error {
	meteoriteData, err := downloadMeteoriteData()

	if err != nil {
		logger.AppServiceLog.Errorw("Meteorite data download error", "error", err)
		return err
	}

	meteoriteDataLength := len(*meteoriteData)

	for i := 0; i < meteoriteDataLength/50; i += 1 {
		dataSlice := (*meteoriteData)[i*50 : (i+1)*50]
		loadMeteoriteData(&dataSlice)
	}

	return nil
}

func loadMeteoriteData(data *[]MeteoriteData) error {
	var payload []db.MeteoriteData

	for _, data := range *data {
		var year *uint = nil
		if data.Year != "" && len(data.Year) >= 4 {
			time, err := time.Parse("2006-01-02T15:04:05.000", data.Year)

			if err == nil {
				intYear := uint(time.Year())
				year = &intYear
			}
		}

		res, err := resty.New().SetQueryParams(map[string]string{ // sample of those who use this manner
			"key":         configs.GoogleMapApiKey,
			"latlng":      fmt.Sprintf("%f,%f", data.GeoLocation.Latitude, data.GeoLocation.Longitude),
			"result_type": "country|administrative_area_level_1",
		}).R().SetResult(&googleReverseGeoCodeInfo{}).Get("https://maps.googleapis.com/maps/api/geocode/json")

		if err != nil {
			logger.AppServiceLog.Errorw("Google reverse geocode fetch error", "error", err)
			return err
		}

		result := res.Result()

		reverseGeoCodeInfo, ok := result.(*googleReverseGeoCodeInfo)

		if !ok {
			logger.AppServiceLog.Errorw("Google reverse geocode type assertion error")
			return errors.New("google reverse geocode type assertion error")
		}

		if resultLength := len(reverseGeoCodeInfo.Results); resultLength > 0 {
			length := len(reverseGeoCodeInfo.Results[0].AddressComponents)

			if length >= 2 && reverseGeoCodeInfo.Results[0].AddressComponents[length-1].LongName != "" && reverseGeoCodeInfo.Results[0].AddressComponents[length-2].LongName != "" {
				payload = append(payload, db.MeteoriteData{
					Name:          data.Name,
					Mass:          data.Mass,
					DiscoveryType: strings.ToLower(data.Fall),
					Year:          year,
					Latitude:      data.GeoLocation.Latitude,
					Longitude:     data.GeoLocation.Longitude,
					Region:        reverseGeoCodeInfo.Results[0].AddressComponents[length-2].LongName,
					Country:       reverseGeoCodeInfo.Results[0].AddressComponents[length-1].LongName,
				})
			}
		}

	}

	tx := db.DB.Clauses(clause.OnConflict{DoNothing: true}).CreateInBatches(&payload, 50)

	if tx.Error != nil {
		logger.AppServiceLog.Errorw("Meteorite data load error", "error", tx.Error)
		return tx.Error
	}

	logger.AppServiceLog.Infow("Meteorite data loaded", "rowsAffected", tx.RowsAffected)
	return nil
}
