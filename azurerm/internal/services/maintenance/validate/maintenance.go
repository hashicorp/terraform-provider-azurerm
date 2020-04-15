package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func TagsWithLowerCaseKey(v interface{}, k string) (warnings []string, errors []error) {
	warnings, errors = tags.Validate(v, k)

	tagsMap := v.(map[string]interface{})
	for key := range tagsMap {
		for _, c := range key {
			if c >= 'A' && c <= 'Z' {
				errors = append(errors, fmt.Errorf("the key of %q can not contain upper case letter. The key %q has upper case letter %q", k, key, c))
			}
		}
	}

	return
}
