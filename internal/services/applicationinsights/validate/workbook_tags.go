// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"errors"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
)

func WorkbookTags(i any, k string) (warnings []string, errs []error) {
	warnings, errs = tags.Validate(i, k)
	if len(errs) > 0 {
		return
	}

	tagsMap := i.(map[string]any)
	if _, ok := tagsMap["hidden-title"]; ok {
		errs = append(errs, errors.New("a tag with the key `hidden-title` should not be used to set the display name. Please Use `display_name` instead"))
	}

	return
}
