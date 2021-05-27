package validate

import (
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
)

func UpgradeTimeout(v interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(
		regexp.MustCompile(`^(\d+\.)?([0-1][0-9]|[2][0-3]):[0-5][0-9]:[0-5][0-9](\.[0-9][0-9])?$`),
		"The timeout duration must be in this format [d.]hh:mm:ss[.ms].",
	)(v, k)
}
