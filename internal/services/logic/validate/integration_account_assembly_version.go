// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func IntegrationAccountAssemblyVersion() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(regexp.MustCompile(`^([0-9]+.[0-9]+.[0-9]+.[0-9]+)$|^([0-9]+.[0-9]+)$`), "must be in the format `major.minor.build.revision` in which `build` and `revision` components are optional")
}
