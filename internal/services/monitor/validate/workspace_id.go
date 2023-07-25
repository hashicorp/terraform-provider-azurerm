// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func WorkspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	segments := strings.Split(v, "|")
	if len(segments) == 2 {
		if _, err := validation.IsUUID(segments[0], key); err != nil {
			errors = append(errors, fmt.Errorf("expected the <subscription id> in %q to be a valid UUID, but got %q", key, segments[0]))
			return
		}
		if _, err := validation.IsUUID(segments[1], key); err != nil {
			errors = append(errors, fmt.Errorf("expected the <workSpace id> in %q to be a valid UUID, but got %q", key, segments[1]))
			return
		}
	} else {
		errors = append(errors, fmt.Errorf("expected %q in the format {<subscription id}|{workSpace id} but got %q", key, v))
	}

	return
}
