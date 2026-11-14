// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandArmServerNetwork(d *pluginsdk.ResourceData) *servers.Network {
	network := servers.Network{}

	if v, ok := d.GetOk("delegated_subnet_id"); ok {
		network.DelegatedSubnetResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		network.PrivateDnsZoneArmResourceId = pointer.To(v.(string))
	}

	publicNetworkAccessEnabled := servers.ServerPublicNetworkAccessStateEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccessEnabled = servers.ServerPublicNetworkAccessStateDisabled
	}
	network.PublicNetworkAccess = pointer.To(publicNetworkAccessEnabled)

	return &network
}

func expandArmServerMaintenanceWindow(input []any) *servers.MaintenanceWindowForPatch {
	if len(input) == 0 {
		return &servers.MaintenanceWindowForPatch{
			CustomWindow: pointer.To(ServerMaintenanceWindowDisabled),
		}
	}
	v := input[0].(map[string]any)

	maintenanceWindow := servers.MaintenanceWindowForPatch{
		CustomWindow: pointer.To(ServerMaintenanceWindowEnabled),
		StartHour:    pointer.To(int64(v["start_hour"].(int))),
		StartMinute:  pointer.To(int64(v["start_minute"].(int))),
		DayOfWeek:    pointer.To(int64(v["day_of_week"].(int))),
	}

	return &maintenanceWindow
}

func expandArmServerStorage(d *pluginsdk.ResourceData) *servers.Storage {
	storage := servers.Storage{}

	autoGrow := servers.StorageAutoGrowDisabled
	if v, ok := d.GetOk("auto_grow_enabled"); ok && v.(bool) {
		autoGrow = servers.StorageAutoGrowEnabled
	}
	storage.AutoGrow = &autoGrow

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageSizeGB = pointer.To(int64(v.(int) / 1024))
	}

	storage.Type = pointer.ToEnum[servers.StorageType](d.Get("storage_type").(string))

	if v, ok := d.GetOk("storage_tier"); ok {
		storage.Tier = pointer.ToEnum[servers.AzureManagedDiskPerformanceTier](v.(string))
	}

	if v, ok := d.GetOk("storage_iops"); ok {
		storage.Iops = pointer.To(int64(v.(int)))
	}

	if v, ok := d.GetOk("storage_throughput"); ok {
		storage.Throughput = pointer.To(int64(v.(int)))
	}

	return &storage
}

func expandArmServerBackup(d *pluginsdk.ResourceData) *servers.Backup {
	backup := servers.Backup{}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		backup.BackupRetentionDays = pointer.To(int64(v.(int)))
	}

	geoRedundantEnabled := servers.GeographicallyRedundantBackupDisabled
	if geoRedundantBackupEnabled := d.Get("geo_redundant_backup_enabled").(bool); geoRedundantBackupEnabled {
		geoRedundantEnabled = servers.GeographicallyRedundantBackupEnabled
	}

	backup.GeoRedundantBackup = &geoRedundantEnabled

	return &backup
}

func expandArmServerBackupForPatch(d *pluginsdk.ResourceData) *servers.BackupForPatch {
	backup := servers.BackupForPatch{}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		backup.BackupRetentionDays = pointer.To(int64(v.(int)))
	}

	geoRedundantEnabled := servers.GeographicallyRedundantBackupDisabled
	if geoRedundantBackupEnabled := d.Get("geo_redundant_backup_enabled").(bool); geoRedundantBackupEnabled {
		geoRedundantEnabled = servers.GeographicallyRedundantBackupEnabled
	}

	backup.GeoRedundantBackup = &geoRedundantEnabled

	return &backup
}

func expandFlexibleServerSku(name string) (*servers.Sku, error) {
	if name == "" {
		return nil, nil
	}
	parts := strings.SplitAfterN(name, "_", 2)

	var tier servers.SkuTier
	switch strings.TrimSuffix(parts[0], "_") {
	case "B":
		tier = servers.SkuTierBurstable
	case "GP":
		tier = servers.SkuTierGeneralPurpose
	case "MO":
		tier = servers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", name, parts[0])
	}

	return &servers.Sku{
		Name: parts[1],
		Tier: tier,
	}, nil
}

func expandFlexibleServerSkuForPatch(name string) (*servers.SkuForPatch, error) {
	if name == "" {
		return nil, nil
	}
	parts := strings.SplitAfterN(name, "_", 2)

	var tier servers.SkuTier
	switch strings.TrimSuffix(parts[0], "_") {
	case "B":
		tier = servers.SkuTierBurstable
	case "GP":
		tier = servers.SkuTierGeneralPurpose
	case "MO":
		tier = servers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", name, parts[0])
	}

	return &servers.SkuForPatch{
		Name: &parts[1],
		Tier: &tier,
	}, nil
}

func expandFlexibleServerHighAvailability(inputs []any, isCreate bool) *servers.HighAvailability {
	if len(inputs) == 0 || inputs[0] == nil {
		return &servers.HighAvailability{
			Mode: pointer.To(servers.PostgreSqlFlexibleServerHighAvailabilityModeDisabled),
		}
	}

	input := inputs[0].(map[string]any)

	result := servers.HighAvailability{
		Mode: pointer.ToEnum[servers.PostgreSqlFlexibleServerHighAvailabilityMode](input["mode"].(string)),
	}

	// service team confirmed it doesn't support to update `high_availability.0.standby_availability_zone` after the PostgreSQL Flexible Server resource is created
	if isCreate {
		if v, ok := input["standby_availability_zone"]; ok && v.(string) != "" {
			result.StandbyAvailabilityZone = pointer.To(v.(string))
		}
	}

	return &result
}

func expandFlexibleServerHighAvailabilityForPatch(inputs []any, isCreate bool) *servers.HighAvailabilityForPatch {
	if len(inputs) == 0 || inputs[0] == nil {
		return &servers.HighAvailabilityForPatch{
			Mode: pointer.To(servers.PostgreSqlFlexibleServerHighAvailabilityModeDisabled),
		}
	}

	input := inputs[0].(map[string]any)

	result := servers.HighAvailabilityForPatch{
		Mode: pointer.ToEnum[servers.PostgreSqlFlexibleServerHighAvailabilityMode](input["mode"].(string)),
	}

	// service team confirmed it doesn't support to update `high_availability.0.standby_availability_zone` after the PostgreSQL Flexible Server resource is created
	if isCreate {
		if v, ok := input["standby_availability_zone"]; ok && v.(string) != "" {
			result.StandbyAvailabilityZone = pointer.To(v.(string))
		}
	}

	return &result
}

func expandFlexibleServerAuthConfig(authRaw []any) *servers.AuthConfig {
	if len(authRaw) == 0 || authRaw[0] == nil {
		return nil
	}

	authConfigs := authRaw[0].(map[string]any)
	out := servers.AuthConfig{}

	activeDirectoryAuthEnabled := servers.MicrosoftEntraAuthDisabled
	if authConfigs["active_directory_auth_enabled"].(bool) {
		activeDirectoryAuthEnabled = servers.MicrosoftEntraAuthEnabled
	}
	out.ActiveDirectoryAuth = &activeDirectoryAuthEnabled

	passwordAuthEnabled := servers.PasswordBasedAuthDisabled
	if authConfigs["password_auth_enabled"].(bool) {
		passwordAuthEnabled = servers.PasswordBasedAuthEnabled
	}
	out.PasswordAuth = &passwordAuthEnabled

	if tenantId, ok := authConfigs["tenant_id"].(string); ok {
		out.TenantId = &tenantId
	}

	return &out
}

func expandFlexibleServerAuthConfigForPatch(authRaw []any) *servers.AuthConfigForPatch {
	if len(authRaw) == 0 || authRaw[0] == nil {
		return nil
	}

	authConfigs := authRaw[0].(map[string]any)
	out := servers.AuthConfigForPatch{}

	activeDirectoryAuthEnabled := servers.MicrosoftEntraAuthDisabled
	if authConfigs["active_directory_auth_enabled"].(bool) {
		activeDirectoryAuthEnabled = servers.MicrosoftEntraAuthEnabled
	}
	out.ActiveDirectoryAuth = &activeDirectoryAuthEnabled

	passwordAuthEnabled := servers.PasswordBasedAuthDisabled
	if authConfigs["password_auth_enabled"].(bool) {
		passwordAuthEnabled = servers.PasswordBasedAuthEnabled
	}
	out.PasswordAuth = &passwordAuthEnabled

	if tenantId, ok := authConfigs["tenant_id"].(string); ok {
		out.TenantId = &tenantId
	}

	return &out
}

func expandFlexibleServerDataEncryption(input []any) *servers.DataEncryption {
	if len(input) == 0 {
		return nil
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

func expandFlexibleServerCluster(input []any) *servers.Cluster {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	v := input[0].(map[string]any)

	cluster := servers.Cluster{}

	if size := v["size"].(int); size != 0 {
		cluster.ClusterSize = pointer.To(int64(size))
	}

	if databaseName := v["default_database_name"].(string); databaseName != "" {
		cluster.DefaultDatabaseName = pointer.To(databaseName)
	}

	return &cluster
}
