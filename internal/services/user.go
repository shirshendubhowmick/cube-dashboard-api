package services

import (
	"encoding/base64"
	apperrors "main/internal/app/appErrors"
	"main/internal/db"
	"main/internal/services/logger"

	"golang.org/x/crypto/bcrypt"
)

func VerifyUserCredentials(username string, password string) (string, *apperrors.ErrorResponse) {
	user := &db.Users{
		UserName: username,
		Active:   true,
	}
	tx := db.DB.Find(&user)

	if tx.Error != nil {
		logger.AppServiceLog.Errorw("Error while fetching user", "error", tx.Error)
		errResponse := apperrors.New("E001")
		return "", &errResponse
	}

	if tx.RowsAffected == 0 {
		logger.AppServiceLog.Errorw("No active user found", "username", username)
		errResponse := apperrors.New("E004")
		return "", &errResponse
	}

	storedPassword, err := base64.StdEncoding.DecodeString(user.Password)

	if err != nil {
		logger.AppServiceLog.Errorw("Error while decoding password", "error", err)
		errResponse := apperrors.New("E001")
		return "", &errResponse
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))

	if err != nil {
		logger.AppServiceLog.Infow("Comparing password failed", "username", username)
		errResponse := apperrors.New("E004")
		return "", &errResponse
	}

	logger.AppServiceLog.Infow("User credentials validated", "username", username)

	// bcrypt.CompareHashAndPassword([]byte("$2a$10$1Z0Z1Z0Z1Z0Z1Z0Z1Z0Z1OuZ1Z0Z1Z0Z1Z0Z1Z0Z1Z0Z1Z0Z1Z0Z1"), []byte(password))
	return "", nil
}
