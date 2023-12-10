package utils

func GenerateSuccessResponse(data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"data": data,
	}
}
