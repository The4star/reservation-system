package handlers

import (
	"encoding/json"
	"net/http"
)

func internalServerErrorJSON(w http.ResponseWriter) {
	resp := roomAvailabilityResponse{
		OK:  false,
		Msg: "Internal server error",
	}

	responseData, _ := json.Marshal(resp)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(responseData)
}
