package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func LogzMonitorName(v interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^[\w\-]{1,32}$`),
		`name must be between 1 and 32 characters in length and may contain only letters, numbers, hyphens and underscores`,
	)(v, k)
}
