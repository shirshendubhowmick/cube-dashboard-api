package services

import (
	"errors"
	"main/internal/db"
	"main/internal/services/logger"
	"time"

	"github.com/google/uuid"
)

var ErrLockAcquireFailed = errors.New("failed to acquire lock")
var ErrDataIngestionNotEligible = errors.New("data ingestion not eligible")

func AccquireLockForIngestion(dataName string, requestedBy uuid.UUID) (*db.DataIngestionRequests, error) {
	tx := db.DB.Begin()

	result := tx.Exec("LOCK TABLE data_ingestion_requests IN ACCESS EXCLUSIVE MODE")

	if result.Error != nil {
		tx.Rollback()
		logger.AppServiceLog.Infow("Unable to accquire lock", "error", result.Error)
		return nil, ErrLockAcquireFailed
	}

	// Check eligibility if not inserted in last one hour
	var count int64
	result = tx.Model(&db.DataIngestionRequests{}).Where("data_name = ?", dataName).Where("updated_at > ? or in_progress = ?", time.Now().Add(-1*time.Hour), true).Count(&count)

	if result.Error != nil {
		tx.Rollback()
		logger.AppServiceLog.Errorw("Error while checking recent data ingestion request", "error", result.Error)
		return nil, result.Error
	}

	if count > 0 {
		tx.Rollback()
		logger.AppServiceLog.Infow("Data ingestion request already in progress", "dataName", dataName)
		return nil, ErrDataIngestionNotEligible
	}

	newRequest := &db.DataIngestionRequests{
		DataName:    dataName,
		InProgress:  true,
		RequestedBy: requestedBy,
	}

	result = tx.Model(&db.DataIngestionRequests{}).Create(newRequest)

	if result.Error != nil {
		tx.Rollback()
		logger.AppServiceLog.Errorw("Error while inserting data ingestion request", "error", result.Error)
		return nil, result.Error
	}

	result = tx.Commit()
	if result.Error != nil {
		logger.AppServiceLog.Errorw("Error while committing transaction", "error", result.Error)
		tx.Rollback()
		return nil, result.Error
	}

	logger.AppServiceLog.Infow("Lock accquired")
	return newRequest, nil
}

func MarkIngestionComplete(requestId uuid.UUID) bool {
	var lockedRequest db.DataIngestionRequests

	tx := db.DB.Begin()

	if err := tx.Model(&db.DataIngestionRequests{}).Where("id = ?", requestId).Set("gorm:query_option", "FOR UPDATE").First(&lockedRequest).Error; err != nil {
		logger.AppServiceLog.Errorw("Error while accquiring lock on data ingestion request", "error", err, "id", requestId)
		return false
	}

	result := tx.Model(&db.DataIngestionRequests{}).Where("id = ?", requestId).Update("in_progress", false)

	if result.Error != nil {
		logger.AppServiceLog.Errorw("Error while updating data ingestion request", "error", result.Error, "id", requestId)
		tx.Rollback()
		return false
	}

	if result.RowsAffected == 0 {
		logger.AppServiceLog.Errorw("No data ingestion request found", "id", requestId)
		tx.Rollback()
		return false
	}

	result = tx.Commit()

	if result.Error != nil {
		logger.AppServiceLog.Errorw("Error while committing transaction", "error", result.Error)
		tx.Rollback()
		return false
	}

	return true
}
