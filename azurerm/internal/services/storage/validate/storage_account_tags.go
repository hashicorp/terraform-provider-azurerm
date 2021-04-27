package validate

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func StorageAccountTags(v interface{}, _ string) (warnings []string, errors []error) {
	tagsMap := v.(map[string]interface{})

	if len(tagsMap) > 50 {
		errors = append(errors, fmt.Errorf("a maximum of 50 tags can be applied to storage account ARM resource"))
	}

	for k, v := range tagsMap {
		if len(k) > 128 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag key is 128 characters: %q is %d characters", k, len(k)))
		}

		value, err := tags.TagValueToString(v)
		if err != nil {
			errors = append(errors, err)
		} else if len(value) > 256 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q is %d characters", k, len(value)))
		}
	}

	return warnings, errors
}
