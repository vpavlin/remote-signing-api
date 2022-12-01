package filestorage

import (
	"reflect"
	"strings"
)

type Config struct {
	Path string `json:"path"`
}

func NewConfig(config interface{}) (*Config, error) {
	c := new(Config)

	tmp := config.(map[string]interface{})

	for key, val := range tmp {
		field := reflect.ValueOf(c).Elem().FieldByName(strings.Title(key))
		if field != (reflect.Value{}) {
			field.SetString(val.(string))
		}
	}

	return c, nil
}
