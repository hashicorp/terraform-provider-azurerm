// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func IntegrationAccountAssemblyVersion() pluginsdk.SchemaValidateFunc {
	return func(i interface{}, k string) (warnings []string, errors []error) {
		v, ok := i.(string)
		if !ok {
			errors = append(errors, fmt.Errorf("expected %q to be a string", k))
			return
		}

		if !regexp.MustCompile(`^([0-9]+.[0-9]+.[0-9]+.[0-9]+)$|^([0-9]+.[0-9]+)$`).MatchString(v) {
			errors = append(errors, fmt.Errorf("%q must be in the format `major.minor.build.revision` in which `build` and `revision` components are optional", k))
			return
		}

		return
	}
}
