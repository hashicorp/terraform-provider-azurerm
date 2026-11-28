// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package storage_account_table_properties

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/storage/2025-08-01/storageaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/storage/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/jackofallops/giovanni/storage/2023-11-03/table/tables"
)

func (s Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			storageClient := metadata.Client.Storage
			var model AccountTablePropertiesModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			accountID, err := commonids.ParseStorageAccountID(model.StorageAccountId)
			if err != nil {
				return err
			}

			// Get the target account to ensure it supports tables
			account, err := storageClient.ResourceManager.StorageAccounts.GetProperties(ctx, *accountID, storageaccounts.DefaultGetPropertiesOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %+v", *accountID, err)
			}
			if account.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *accountID)
			}

			if account.Model.Sku == nil || account.Model.Sku.Tier == nil || string(account.Model.Sku.Name) == "" {
				return fmt.Errorf("could not read SKU details for %s", *accountID)
			}

			accountTier := *account.Model.Sku.Tier
			accountReplicationTypeParts := strings.Split(string(account.Model.Sku.Name), "_")
			if len(accountReplicationTypeParts) != 2 {
				return fmt.Errorf("could not read SKU replication type for %s", *accountID)
			}
			accountReplicationType := accountReplicationTypeParts[1]

			accountDetails, err := storageClient.GetAccount(ctx, *accountID)
			if err != nil {
				return err
			}
			if accountDetails == nil {
				return fmt.Errorf("unable to locate %s", *accountID)
			}

			supportLevel := helpers.AvailableFunctionalityForAccount(accountDetails.Kind, accountTier, accountReplicationType)

			if !supportLevel.SupportTable {
				return fmt.Errorf("account %s does not support tables", *accountID)
			}

			client, err := storageClient.TablesDataPlaneClient(ctx, *accountDetails, storageClient.DataPlaneOperationSupportingAnyAuthMethod())
			if err != nil {
				return fmt.Errorf("creating Tables Data Plane Client for %s: %+v", accountID, err)
			}

			props := DefaultValueForAccountTableProperties()
			expandTablePropertiesModel(&model, &props)

			if err = client.UpdateServiceProperties(ctx, props); err != nil {
				return fmt.Errorf("updating Table Properties for %s: %+v", accountID, err)
			}

			// Poll until properties are confirmed set
			if err = pollForTableProperties(ctx, client, props); err != nil {
				return fmt.Errorf("waiting for Table Properties to be set for %s: %+v", accountID, err)
			}

			metadata.SetID(accountID)
			return pluginsdk.SetResourceIdentityData(metadata.ResourceData, accountID, pluginsdk.ResourceTypeForIdentityVirtual)
		},
	}
}

func expandTablePropertiesModel(model *AccountTablePropertiesModel, props *tables.StorageServiceProperties) {
	if len(model.CorsRule) >= 1 {
		corsRules := make([]tables.CorsRule, 0)
		for _, corsRule := range model.CorsRule {
			corsRules = append(corsRules, tables.CorsRule{
				AllowedOrigins:  strings.Join(corsRule.AllowedOrigins, ","),
				AllowedMethods:  strings.Join(corsRule.AllowedMethods, ","),
				AllowedHeaders:  strings.Join(corsRule.AllowedHeaders, ","),
				ExposedHeaders:  strings.Join(corsRule.ExposedHeaders, ","),
				MaxAgeInSeconds: int(corsRule.MaxAgeSeconds),
			})
		}

		props.Cors.CorsRule = corsRules
	}

	if len(model.HourMetrics) == 1 {
		metrics := model.HourMetrics[0]
		props.HourMetrics.Enabled = true
		props.HourMetrics.Version = metrics.Version
		if metrics.RetentionPolicyDays != 0 {
			props.HourMetrics.RetentionPolicy = tables.RetentionPolicy{
				Days:    int(metrics.RetentionPolicyDays),
				Enabled: true,
			}
		}

		props.HourMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
	}

	if len(model.MinuteMetrics) != 0 {
		metrics := model.MinuteMetrics[0]
		props.MinuteMetrics.Enabled = true
		props.MinuteMetrics.Version = metrics.Version
		if metrics.RetentionPolicyDays != 0 {
			props.MinuteMetrics.RetentionPolicy = tables.RetentionPolicy{
				Days:    int(metrics.RetentionPolicyDays),
				Enabled: true,
			}
		}

		props.MinuteMetrics.IncludeAPIs = pointer.To(metrics.IncludeAPIS)
	}

	if len(model.Logging) != 0 {
		logging := model.Logging[0]
		props.Logging.Version = logging.Version
		props.Logging.Delete = logging.Delete
		props.Logging.Read = logging.Read
		props.Logging.Write = logging.Write
		if logging.RetentionPolicyDays != 0 {
			props.Logging.RetentionPolicy = tables.RetentionPolicy{
				Enabled: true,
				Days:    int(logging.RetentionPolicyDays),
			}
		}
	}
}
