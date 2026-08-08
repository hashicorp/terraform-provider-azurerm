// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func StorageShareDirectoryName(v interface{}, k string) ([]string, []error) {
	// Per: https://learn.microsoft.com/en-us/rest/api/storageservices/naming-and-referencing-shares--directories--files--and-metadata#directory-and-file-names
	return validation.All(
		validation.StringDoesNotMatch(regexp.MustCompile(`^\.+$`), "must not only contain dots"),
		// Note that we didn't forbid the forward slash in the non-head segment here as it seems to be allowed as the directory name for constructing directory hierarchy.
		validation.StringMatch(regexp.MustCompile(`^[^"/\:|<>*?]+(/[^"\:|<>*?]+)*$`), `must not contain following characters: "\:|<>*?`),
	)(v, k)
}
