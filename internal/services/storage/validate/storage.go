package validate

import (
	"fmt"
	"regexp"
)

func StorageShareDirectoryName(v interface{}, k string) (warnings []string, errors []error) {
	value := v.(string)
	// Per: https://learn.microsoft.com/en-us/rest/api/storageservices/naming-and-referencing-shares--directories--files--and-metadata#directory-and-file-names
	if regexp.MustCompile(`^\.+$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%s must not only contain dots`, k))
	}
	if !regexp.MustCompile(`^[^"\/:|<>*?]{1,255}$`).MatchString(value) {
		errors = append(errors, fmt.Errorf(`%s must not contain following characters: "\/:|<>*?`, k))
	}

	return warnings, errors
}
