package filestorage

import (
	"log"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Path string `json:"path"`
}

func NewConfig(config interface{}) (*Config, error) {
	c := new(Config)

	log.Println(config)

	tmp := config.(map[string]interface{})

	for key, val := range tmp {
		logrus.Info(strings.Title(key), val)
		reflect.ValueOf(c).Elem().FieldByName(strings.Title(key)).SetString(val.(string))
	}

	logrus.Info(c)
	return c, nil
}
