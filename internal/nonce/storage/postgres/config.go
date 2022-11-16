package postgres

import (
	"reflect"
	"strings"
)

type Config struct {
	Connection string `json:"connection"`
}

func NewConfig(config interface{}) (*Config, error) {
	c := new(Config)

	tmp := config.(map[string]interface{})

	for key, val := range tmp {
		reflect.ValueOf(c).Elem().FieldByName(strings.Title(key)).SetString(val.(string))
	}

	return c, nil
}
