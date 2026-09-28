// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

// BucketName validates that the given value is a valid Azure NetApp Files Bucket name.
// The bucket name is S3-compatible: 3-63 characters, DNS-compliant, lowercase
// letters/numbers/hyphens/periods, must start and end with a letter or number,
// must not contain consecutive periods or "."- / "-." sequences, and must not
// look like an IPv4 address.
func BucketName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(3, 63),
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9.\-]+$`), "must contain only lowercase letters, numbers, hyphens or periods"),
		validation.StringMatch(regexp.MustCompile(`^[a-z0-9]`), "must start with a lowercase letter or number"),
		validation.StringMatch(regexp.MustCompile(`[a-z0-9]$`), "must end with a lowercase letter or number"),
		validation.StringDoesNotMatch(regexp.MustCompile(`\.\.|\.\-|\-\.`), `must not contain consecutive periods or ".-" / "-." sequences`),
		validation.StringDoesNotMatch(regexp.MustCompile(`^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$`), "must not be formatted as an IPv4 address"),
	)(v, k)
}
