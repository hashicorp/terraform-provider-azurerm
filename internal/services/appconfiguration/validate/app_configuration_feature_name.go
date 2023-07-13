// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func AppConfigurationFeatureName(input interface{}, key string) ([]string, []error) {
	v, ok := input.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", key)}
	}

	if idx := strings.Index(v, "%"); idx != -1 {
		return nil, []error{fmt.Errorf(`character "%%" is not allowed in %q`, key)}
	}

	if idx := strings.Index(v, ":"); idx != -1 {
		return nil, []error{fmt.Errorf(`character ":" is not allowed in %q`, key)}
	}

	return validation.StringIsNotWhiteSpace(input, key)
}

func AppConfigurationFeatureKey(input interface{}, key string) ([]string, []error) {
	v, ok := input.(string)
	if !ok {
		return nil, []error{fmt.Errorf("expected type of %q to be string", key)}
	}

	if idx := strings.Index(v, "%"); idx != -1 {
		return nil, []error{fmt.Errorf(`character "%%" is not allowed in %q`, key)}
	}

	return validation.StringIsNotWhiteSpace(input, key)
}
