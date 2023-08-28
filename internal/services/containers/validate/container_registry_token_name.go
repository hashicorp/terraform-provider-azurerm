// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ContainerRegistryTokenName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9-]{4,48}$`), fmt.Sprintf("only alpha numeric characters (optionally separated by dash) in length of 5 to 50 are allowed in %q", k))(v, k)
}
