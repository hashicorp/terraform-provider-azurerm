package validate

import (
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func CdnFrontdoorEndpointName(i interface{}, k string) (_ []string, errors []error) {
	regex := regexp.MustCompile(`(^[\da-zA-Z])([-\da-zA-Z]{0,44})([\da-zA-Z]$)`)
	message := fmt.Sprintf(`%q must be between 2 and 46 characters in length and begin with a letter or number, end with a letter or number and may contain only letters, numbers or hyphens`, k)
	return validation.StringMatch(regex, message)(i, k)
}
