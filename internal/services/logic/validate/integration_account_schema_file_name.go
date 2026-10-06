// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IntegrationAccountSchemaFileName() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`\.xsd$`), "must end with `.xsd`")
}
