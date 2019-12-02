package sql

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

func ValidateUUIdString(val interface{}, key string) (warnings []string, errors []error) {
	v := val.(string)
	var _, err = uuid.FromString(v)
	if err != nil {
		errors = append(errors, fmt.Errorf("%q is not in correct format:%+v", key, err))
	}
	return
}
