package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/the4star/reservation-system/internal/config"
	"github.com/the4star/reservation-system/internal/models"
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

func sendEmail(app *config.AppConfig, to, from, subject, content, template string) {
	msg := models.MailData{
		To:       to,
		From:     from,
		Subject:  subject,
		Content:  content,
		Template: template,
	}

	app.MailChan <- msg
}
