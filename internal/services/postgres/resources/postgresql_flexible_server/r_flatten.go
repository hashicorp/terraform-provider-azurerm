// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
)

func flattenFlexibleServerSku(sku *servers.Sku) (string, error) {
	if sku == nil || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case servers.SkuTierBurstable:
		tier = "B"
	case servers.SkuTierGeneralPurpose:
		tier = "GP"
	case servers.SkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, sku.Name}, "_"), nil
}

func flattenArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []any {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == ServerMaintenanceWindowDisabled {
		return make([]any, 0)
	}

	dayOfWeek := pointer.From(input.DayOfWeek)
	startHour := pointer.From(input.StartHour)
	startMinute := pointer.From(input.StartMinute)
	return []any{
		map[string]any{
			"day_of_week":  dayOfWeek,
			"start_hour":   startHour,
			"start_minute": startMinute,
		},
	}
}

func flattenFlexibleServerHighAvailability(ha *servers.HighAvailability) []any {
	if ha == nil || ha.Mode == nil || *ha.Mode == servers.PostgreSqlFlexibleServerHighAvailabilityModeDisabled {
		return []any{}
	}

	zone := pointer.From(ha.StandbyAvailabilityZone)

	return []any{
		map[string]any{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}

func flattenFlexibleServerAuthConfig(ac *servers.AuthConfig) any {
	out := make(map[string]any, 0)

	if ac == nil {
		return out
	}

	aadEnabled := false
	if ac.ActiveDirectoryAuth != nil {
		aadEnabled = *ac.ActiveDirectoryAuth == servers.MicrosoftEntraAuthEnabled
	}
	out["active_directory_auth_enabled"] = aadEnabled

	// It is by design if PasswordAuthEnabled is not returned or undefined, we consider it as true.
	pwdEnabled := true
	if ac.PasswordAuth != nil {
		pwdEnabled = *ac.PasswordAuth == servers.PasswordBasedAuthEnabled
	}
	out["password_auth_enabled"] = pwdEnabled

	if ac.TenantId != nil {
		out["tenant_id"] = *ac.TenantId
	}

	result := make([]any, 0)
	result = append(result, out)
	return result
}

func flattenFlexibleServerDataEncryption(de *servers.DataEncryption) ([]any, error) {
	if de == nil || *de.Type != servers.DataEncryptionTypeAzureKeyVault {
		return []any{}, nil
	}

	item := map[string]any{}
	if de.PrimaryKeyURI != nil {
		item["key_vault_key_id"] = *de.PrimaryKeyURI
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

func flattenFlexibleServerCluster(cluster *servers.Cluster) []any {
	if cluster == nil {
		return []any{}
	}

	item := map[string]any{
		"size":                  pointer.From(cluster.ClusterSize),
		"default_database_name": pointer.From(cluster.DefaultDatabaseName),
	}

	return []any{item}
}
