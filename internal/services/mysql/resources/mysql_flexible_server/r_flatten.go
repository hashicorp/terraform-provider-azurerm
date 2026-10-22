// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/keyvault"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func resourceMysqlFlexibleServerFlatten(d *pluginsdk.ResourceData, id *servers.FlexibleServerId, server *servers.Server) error {
	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if server != nil {
		d.Set("location", location.NormalizeNilable(&server.Location))
		if props := server.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("zone", props.AvailabilityZone)
			d.Set("version", pointer.FromEnum(props.Version))
			d.Set("fqdn", props.FullyQualifiedDomainName)
			d.Set("source_server_id", props.SourceServerResourceId)

			if network := props.Network; network != nil {
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				d.Set("private_dns_zone_id", network.PrivateDnsZoneResourceId)
				d.Set("public_network_access", pointer.FromEnum(network.PublicNetworkAccess))
			}

			cmk, err := flattenFlexibleServerDataEncryption(props.DataEncryption)
			if err != nil {
				return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
			}
			if err := d.Set("customer_managed_key", cmk); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			identity, err := flattenFlexibleServerIdentity(server.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}

			if err := d.Set("maintenance_window", flattenArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
				return fmt.Errorf("setting `maintenance_window`: %+v", err)
			}

			if err := d.Set("storage", flattenArmServerStorage(props.Storage)); err != nil {
				return fmt.Errorf("setting `storage`: %+v", err)
			}

			if backup := props.Backup; backup != nil {
				d.Set("backup_retention_days", backup.BackupRetentionDays)
				d.Set("geo_redundant_backup_enabled", *backup.GeoRedundantBackup == servers.EnableStatusEnumEnabled)
			}

			if err := d.Set("high_availability", flattenFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
				return fmt.Errorf("setting `high_availability`: %+v", err)
			}
			d.Set("replication_role", pointer.FromEnum(props.ReplicationRole))
			d.Set("replica_capacity", props.ReplicaCapacity)
		}
		sku, err := flattenFlexibleServerSku(server.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name`: %+v", err)
		}
		d.Set("sku_name", sku)

		d.Set("administrator_password_wo_version", d.Get("administrator_password_wo_version").(int))

		if err := tags.FlattenAndSet(d, server.Tags); err != nil {
			return err
		}
	}

	return pluginsdk.SetResourceIdentityData(d, id)
}

func flattenArmServerStorage(storage *servers.Storage) []any {
	if storage == nil {
		return []any{}
	}

	var size, iops int64
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []any{
		map[string]any{
			"size_gb":             size,
			"iops":                iops,
			"auto_grow_enabled":   *storage.AutoGrow == servers.EnableStatusEnumEnabled,
			"io_scaling_enabled":  *storage.AutoIoScaling == servers.EnableStatusEnumEnabled,
			"log_on_disk_enabled": *storage.LogOnDisk == servers.EnableStatusEnumEnabled,
		},
	}
}

func flattenFlexibleServerSku(sku *servers.MySQLServerSku) (string, error) {
	if sku == nil || sku.Name == "" || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case servers.ServerSkuTierBurstable:
		tier = "B"
	case servers.ServerSkuTierGeneralPurpose:
		tier = "GP"
	case servers.ServerSkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, sku.Name}, "_"), nil
}

func flattenArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []any {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]any, 0)
	}

	return []any{
		map[string]any{
			"day_of_week":  pointer.From(input.DayOfWeek),
			"start_hour":   pointer.From(input.StartHour),
			"start_minute": pointer.From(input.StartMinute),
		},
	}
}

func flattenFlexibleServerHighAvailability(ha *servers.HighAvailability) []any {
	if ha == nil || *ha.Mode == servers.HighAvailabilityModeDisabled {
		return []any{}
	}

	return []any{
		map[string]any{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": pointer.From(ha.StandbyAvailabilityZone),
		},
	}
}

func flattenFlexibleServerDataEncryption(de *servers.DataEncryption) ([]any, error) {
	if de == nil || *de.Type == servers.DataEncryptionTypeSystemManaged {
		return []any{}, nil
	}

	item := map[string]any{}
	if de.PrimaryKeyURI != nil {
		keyID, err := keyvault.ParseNestedItemID(*de.PrimaryKeyURI, keyvault.VersionTypeAny, keyvault.NestedItemTypeKey)
		if err != nil {
			return nil, err
		}

		item["key_vault_key_id"] = keyID.ID()
	}

	if identity := de.PrimaryUserAssignedIdentityId; identity != nil {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*identity)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", *identity, err)
		}
		item["primary_user_assigned_identity_id"] = parsed.ID()
	}

	if de.GeoBackupKeyURI != nil {
		item["geo_backup_key_vault_key_id"] = *de.GeoBackupKeyURI
	}
	if identity := de.GeoBackupUserAssignedIdentityId; identity != nil {
		parsed, err := commonids.ParseUserAssignedIdentityIDInsensitively(*identity)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", *identity, err)
		}
		item["geo_backup_user_assigned_identity_id"] = parsed.ID()
	}

	return []any{item}, nil
}

func flattenFlexibleServerIdentity(input *servers.MySQLServerIdentity) (*[]any, error) {
	var transform *identity.UserAssignedMap

	if input != nil {
		transform = &identity.UserAssignedMap{
			Type:        identity.Type(string(*input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		for k := range *input.UserAssignedIdentities {
			transform.IdentityIds[k] = identity.UserAssignedIdentityDetails{}
		}
	}

	return identity.FlattenUserAssignedMap(transform)
}
