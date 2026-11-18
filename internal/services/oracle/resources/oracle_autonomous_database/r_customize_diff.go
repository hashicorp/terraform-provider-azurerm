// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (Resource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Second,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff == nil {
				return nil
			}

			rawConfig := metadata.ResourceDiff.GetRawConfig()
			if rawConfig.IsNull() || !rawConfig.IsKnown() {
				return nil
			}

			workload := rawConfig.GetAttr("db_workload")
			if !workload.IsKnown() || workload.IsNull() || workload.AsString() != string(autonomousdatabases.WorkloadTypeAJD) {
				return nil
			}

			// Check configuration rather than the diff so omitted arguments are not
			// mistaken for explicitly configured schema defaults. Unknown values
			// are validated when Terraform resolves them during apply planning.
			for _, name := range []string{"character_set", "national_character_set"} {
				value := rawConfig.GetAttr(name)
				if value.IsKnown() && !value.IsNull() && value.AsString() != "" {
					return fmt.Errorf("`%s` must be omitted when `db_workload` is `AJD`", name)
				}
			}
			return nil
		},
	}
}
