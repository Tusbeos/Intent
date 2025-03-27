package middleware

import (
	"encoding/json"
)

// extractUserID lấy `UserID` từ response JSON
func extractUserID(rec *CustomResponseRecorder) int {
	var resBody map[string]interface{}
	json.Unmarshal(rec.Body.Bytes(), &resBody)

	if userID, exists := resBody["user_id"].(float64); exists {
		return int(userID)
	}

	return 0
}
