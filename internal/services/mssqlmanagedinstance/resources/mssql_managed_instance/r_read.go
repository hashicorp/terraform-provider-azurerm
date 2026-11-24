// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mssql_managed_instance

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/maintenance/2023-04-01/publicmaintenanceconfigurations"
	"github.com/hashicorp/go-azure-sdk/resource-manager/sql/2025-01-01/managedinstances"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
)

func (r Resource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.MSSQLManagedInstance.ManagedInstancesClient

			id, err := commonids.ParseSqlManagedInstanceID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var state MsSqlManagedInstanceModel
			if err := metadata.Decode(&state); err != nil {
				return err
			}

			existing, err := client.Get(ctx, *id, managedinstances.GetOperationOptions{})
			if err != nil {
				if response.WasNotFound(existing.HttpResponse) {
					return metadata.MarkAsGone(id)
				}
				return fmt.Errorf("retrieving %s: %v", id, err)
			}
			model := MsSqlManagedInstanceModel{}

			if existing.Model != nil {
				model = MsSqlManagedInstanceModel{
					Name:              id.ManagedInstanceName,
					Location:          location.Normalize(existing.Model.Location),
					ResourceGroupName: id.ResourceGroupName,
					Identity:          r.flattenIdentity(existing.Model.Identity),
					Tags:              pointer.From(existing.Model.Tags),

					// This value is not returned, so we'll just set whatever is in the state/config
					AdministratorLoginPassword: state.AdministratorLoginPassword,
					// This value is not returned, so we'll just set whatever is in the state/config
					DnsZonePartnerId: state.DnsZonePartnerId,
				}

				if sku := existing.Model.Sku; sku != nil {
					model.SkuName = r.normalizeSku(sku.Name)
				}

				if props := existing.Model.Properties; props != nil {
					model.LicenseType = pointer.FromEnum(props.LicenseType)
					model.ProxyOverride = pointer.FromEnum(props.ProxyOverride)
					model.StorageAccountType = backupStorageRedundancyToStorageAccType(pointer.From(props.RequestedBackupStorageRedundancy))

					model.AdministratorLogin = pointer.From(props.AdministratorLogin)

					// read from state since when `azuread_authentication_only` is enabled via resource `azurerm_mssql_managed_instance_active_directory_administrator`,
					// the API returns the value of `AzureActiveDirectoryAdministrator` which causes diff.
					// TODO: revisit this (on migration to Framework?), currently this resource cannot determine drift on `azure_active_directory_administrator`
					model.AzureActiveDirectoryAdministrator = state.AzureActiveDirectoryAdministrator

					model.Collation = pointer.From(props.Collation)
					model.DnsZone = pointer.From(props.DnsZone)
					model.Fqdn = pointer.From(props.FullyQualifiedDomainName)

					if props.MaintenanceConfigurationId != nil {
						maintenanceConfigId, err := publicmaintenanceconfigurations.ParsePublicMaintenanceConfigurationIDInsensitively(*props.MaintenanceConfigurationId)
						if err != nil {
							return err
						}
						model.MaintenanceConfigurationName = maintenanceConfigId.PublicMaintenanceConfigurationName
					}

					model.MinimumTlsVersion = pointer.From(props.MinimalTlsVersion)
					model.PublicDataEndpointEnabled = pointer.From(props.PublicDataEndpointEnabled)
					model.GeneralPurposeV2Enabled = pointer.From(props.IsGeneralPurposeV2)
					model.StorageIOps = pointer.From(props.StorageIOps)
					model.StorageSizeInGb = pointer.From(props.StorageSizeInGB)
					model.SubnetId = pointer.From(props.SubnetId)
					model.TimezoneId = pointer.From(props.TimezoneId)
					model.VCores = pointer.From(props.VCores)
					model.ZoneRedundantEnabled = pointer.From(props.ZoneRedundant)

					model.ServicePrincipalType = ""
					if props.ServicePrincipal != nil {
						model.ServicePrincipalType = pointer.FromEnum(props.ServicePrincipal.Type)
					}
					model.DatabaseFormat = pointer.FromEnum(props.DatabaseFormat)
					model.HybridSecondaryUsage = pointer.FromEnum(props.HybridSecondaryUsage)
				}
			}

			return metadata.Encode(&model)
		},
	}
}

func (r Resource) flattenIdentity(input *identity.LegacySystemAndUserAssignedMap) []identity.SystemOrUserAssignedList {
	if input == nil {
		return []identity.SystemOrUserAssignedList{}
	}

	identityIds := make([]string, 0)
	for k := range input.IdentityIds {
		parsedId, err := commonids.ParseUserAssignedIdentityIDInsensitively(k)
		if err != nil {
			continue
		}
		identityIds = append(identityIds, parsedId.ID())
	}

	return []identity.SystemOrUserAssignedList{{
		Type:        input.Type,
		PrincipalId: input.PrincipalId,
		TenantId:    input.TenantId,
		IdentityIds: identityIds,
	}}
}

const (
	StorageAccountTypeGRS  = "GRS"
	StorageAccountTypeLRS  = "LRS"
	StorageAccountTypeZRS  = "ZRS"
	StorageAccountTypeGZRS = "GZRS"
)

type AzureActiveDirectoryAdministrator struct {
	LoginUserName                    string `tfschema:"login_username"`
	ObjectID                         string `tfschema:"object_id"`
	AzureADAuthenticationOnlyEnabled bool   `tfschema:"azuread_authentication_only_enabled"`
	TenantID                         string `tfschema:"tenant_id"`
	PrincipalType                    string `tfschema:"principal_type"`
}

func (r Resource) normalizeSku(sku string) string {
	switch sku {
	case "MIBC64G8IH", "BC_G8IH":
		return "BC_Gen8IH"
	case "MIBC64G8IM", "BC_G8IM":
		return "BC_Gen8IM"
	case "MIGP4G8IH", "GP_G8IH":
		return "GP_Gen8IH"
	case "MIGP4G8IM", "GP_G8IM":
		return "GP_Gen8IM"
	}

	return sku
}

// the StorageAccountType property has changed to RequestedBackupStorageRedundancy with 1-1 mapping of the following values:
// GRS -> Geo
// ZRS -> Zone
// LRS -> Local
// GZRS -> GeoZone
func storageAccTypeToBackupStorageRedundancy(storageAccountType string) managedinstances.BackupStorageRedundancy {
	switch storageAccountType {
	case StorageAccountTypeZRS:
		return managedinstances.BackupStorageRedundancyZone
	case StorageAccountTypeLRS:
		return managedinstances.BackupStorageRedundancyLocal
	case StorageAccountTypeGZRS:
		return managedinstances.BackupStorageRedundancyGeoZone
	}
	return managedinstances.BackupStorageRedundancyGeo
}

func backupStorageRedundancyToStorageAccType(backupStorageRedundancy managedinstances.BackupStorageRedundancy) string {
	switch backupStorageRedundancy {
	case managedinstances.BackupStorageRedundancyZone:
		return StorageAccountTypeZRS
	case managedinstances.BackupStorageRedundancyLocal:
		return StorageAccountTypeLRS
	case managedinstances.BackupStorageRedundancyGeoZone:
		return StorageAccountTypeGZRS
	}
	return StorageAccountTypeGRS
}
