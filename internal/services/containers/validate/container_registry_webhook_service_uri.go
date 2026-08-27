// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func ContainerRegistryWebhookServiceUri(v interface{}, k string) ([]string, []error) {
	return validation.StringMatch(regexp.MustCompile(`^https?://[^\s]+$`), "must start with http:// or https:// and must not contain whitespaces")(v, k)
}
