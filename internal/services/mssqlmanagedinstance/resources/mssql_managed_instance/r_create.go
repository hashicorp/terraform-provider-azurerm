// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstanceadministrators"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 24 * time.Hour,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model MsSqlManagedInstanceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := commonids.NewSqlManagedInstanceID(subscriptionId, model.ResourceGroupName, model.Name)

			if !metadata.Client.Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
				existing, err := client.Get(ctx, id, managedinstances.GetOperationOptions{})
				if err != nil && !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}

				if !response.WasNotFound(existing.HttpResponse) {
					return metadata.ResourceRequiresImport(r.ResourceType(), id)
				}
			}

			sku, err := r.expandSkuName(model.SkuName)
			if err != nil {
				return fmt.Errorf("expanding `sku_name` for SQL Managed Instance Server %q: %v", id.ID(), err)
			}

			maintenanceConfigId := publicmaintenanceconfigurations.NewPublicMaintenanceConfigurationID(subscriptionId, model.MaintenanceConfigurationName)

			isGeneralPurposeV2 := expandMsSqlManagedInstanceGeneralPurposeV2Enabled(model.GeneralPurposeV2Enabled, model.SkuName)

			parameters := managedinstances.ManagedInstance{
				Sku:      sku,
				Identity: r.expandIdentity(model.Identity),
				Location: location.Normalize(model.Location),
				Properties: &managedinstances.ManagedInstanceProperties{
					AdministratorLogin:               pointer.To(model.AdministratorLogin),
					AdministratorLoginPassword:       pointer.To(model.AdministratorLoginPassword),
					Collation:                        pointer.To(model.Collation),
					DnsZonePartner:                   pointer.To(model.DnsZonePartnerId),
					IsGeneralPurposeV2:               isGeneralPurposeV2,
					LicenseType:                      pointer.ToEnum[managedinstances.ManagedInstanceLicenseType](model.LicenseType),
					MaintenanceConfigurationId:       pointer.To(maintenanceConfigId.ID()),
					MinimalTlsVersion:                pointer.To(model.MinimumTlsVersion),
					ProxyOverride:                    pointer.ToEnum[managedinstances.ManagedInstanceProxyOverride](model.ProxyOverride),
					PublicDataEndpointEnabled:        pointer.To(model.PublicDataEndpointEnabled),
					RequestedBackupStorageRedundancy: pointer.To(storageAccTypeToBackupStorageRedundancy(model.StorageAccountType)),
					StorageSizeInGB:                  pointer.To(model.StorageSizeInGb),
					SubnetId:                         pointer.To(model.SubnetId),
					TimezoneId:                       pointer.To(model.TimezoneId),
					VCores:                           pointer.To(model.VCores),
					ZoneRedundant:                    pointer.To(model.ZoneRedundantEnabled),
					// `Administrators` is only valid when specified during creation`
					Administrators:       expandMsSqlManagedInstanceExternalAdministrators(model.AzureActiveDirectoryAdministrator),
					DatabaseFormat:       pointer.ToEnum[managedinstances.ManagedInstanceDatabaseFormat](model.DatabaseFormat),
					HybridSecondaryUsage: pointer.ToEnum[managedinstances.HybridSecondaryUsage](model.HybridSecondaryUsage),
				},
				Tags: pointer.To(model.Tags),
			}

			// If StorageIOps is carried in payload and set to 0, the service will return HTTP 405.
			if pointer.From(isGeneralPurposeV2) && model.StorageIOps != 0 {
				parameters.Properties.StorageIOps = pointer.To(model.StorageIOps)
			}

			if parameters.Identity != nil && len(parameters.Identity.IdentityIds) > 0 {
				for k := range parameters.Identity.IdentityIds {
					parameters.Properties.PrimaryUserAssignedIdentityId = pointer.To(k)
					break
				}
			}

			if model.ServicePrincipalType != "" {
				parameters.Properties.ServicePrincipal = &managedinstances.ServicePrincipal{
					Type: pointer.ToEnum[managedinstances.ServicePrincipalType](model.ServicePrincipalType),
				}
			}

			if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, metadata.SetIDCallback(&id)); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r Resource) expandIdentity(input []identity.SystemOrUserAssignedList) *identity.LegacySystemAndUserAssignedMap {
	if len(input) == 0 {
		return nil
	}

	// Workaround for issue https://github.com/Azure/azure-rest-api-specs/issues/16838
	if input[0].Type == identity.TypeNone {
		return nil
	}

	var identityIds map[string]identity.UserAssignedIdentityDetails
	if len(input[0].IdentityIds) != 0 {
		identityIds = map[string]identity.UserAssignedIdentityDetails{}
		for _, id := range input[0].IdentityIds {
			identityIds[id] = identity.UserAssignedIdentityDetails{}
		}
	}

	return &identity.LegacySystemAndUserAssignedMap{
		Type:        input[0].Type,
		IdentityIds: identityIds,
	}
}

func (r Resource) expandSkuName(skuName string) (*managedinstances.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 2 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier string
	switch parts[0] {
	case "GP":
		tier = "GeneralPurpose"
	case "BC":
		tier = "BusinessCritical"
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	return &managedinstances.Sku{
		Name:   skuName,
		Tier:   pointer.To(tier),
		Family: pointer.To(parts[1]),
	}, nil
}

func expandMsSqlManagedInstanceAadAuthenticationOnly(input []AzureActiveDirectoryAdministrator) bool {
	if len(input) == 0 {
		return false
	}

	if ok := input[0].AzureADAuthenticationOnlyEnabled; ok {
		return input[0].AzureADAuthenticationOnlyEnabled
	}

	return false
}

func expandMsSqlManagedInstanceExternalAdministrators(input []AzureActiveDirectoryAdministrator) *managedinstances.ManagedInstanceExternalAdministrator {
	if len(input) == 0 {
		return nil
	}

	admin := input[0]
	adminParams := managedinstances.ManagedInstanceExternalAdministrator{
		AdministratorType:         pointer.To(managedinstances.AdministratorTypeActiveDirectory),
		PrincipalType:             pointer.ToEnum[managedinstances.PrincipalType](admin.PrincipalType),
		Login:                     pointer.To(admin.LoginUserName),
		Sid:                       pointer.To(admin.ObjectID),
		AzureADOnlyAuthentication: pointer.To(admin.AzureADAuthenticationOnlyEnabled),
	}

	if admin.TenantID != "" {
		adminParams.TenantId = pointer.To(admin.TenantID)
	}

	return &adminParams
}

func expandMsSqlManagedInstanceAdministrators(input []AzureActiveDirectoryAdministrator) *managedinstanceadministrators.ManagedInstanceAdministrator {
	if len(input) == 0 {
		return nil
	}

	admin := input[0]

	adminProps := managedinstanceadministrators.ManagedInstanceAdministrator{
		Properties: &managedinstanceadministrators.ManagedInstanceAdministratorProperties{
			AdministratorType: managedinstanceadministrators.ManagedInstanceAdministratorTypeActiveDirectory,
			Login:             admin.LoginUserName,
			Sid:               admin.ObjectID,
		},
	}

	if admin.TenantID != "" {
		adminProps.Properties.TenantId = pointer.To(admin.TenantID)
	}

	return pointer.To(adminProps)
}

func expandMsSqlManagedInstanceGeneralPurposeV2Enabled(input bool, sku string) *bool {
	// The API returns an error when `isGeneralPurposeV2` is set on a non-general purpose SKU regardless of the value
	if strings.HasPrefix(sku, "BC_") {
		return nil
	}
	return &input
}
