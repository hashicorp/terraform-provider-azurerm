// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package oracle_autonomous_database_cross_region_disaster_recovery

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-azure-sdk/resource-manager/oracledatabase/2025-09-01/autonomousdatabases"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

var _ sdk.ResourceWithCustomImporter = Resource{}

type Resource struct{}

func (Resource) ResourceType() string {
	return "azurerm_oracle_autonomous_database_cross_region_disaster_recovery"
}

func (Resource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return autonomousdatabases.ValidateAutonomousDatabaseID
}

func (r Resource) CustomImporter() sdk.ResourceRunFunc {
	return func(ctx context.Context, metadata sdk.ResourceMetaData) error {
		id, err := autonomousdatabases.ParseAutonomousDatabaseID(metadata.ResourceData.Id())
		if err != nil {
			return err
		}

		resp, err := metadata.Client.Oracle.OracleClient.AutonomousDatabases.Get(ctx, *id)
		if err != nil {
			return fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if resp.Model == nil || resp.Model.Properties == nil {
			return fmt.Errorf("retrieving %s: unable to determine disaster recovery type", id)
		}

		db, ok := resp.Model.Properties.(autonomousdatabases.AutonomousDatabaseCrossRegionDisasterRecoveryProperties)
		if !ok {
			return fmt.Errorf("retrieving %s: was not of type `%s`", id, autonomousdatabases.DataBaseTypeCrossRegionDisasterRecovery)
		}

		if db.RemoteDisasterRecoveryType != autonomousdatabases.DisasterRecoveryTypeAdg {
			return fmt.Errorf("importing %s: expected disaster recovery type `%s`, got `%s`", id, autonomousdatabases.DisasterRecoveryTypeAdg, db.RemoteDisasterRecoveryType)
		}

		return nil
	}
}
