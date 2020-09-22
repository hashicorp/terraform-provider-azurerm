package validate

import (
	"fmt"
	"regexp"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/storage/parsers"
)

func StorageSyncName(v interface{}, _ string) (warnings []string, errors []error) {
	input := v.(string)

	if !regexp.MustCompile("^[0-9a-zA-Z-_. ]*[0-9a-zA-Z-_]$").MatchString(input) {
		errors = append(errors, fmt.Errorf("name (%q) can only consist of letters, numbers, spaces, and any of the following characters: '.-_' and that does not end with characters: '. '", input))
	}

	return warnings, errors
}

func StorageSyncId(i interface{}, k string) (warnings []string, errors []error) {
	v, ok := i.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected type of %q to be string", k))
		return warnings, errors
	}

	if _, err := parsers.ParseStorageSyncID(v); err != nil {
		errors = append(errors, fmt.Errorf("can not parse %q as a Storage Sync resource id: %v", k, err))
	}

	return warnings, errors
}
