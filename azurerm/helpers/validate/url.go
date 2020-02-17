package validate

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

// deprecated use validation.IsURLWithHTTPS instead
func URLIsHTTPS(i interface{}, k string) (_ []string, errors []error) {
	return validation.IsURLWithHTTPS(i, k)
}

// todo ad to sdk
// deprecated use validation.IsURLWithScheme instead
func URLIsHTTPOrHTTPS(i interface{}, k string) (_ []string, errors []error) {
	return validation.IsURLWithScheme([]string{"http", "https"})(i, k)
}
