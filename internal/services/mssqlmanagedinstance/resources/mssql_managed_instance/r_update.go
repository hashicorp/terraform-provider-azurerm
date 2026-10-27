// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceazureadonlyauthentications"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			adminClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAdministratorsClient
			azureADAuthenticationOnlyClient := metadata.Client.MSSQLManagedInstance.ManagedInstanceAzureADOnlyAuthenticationsClient

			id, err := commonids.ParseSqlManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MsSqlManagedInstanceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id, managedinstances.DefaultGetOperationOptions())
			if err != nil {
				return fmt.Errorf("retrieving %s: %v", id, err)
			}

			if existing.Model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", id)
			}

			if existing.Model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", id)
			}
			props := existing.Model.Properties

			// `Administrators` is only valid when specified during creation
			// This not ideal, but matches previous behaviour (when the request body was built from scratch rather than modifying a returned one)
			// Without this, we receive this error: `Invalid value given for parameter AzureADOnlyAuthentication`
			props.Administrators = nil

			// The API does not allow `storageThroughputMBps` in the PUT request
			props.StorageThroughputMBps = nil

			if metadata.ResourceData.HasChange("sku_name") {
				sku, err := r.expandSkuName(state.SkuName)
				if err != nil {
					return fmt.Errorf("expanding `sku_name` for SQL Managed Instance Server %q: %v", *id, err)
				}
				existing.Model.Sku = sku
			}

			if metadata.ResourceData.HasChange("license_type") {
				props.LicenseType = pointer.ToEnum[managedinstances.ManagedInstanceLicenseType](state.LicenseType)
			}

			if metadata.ResourceData.HasChange("storage_size_in_gb") {
				props.StorageSizeInGB = pointer.To(state.StorageSizeInGb)
			}

			if metadata.ResourceData.HasChange("subnet_id") {
				props.SubnetId = pointer.To(state.SubnetId)
			}

			if metadata.ResourceData.HasChange("vcores") {
				props.VCores = pointer.To(state.VCores)
			}

			if metadata.ResourceData.HasChange("administrator_login_password") {
				props.AdministratorLoginPassword = pointer.To(state.AdministratorLoginPassword)
			}

			if metadata.ResourceData.HasChange("identity") {
				existing.Model.Identity = r.expandIdentity(state.Identity)

				// @sreallymatt: This is pre-existing logic that is technically broken, but no one seems to be complaining about it ¯\_(ツ)_/¯
				// TODO: revisit, we'll likely want to introduce a `primary_user_assigned_identity_id` property
				// rather than grabbing the ID from the first loop iteration given maps are unordered...
				if existing.Model.Identity != nil && len(existing.Model.Identity.IdentityIds) > 0 {
					for k := range existing.Model.Identity.IdentityIds {
						props.PrimaryUserAssignedIdentityId = pointer.To(k)
						break
					}
				}
			}

			if metadata.ResourceData.HasChange("minimum_tls_version") {
				props.MinimalTlsVersion = pointer.To(state.MinimumTlsVersion)
			}

			if metadata.ResourceData.HasChange("proxy_override") {
				props.ProxyOverride = pointer.ToEnum[managedinstances.ManagedInstanceProxyOverride](state.ProxyOverride)
			}

			if metadata.ResourceData.HasChange("public_data_endpoint_enabled") {
				props.PublicDataEndpointEnabled = pointer.To(state.PublicDataEndpointEnabled)
			}

			if metadata.ResourceData.HasChange("storage_account_type") {
				props.RequestedBackupStorageRedundancy = pointer.To(storageAccTypeToBackupStorageRedundancy(state.StorageAccountType))
			}

			if metadata.ResourceData.HasChange("zone_redundant_enabled") {
				props.ZoneRedundant = pointer.To(state.ZoneRedundantEnabled)
			}

			if metadata.ResourceData.HasChange("dns_zone_partner_id") {
				props.DnsZonePartner = pointer.To(state.DnsZonePartnerId)
			}

			if metadata.ResourceData.HasChange("tags") {
				existing.Model.Tags = pointer.To(state.Tags)
			}

			if metadata.ResourceData.HasChange("maintenance_configuration_name") {
				maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(id.SubscriptionId, state.MaintenanceConfigurationName)
				props.MaintenanceConfigurationId = pointer.To(maintenanceConfigId.ID())
			}

			if metadata.ResourceData.HasChange("service_principal_type") {
				props.ServicePrincipal = &managedinstances.ServicePrincipal{}
				if state.ServicePrincipalType == "" {
					props.ServicePrincipal.Type = pointer.To(managedinstances.ServicePrincipalTypeNone)
				} else {
					props.ServicePrincipal.Type = pointer.ToEnum[managedinstances.ServicePrincipalType](state.ServicePrincipalType)
				}
			}

			if metadata.ResourceData.HasChange("azure_active_directory_administrator") {
				// Need to check if Microsoft AAD authentication only is enabled or not before calling delete, else you will get the following error:
				// InvalidManagedServerAADOnlyAuthNoAADAdminPropertyName: AAD Admin is not configured,
				// AAD Admin must be set before enabling/disabling AAD Authentication Only.
				log.Printf("[INFO] Checking if AAD Administrator exists")
				aadAdminExists := false
				resp, err := adminClient.Get(ctx, *id)
				if err != nil {
					if !response.WasNotFound(resp.HttpResponse) {
						return fmt.Errorf("retrieving the Administrators of %s: %+v", *id, err)
					}
				} else {
					aadAdminExists = true
				}

				if aadAdminExists {
					// Before deleting an AAD admin, it is necessary to disable `AzureADOnlyAuthentication` first, as deleting an AAD admin when `AzureADOnlyAuthentication` feature is enabled is not supported.
					// Use `CreateOrUpdateThenPoll` instead of `DeleteThenPoll`, because the actual deletion behaviour of the API is not to really delete the record, but to update `AzureADOnlyAuthentication` to false. Therefore, using `DeleteThenPoll` will cause pull till done to never end until it times out.
					aadAuthOnlyParams := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
						Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{},
					}
					if err = azureADAuthenticationOnlyClient.CreateOrUpdateThenPoll(ctx, *id, aadAuthOnlyParams); err != nil {
						return fmt.Errorf("disabling `azuread_authentication_only` for %s: %+v", *id, err)
					}

					if err := adminClient.DeleteThenPoll(ctx, *id); err != nil {
						return fmt.Errorf("removing the AAD Administrator for %s: %+v", *id, err)
					}
				}

				aadAdminProps := expandMsSqlManagedInstanceAdministrators(state.AzureActiveDirectoryAdministrator)
				if aadAdminProps != nil {
					if err := adminClient.CreateOrUpdateThenPoll(ctx, *id, *aadAdminProps); err != nil {
						return fmt.Errorf("creating AAD Administrator for %s: %+v", *id, err)
					}
				}

				if aadOnlyAuthenticationsEnabled := expandMsSqlManagedInstanceAadAuthenticationOnly(state.AzureActiveDirectoryAdministrator); aadOnlyAuthenticationsEnabled {
					aadOnlyAuthenticationsProps := managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthentication{
						Properties: &managedinstanceazureadonlyauthentications.ManagedInstanceAzureADOnlyAuthProperties{
							AzureADOnlyAuthentication: true,
						},
					}

					if err := azureADAuthenticationOnlyClient.CreateOrUpdateThenPoll(ctx, *id, aadOnlyAuthenticationsProps); err != nil {
						return fmt.Errorf("setting `azuread_authentication_only_enabled` for %s: %+v", *id, err)
					}
				}
			}

			if metadata.ResourceData.HasChange("database_format") {
				props.DatabaseFormat = pointer.ToEnum[managedinstances.ManagedInstanceDatabaseFormat](state.DatabaseFormat)
			}

			if metadata.ResourceData.HasChange("hybrid_secondary_usage") {
				props.HybridSecondaryUsage = pointer.ToEnum[managedinstances.HybridSecondaryUsage](state.HybridSecondaryUsage)
			}

			effectiveIsGeneralPurposeV2 := expandMsSqlManagedInstanceGeneralPurposeV2Enabled(state.GeneralPurposeV2Enabled, state.SkuName)

			if metadata.ResourceData.HasChanges("general_purpose_v2_enabled", "sku_name") {
				props.IsGeneralPurposeV2 = effectiveIsGeneralPurposeV2
			}

			// When updating MI from GPv2 to non GPv2, `StorageIOps` is not a valid value in payload.
			if pointer.From(effectiveIsGeneralPurposeV2) && state.StorageIOps != 0 {
				if metadata.ResourceData.HasChange("storage_iops") {
					props.StorageIOps = pointer.To(state.StorageIOps)
				}
			} else {
				props.StorageIOps = nil
			}

			if err := client.CreateOrUpdateThenPoll(ctx, *id, *existing.Model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}
