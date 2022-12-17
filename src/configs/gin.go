package configs

import (
	log "github.com/sirupsen/logrus"
	"os"
)

func GetAppPort() string {
	appPort, ok := os.LookupEnv("PORT")
	if !ok {
		log.Errorf("Missing PORT env, defaulting to 3000")
		appPort = "3000"
	}
	return appPort
}
