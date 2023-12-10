package controllers

import (
	"fmt"
	apperrors "main/internal/app/appErrors"
	"main/internal/app/utils"
	"main/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LoginRequestBody struct {
	User     string `json:"user" validate:"required,min:8,max:16"`
	Password string `json:"password" validate:"required,min:8,max:16"`
}

func CreateUserSession(ginContext *gin.Context) {
	var body LoginRequestBody
	if err := ginContext.ShouldBindJSON(&body); err != nil {
		errResponse := apperrors.New("E003")
		ginContext.AbortWithStatusJSON(errResponse.HttpStatusCode, errResponse)
		return
	}
	user, userJwt, csrfToken, maxAge, cubeJwt, errorResponse := services.VerifyUserCredentials(body.User, body.Password)

	if errorResponse != nil {
		ginContext.AbortWithStatusJSON(errorResponse.HttpStatusCode, errorResponse)
		return
	}

	fmt.Println(userJwt)
	fmt.Println(cubeJwt)

	ginContext.SetSameSite(http.SameSiteStrictMode)
	ginContext.SetCookie("__xauth", userJwt, int(maxAge), "/", "", false, true)

	ginContext.JSON(http.StatusCreated, utils.GenerateSuccessResponse(map[string]interface{}{
		"user":         user,
		"cubeApiToken": cubeJwt,
		"csrfToken":    csrfToken,
	}))

}
