// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/parse"
)

func CertificateContactsID(input interface{}, k string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", k))
		return
	}

	if _, err := parse.CertificateContactsID(v); err != nil {
		errors = append(errors, err)
	}

	return
}
