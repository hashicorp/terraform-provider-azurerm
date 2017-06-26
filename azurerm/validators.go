package azurerm

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/satori/uuid"
)

func evaluateSchemaValidateFunc(i interface{}, k string, validateFunc schema.SchemaValidateFunc) (bool, error) {
	_, es := validateFunc(i, k)

	if len(es) > 0 {
		return false, es[0]
	}

	return true, nil
}

func validateUUID(v interface{}, k string) (ws []string, errors []error) {
	if _, err := uuid.FromString(v.(string)); err != nil {
		errors = append(errors, fmt.Errorf("%q is an invalid UUUID: %s", k, err))
	}
	return
}
