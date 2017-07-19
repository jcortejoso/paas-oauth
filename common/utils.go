package common

import (
	"net/http"
	"regexp"
	"os"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type HttpError struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      int    `json:"-"`
}

func NewHttpError(description string, status int) *HttpError {
	return &HttpError{
		Title:       http.StatusText(status),
		Description: description,
		Status:      status,
	}
}

var (
	emailRe = regexp.MustCompile(`^[a-z0-9“”._%+-]+@(?:[a-z0-9-\[]+\.)+[a-z0-9-\]]{2,}$`)
)

func ValidateEmail(email string) bool {
	return emailRe.MatchString(email)
}

const (
	uiconfigJson = "/opt/mesosphere/etc/ui-config.json"
)

func OpenDcosConfig(cfg *map[string]interface{}) *HttpError {
	f, err := os.Open(uiconfigJson)
	if err != nil {
		return NewHttpError("ui-config.json read failed", http.StatusInternalServerError)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(cfg)
	if err != nil {
		log.Printf("Decode: %v", err)
		return NewHttpError("JSON decode error", http.StatusInternalServerError)
	}
	return nil
}
