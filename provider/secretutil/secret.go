package secretutil

import (
	"errors"

	"github.com/tidwall/gjson"
)

func GetValueFromJSON(secret []byte, property string) ([]byte, error) {
	result := gjson.Get(string(secret), property)
	if !result.Exists() {
		return nil, errors.New("property not found in secret")
	}

	return []byte(result.String()), nil
}
