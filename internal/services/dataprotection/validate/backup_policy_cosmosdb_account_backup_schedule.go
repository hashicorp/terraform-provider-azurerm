// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"regexp"

	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func BackupPolicyCosmosdbAccountBackupSchedule() pluginsdk.SchemaValidateFunc {
	return validation.StringMatch(
		regexp.MustCompile(`^R/\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:Z|[+-]\d{2}:\d{2})/P1W$`),
		"must use the format `R/YYYY-MM-DDThh:mm:ss[Z|(+/-)hh:mm]/P1W`; seconds are required and fractional seconds are not supported",
	)
}
