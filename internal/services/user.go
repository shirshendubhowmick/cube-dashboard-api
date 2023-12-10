package services

import (
	"encoding/base64"
	"main/configs"
	apperrors "main/internal/app/appErrors"
	"main/internal/db"
	"main/internal/services/jwt"
	"main/internal/services/logger"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

func VerifyUserCredentials(username string, password string) (*db.Users, string, string, time.Duration, string, *apperrors.ErrorResponse) {
	zeroTime := time.Duration(0)
	user := &db.Users{
		UserName: username,
		Active:   true,
	}
	tx := db.DB.Find(&user)

	if tx.Error != nil {
		logger.AppServiceLog.Errorw("Error while fetching user", "error", tx.Error)
		errResponse := apperrors.New("E001")
		return nil, "", "", zeroTime, "", &errResponse
	}

	if tx.RowsAffected == 0 {
		logger.AppServiceLog.Errorw("No active user found", "username", username)
		errResponse := apperrors.New("E004")
		return nil, "", "", zeroTime, "", &errResponse
	}

	storedPassword, err := base64.StdEncoding.DecodeString(user.Password)

	if err != nil {
		logger.AppServiceLog.Errorw("Error while decoding password", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", zeroTime, "", &errResponse
	}

	err = bcrypt.CompareHashAndPassword(storedPassword, []byte(password))

	if err != nil {
		logger.AppServiceLog.Infow("Comparing password failed", "username", username)
		errResponse := apperrors.New("E004")
		return nil, "", "", zeroTime, "", &errResponse
	}

	logger.AppServiceLog.Infow("User credentials validated", "username", username)

	guid := xid.New()

	userJwt, maxAge, err := jwt.Sign(jwt.GenerationData{
		Payload: jwt.Payload{
			"username":  username,
			"csrfToken": guid.String(),
		},
		Key: configs.APIJWTSecret,
	})

	if err != nil {
		logger.AppServiceLog.Errorw("Error generating user JWT", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", zeroTime, "", &errResponse
	}

	cubeJwt, _, err := jwt.Sign(jwt.GenerationData{
		Key: configs.CubeAPISecret,
	})

	if err != nil {
		logger.AppServiceLog.Errorw("Error generating cube JWT", "error", err)
		errResponse := apperrors.New("E001")
		return nil, "", "", zeroTime, "", &errResponse
	}

	return user, userJwt, guid.String(), maxAge, cubeJwt, nil
}
