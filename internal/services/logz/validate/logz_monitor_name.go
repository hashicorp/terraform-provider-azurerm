package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LogzMonitorName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w\-]{1,32}$`),
		`The name length must be from 1 to 32 characters. The name can only contain letters, numbers, hyphens and underscore.`,
	)(v, k)
}
