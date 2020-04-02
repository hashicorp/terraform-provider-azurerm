package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func PeerAsnName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\w]+$`), "PeerAsnName has invalid characters. Allowed characters are a-z, A-Z, 0-9, _")
}

func PeerName() func(i interface{}, k string) (warnings []string, errors []error) {
	return validation.StringMatch(regexp.MustCompile(`^[\w-. ]+$`), "PeerName has invalid characters. Allowed characters are a-z, A-Z, 0-9,  , ., -, _")
}
