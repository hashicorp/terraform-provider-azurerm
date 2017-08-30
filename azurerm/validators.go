package azurerm

import (
	"fmt"
	"time"

	"github.com/satori/uuid"
	"github.com/Azure/go-autorest/autorest/date"
)

func validateRFC3339Date(v interface{}, k string) (ws []string, errors []error) {
	dateString := v.(string)

	if _, err := date.ParseTime(time.RFC3339, dateString); err != nil {
		errors = append(errors, fmt.Errorf("`%s` is an invalid RFC3339 date: %+v", k, err))
	}

	return
}

func validateUUID(v interface{}, k string) (ws []string, errors []error) {
	if _, err := uuid.FromString(v.(string)); err != nil {
		errors = append(errors, fmt.Errorf("%q is an invalid UUUID: %s", k, err))
	}
	return
}
