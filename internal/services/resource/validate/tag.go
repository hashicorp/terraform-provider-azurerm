package validate

import (
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
)

func TagKey(i interface{}, k string) ([]string, []error) {
	v := i.(string)

	return tags.Validate(map[string]interface{}{v: "value"}, k)
}

func TagValue(v interface{}, k string) ([]string, []error) {
	return tags.Validate(map[string]interface{}{"key": v}, k)
}
