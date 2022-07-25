package filestorage

import (
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type Config struct {
	Path string `json:"path"`
}

func NewConfig(config interface{}) (*Config, error) {
	c := new(Config)

	tmp := config.(map[string]interface{})

	for key, val := range tmp {
		reflect.ValueOf(c).Elem().FieldByName(strings.Title(key)).SetString(val.(string))
	}

	logrus.Info(c)
	return c, nil
}
