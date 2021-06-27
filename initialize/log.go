package initialize

import (
	"github.com/sirupsen/logrus"
	"io"
	"ocpp-client/log"
	"os"
)

func initLog() {
	log.Logger = logrus.New()
	file, err := os.OpenFile(os.Getenv("WIN_LOG_PATH"), os.O_CREATE|os.O_RDWR|os.O_APPEND, os.ModeAppend|os.ModePerm)
	if err != nil {
		panic(err)
	}
	log.Logger.SetOutput(io.MultiWriter(os.Stdout, file))
	log.Logger.SetLevel(logrus.DebugLevel)
	log.Logger.SetFormatter(&logrus.JSONFormatter{})
	//log.Logger.SetReportCaller(true)
	// elasticsearch hook
	//log.Logger.AddHook(hooks.NewEsHook(hooks.NewEsCfg()))
}
