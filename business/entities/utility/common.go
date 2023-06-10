package utility

import (
	"encoding/json"
)

// maps obj1 to obj2
func MapObjectToAnother(obj1 interface{}, obj2 interface{}) error {
	b, err := json.Marshal(obj1)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, obj2)
	return err
}
