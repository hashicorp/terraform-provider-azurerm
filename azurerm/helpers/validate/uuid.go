package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

//uuid regex helper
var UUIDRegExp = regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")

// deprecated use validation.IsUUID instead
func UUID(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsUUID(i, k)
}

func GUID(i interface{}, k string) (warnings []string, errors []error) {
	return validation.IsUUID(i, k)
}

// deprecated use validation.Any(validation.IsUUID, validation.StringIsEmpty) instead
func UUIDOrEmpty(i interface{}, k string) (warnings []string, errors []error) {
	return validation.Any(validation.IsUUID, validation.StringIsEmpty)(i, k)
}
