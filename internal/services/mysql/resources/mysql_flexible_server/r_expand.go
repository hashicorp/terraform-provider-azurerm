// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandArmServerNetwork(d *pluginsdk.ResourceData) *servers.Network {
	network := servers.Network{}

	if v, ok := d.GetOk("delegated_subnet_id"); ok {
		network.DelegatedSubnetResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		network.PrivateDnsZoneResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("public_network_access"); ok {
		network.PublicNetworkAccess = pointer.ToEnum[servers.EnableStatusEnum](v.(string))
	}

	return &network
}

func expandArmServerMaintenanceWindow(input []any) *servers.MaintenanceWindow {
	if len(input) == 0 {
		return &servers.MaintenanceWindow{
			CustomWindow: pointer.To(ServerMaintenanceWindowDisabled),
		}
	}
	v := input[0].(map[string]any)

	maintenanceWindow := servers.MaintenanceWindow{
		CustomWindow: pointer.To(ServerMaintenanceWindowEnabled),
		StartHour:    pointer.To(int64(v["start_hour"].(int))),
		StartMinute:  pointer.To(int64(v["start_minute"].(int))),
		DayOfWeek:    pointer.To(int64(v["day_of_week"].(int))),
	}

	return &maintenanceWindow
}

func expandArmServerStorage(inputs []any) *servers.Storage {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]any)
	autoGrow := servers.EnableStatusEnumDisabled
	if v := input["auto_grow_enabled"].(bool); v {
		autoGrow = servers.EnableStatusEnumEnabled
	}

	autoIoScaling := servers.EnableStatusEnumDisabled
	if v := input["io_scaling_enabled"].(bool); v {
		autoIoScaling = servers.EnableStatusEnumEnabled
	}

	logOnDisk := servers.EnableStatusEnumDisabled
	if v := input["log_on_disk_enabled"].(bool); v {
		logOnDisk = servers.EnableStatusEnumEnabled
	}

	storage := servers.Storage{
		AutoGrow:      &autoGrow,
		AutoIoScaling: &autoIoScaling,
		LogOnDisk:     &logOnDisk,
	}

	if v := input["size_gb"].(int); v != 0 {
		storage.StorageSizeGB = pointer.To(int64(v))
	}

	if v := input["iops"].(int); v != 0 {
		storage.Iops = pointer.To(int64(v))
	}

	return &storage
}

func expandArmServerBackup(d *pluginsdk.ResourceData) *servers.Backup {
	geoRedundantBackup := servers.EnableStatusEnumDisabled
	if d.Get("geo_redundant_backup_enabled").(bool) {
		geoRedundantBackup = servers.EnableStatusEnumEnabled
	}

	backup := servers.Backup{
		GeoRedundantBackup: &geoRedundantBackup,
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		backup.BackupRetentionDays = pointer.To(int64(v.(int)))
	}

	return &backup
}

func expandFlexibleServerSku(name string) (*servers.MySQLServerSku, error) {
	if name == "" {
		return nil, nil
	}
	parts := strings.SplitAfterN(name, "_", 2)

	var tier servers.ServerSkuTier
	switch strings.TrimSuffix(parts[0], "_") {
	case "B":
		tier = servers.ServerSkuTierBurstable
	case "GP":
		tier = servers.ServerSkuTierGeneralPurpose
	case "MO":
		tier = servers.ServerSkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", name, parts[0])
	}

	return &servers.MySQLServerSku{
		Name: parts[1],
		Tier: tier,
	}, nil
}

func expandFlexibleServerHighAvailability(inputs []any) *servers.HighAvailability {
	if len(inputs) == 0 || inputs[0] == nil {
		return &servers.HighAvailability{
			Mode: pointer.To(servers.HighAvailabilityModeDisabled),
		}
	}

	input := inputs[0].(map[string]any)

	mode := servers.HighAvailabilityMode(input["mode"].(string))

	result := servers.HighAvailability{
		Mode: &mode,
	}

	// for updating mode from ZoneRedundant to SameZone, the standby az will be changed
	// if we keep setting the standby az of ZoneRedundant, ha could not be changed to SameZone
	if mode == servers.HighAvailabilityModeSameZone {
		return &result
	}

	if v, ok := input["standby_availability_zone"]; ok && v.(string) != "" {
		result.StandbyAvailabilityZone = pointer.To(v.(string))
	}

	return &result
}

func expandFlexibleServerDataEncryption(input []any) *servers.DataEncryption {
	if len(input) == 0 || input[0] == nil {
		return &servers.DataEncryption{
			Type: pointer.To(servers.DataEncryptionTypeSystemManaged),
		}
	}
	v := input[0].(map[string]any)

	dataEncryption := servers.DataEncryption{
		Type: pointer.To(servers.DataEncryptionTypeAzureKeyVault),
	}

	if keyVaultKeyId := v["key_vault_key_id"].(string); keyVaultKeyId != "" {
		dataEncryption.PrimaryKeyURI = pointer.To(keyVaultKeyId)
	}

	if primaryUserAssignedIdentityId := v["primary_user_assigned_identity_id"].(string); primaryUserAssignedIdentityId != "" {
		dataEncryption.PrimaryUserAssignedIdentityId = pointer.To(primaryUserAssignedIdentityId)
	}

	if geoBackupKeyVaultKeyId := v["geo_backup_key_vault_key_id"].(string); geoBackupKeyVaultKeyId != "" {
		dataEncryption.GeoBackupKeyURI = pointer.To(geoBackupKeyVaultKeyId)
	}

	if geoBackupUserAssignedIdentityId := v["geo_backup_user_assigned_identity_id"].(string); geoBackupUserAssignedIdentityId != "" {
		dataEncryption.GeoBackupUserAssignedIdentityId = pointer.To(geoBackupUserAssignedIdentityId)
	}

	return &dataEncryption
}

func expandFlexibleServerIdentity(input []any) (*servers.MySQLServerIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	out := servers.MySQLServerIdentity{
		Type: pointer.ToEnum[servers.ManagedServiceIdentityType](string(expanded.Type)),
	}
	if expanded.Type == identity.TypeUserAssigned {
		ids := make(map[string]any)
		for k := range expanded.IdentityIds {
			ids[k] = struct{}{}
		}
		out.UserAssignedIdentities = &ids
	}

	return &out, nil
}
