package services

import (
	"encoding/base64"
	"main/configs"
	apperrors "main/internal/app/appErrors"
	"main/internal/db"
	"main/internal/services/jwt"
	"main/internal/services/logger"

	"golang.org/x/crypto/bcrypt"
)

func VerifyUserCredentials(username string, password string) (*db.Users, string, string, *apperrors.ErrorResponse) {
	user := &db.Users{
		UserName: username,
		Active:   true,
	}
	tx := db.DB.Find(&user)

	if tx.Error != nil {
		logger.AppServiceLog.Errorw("Error while fetching user", "error", tx.Error)
		errResponse := apperrors.New("E001")
		return nil, "", "", &errResponse
	}

	if tx.RowsAffected == 0 {
		logger.AppServiceLog.Errorw("No active user found", "username", username)
		errResponse := apperrors.New("E004")
		return nil, "", "", &errResponse
	}

	storedPassword, err := base64.StdEncoding.DecodeString(user.Password)

	if err != nil {
		logger.AppServiceLog.Errorw("Error while decoding password", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", &errResponse
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))

	if err != nil {
		logger.AppServiceLog.Infow("Comparing password failed", "username", username)
		errResponse := apperrors.New("E004")
		return nil, "", "", &errResponse
	}

	logger.AppServiceLog.Infow("User credentials validated", "username", username)

	userJwt, _, err := jwt.Sign(jwt.GenerationData{
		Payload: jwt.Payload{
			"username": username,
		},
		Key: configs.APIJWTSecret,
	})

	if err != nil {
		logger.AppServiceLog.Errorw("Error generating user JWT", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", &errResponse
	}

	cubeJwt, _, err := jwt.Sign(jwt.GenerationData{
		Key: configs.CubeAPISecret,
	})

	if err != nil {
		logger.AppServiceLog.Errorw("Error generating cube JWT", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", &errResponse
	}

	return user, userJwt, cubeJwt, nil
}
