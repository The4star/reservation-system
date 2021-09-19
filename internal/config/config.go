package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/the4star/reservation-system/internal/models"
)

// AppConfig holds the application config.
type AppConfig struct {
	InProduction  bool
	UseCache      bool
	TemplateCache map[string]*template.Template
	Session       *scs.SessionManager
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	MailChan      chan models.MailData
}
