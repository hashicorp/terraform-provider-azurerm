// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ClusterResourceGroupName(v interface{}, k string) ([]string, []error) {
	return validation.All(
		validation.StringLenBetween(0, 90),
		validation.StringDoesNotMatch(regexp.MustCompile(`\.$`), "may not end with a period"),
		// regex pulled from https://docs.microsoft.com/en-us/rest/api/resources/resourcegroups/createorupdate
		// ARO only allow for lower cases https://github.com/Azure/ARO-RP/blob/e5c40654277c77fe78ba669610ac05774e448683/pkg/frontend/openshiftcluster_putorpatch.go#L189
		validation.StringMatch(regexp.MustCompile(`^[0-9a-z-._()]+$`), "may only contain lowercase alpha characters, digit, dash, underscores, parentheses and periods"),
	)(v, k)
}
