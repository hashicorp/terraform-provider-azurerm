// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BackupPolicyCosmosdbAccountFullBackupSchedule() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^R/\d{4}-\d{2}-\d{2}T(?:\d{2}:\d{2}|\d{2}:\d{2}:\d{2}(?:\.\d{3})?)(?:Z|[+-]\d{2}:\d{2})/P1W$`),
		"must use the format `R/YYYY-MM-DDThh:mm[Z|(+/-)hh:mm]/P1W`, `R/YYYY-MM-DDThh:mm:ss[Z|(+/-)hh:mm]/P1W`, or `R/YYYY-MM-DDThh:mm:ss.fff[Z|(+/-)hh:mm]/P1W`",
	)
}
