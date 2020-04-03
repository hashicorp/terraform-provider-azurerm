package validate

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
)

func DatashareTags(i interface{}, k string) (warnings []string, errors []error) {
	tagsMap := i.(map[string]interface{})

	if len(tagsMap) > 50 {
		errors = append(errors, fmt.Errorf("a maximum of 50 tags can be applied to each ARM resource"))
	}

	for k, v := range tagsMap {
		if len(k) > 512 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag key is 512 characters: %q has %d characters", k, len(k)))
			return warnings, errors
		}

		if strings.ToLower(k) != k {
			errors = append(errors, fmt.Errorf("a tag key %q expected to be all in lowercase", k))
			return warnings, errors
		}

		value, err := tags.TagValueToString(v)
		if err != nil {
			errors = append(errors, err)
			return warnings, errors
		}
		if len(value) > 256 {
			errors = append(errors, fmt.Errorf("the maximum length for a tag value is 256 characters: the value for %q has %d characters", k, len(value)))
			return warnings, errors
		}
	}

	return warnings, errors
}

func DataShareAccountName() schema.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^[^<>%&:\\?/#*$^();,.\|+={}\[\]!~@]{3,90}$`), `Data share account name should have length of 3 - 90, and cannot contain <>%&:\?/#*$^();,.|+={}[]!~@.`,
	)
}
