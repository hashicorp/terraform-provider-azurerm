// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
)

func WorkbookTags(i interface{}, k string) (warnings []string, errors []error) {
	warnings, errors = tags.Validate(i, k)
	if len(errors) > 0 {
		return
	}

	tagsMap := i.(map[string]interface{})
	if _, ok := tagsMap["hidden-title"]; ok {
		errors = append(errors, fmt.Errorf("a tag with the key `hidden-title` should not be used to set the display name. Please Use `display_name` instead"))
	}

	return
}
