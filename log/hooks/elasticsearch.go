package hooks

import (
	"context"
	"fmt"
	"github.com/olivere/elastic"
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"strings"
	"time"
)

type EsCfg struct {
	esAddrs []string // es地址(目前默认是单机)
}

func NewEsCfg() EsCfg {
	return EsCfg{
		esAddrs: []string{os.Getenv("ES_ADDR")},
	}
}

type esHook struct {
	cmd    string
	client *elastic.Client
}

func (e *esHook) Levels() []logrus.Level {
	return []logrus.Level{logrus.PanicLevel, logrus.FatalLevel, logrus.ErrorLevel}
}

func (e *esHook) Fire(entry *logrus.Entry) error {
	doc := newEsLog(entry)
	go e.sendES(doc)
	return nil
}

type appLogDocModel map[string]interface{}

func (e *esHook) sendES(doc appLogDocModel) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("send entry to es failed: ", r)
		}
	}()
	_, err := e.client.Index().Index(doc.indexName()).Type("_doc").BodyJson(doc).Do(context.Background())
	if err != nil {
		return
	}
}

func NewEsHook(cc EsCfg) logrus.Hook {
	es, err := elastic.NewClient(
		elastic.SetURL(cc.esAddrs...),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(15*time.Second),
		elastic.SetErrorLog(log.New(os.Stderr, "ES:", log.LstdFlags)))
	if err != nil {
		panic(err)
	}
	return &esHook{client: es, cmd: strings.Join(os.Args, " ")}
}

func newEsLog(e *logrus.Entry) appLogDocModel {
	ins := make(map[string]interface{})
	for k, v := range e.Data {
		ins[k] = v
	}
	ins["time"] = time.Now().Local()
	ins["level"] = e.Level
	ins["message"] = e.Message
	ins["caller"] = fmt.Sprintf("%s:%d %#v", e.Caller.File, e.Caller.Line, e.Caller.Function)
	return ins
}

func (a *appLogDocModel) indexName() string {
	return "ocpp-client-" + time.Now().Local().Format("2006-01-02")
}
