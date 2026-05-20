// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"errors"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

// NOTE: these methods are deprecated, but provided to ease compatibility for open PR's

func evaluateSchemaValidateFunc(i interface{}, k string, validateFunc pluginsdk.SchemaValidateFunc) (bool, error) { // nolint: unparam
	_, errs := validateFunc(i, k)

	errorStrings := []string{}
	for _, e := range errs {
		errorStrings = append(errorStrings, e.Error())
	}

	if len(errs) > 0 {
		return false, errors.New(strings.Join(errorStrings, "\n"))
	}

	return true, nil
}
