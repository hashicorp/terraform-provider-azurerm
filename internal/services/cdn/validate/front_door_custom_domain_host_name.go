// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	helperValidate "github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
)

func FrontDoorCustomDomainHostName(i interface{}, k string) (_ []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", k)}
	}

	if len(v) > 253 {
		return nil, []error{fmt.Errorf("%q must be a valid fully qualified domain name, got %q", k, v)}
	}

	labels := strings.Split(v, ".")
	if len(labels) < 2 {
		return nil, []error{fmt.Errorf("%q must be a valid fully qualified domain name, got %q", k, v)}
	}

	for index, label := range labels {
		if label == "" || len(label) > 63 {
			return nil, []error{fmt.Errorf("%q must be a valid fully qualified domain name, got %q", k, v)}
		}

		if index == 0 && label == "*" {
			continue
		}

		if m, _ := helperValidate.RegExHelper(label, k, `^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`); !m {
			return nil, []error{fmt.Errorf("%q must be a valid fully qualified domain name, got %q", k, v)}
		}
	}

	return nil, nil
}
