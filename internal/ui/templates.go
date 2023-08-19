package ui

import (
	"html/template"

	"github.com/sirupsen/logrus"
)

var (
	templateFS = "/Users/dasith/Developer/projects/byte3org/openid-server/internal/ui/templates/*.html"
	templates  = template.Must(template.ParseGlob(templateFS))
)

const (
	queryAuthRequestID = "auth_request_id"
)

func errMsg(err error) string {
	if err == nil {
		return ""
	}
	logrus.Error(err)
	return err.Error()
}
