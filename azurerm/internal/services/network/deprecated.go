package network

import (
	"fmt"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
)

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

func evaluateSchemaValidateFunc(i interface{}, k string, validateFunc pluginsdk.SchemaValidateFunc) (bool, error) { // nolint: unparam
	_, errors := validateFunc(i, k)

	errorStrings := []string{}
	for _, e := range errors {
		errorStrings = append(errorStrings, e.Error())
	}

	if len(errors) > 0 {
		return false, fmt.Errorf(strings.Join(errorStrings, "\n"))
	}

	return true, nil
}
