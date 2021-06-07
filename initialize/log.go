package initialize

import (
	"github.com/sirupsen/logrus"
	"ocpp-client/log"
	"os"
)

func init() {
	log.Logger = logrus.New()
	log.Logger.SetOutput(os.Stdout)
	log.Logger.SetLevel(logrus.DebugLevel)
	log.Logger.SetFormatter(&logrus.JSONFormatter{})
}
