package controllers

import (
	apperrors "main/internal/app/appErrors"
	"main/internal/db"
	"main/internal/services"
	"main/internal/services/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

const MeteoriteData = "meteorite"

func HandleMeteoriteData(ginContext *gin.Context) {
	value, ok := ginContext.Get("user")
	if !ok {
		logger.AppControllerLog.Errorw("Error getting user from gin context")
		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	user, ok := value.(db.Users)

	if !ok {
		logger.AppControllerLog.Errorw("Error asserting user from gin context")
		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	request, err := services.AccquireLockForIngestion(MeteoriteData, user.ID)

	if err != nil {
		logger.AppControllerLog.Errorw("Error while accquiring lock", "error", err)

		if err == services.ErrDataIngestionNotEligible {
			errResponse := apperrors.New("E007")
			ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
			return
		}

		errResponse := apperrors.New("E001")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}

	go services.ProcessMeteoriteData(request.ID)

	ginContext.Status(http.StatusAccepted)
}
