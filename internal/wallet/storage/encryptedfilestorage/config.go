package encryptedfilestorage

import (
	"log"
	"reflect"
	"strings"
)

type Config struct {
	Path            string `json:"path"`
	PasswordEnvName string `json:"passwordEnvName"`
}

func NewConfig(config interface{}) (*Config, error) {
	c := new(Config)

	tmp := config.(map[string]interface{})

	log.Println(tmp)

	for key, val := range tmp {
		reflect.ValueOf(c).Elem().FieldByName(strings.Title(key)).SetString(val.(string))
	}

	return c, nil
}
