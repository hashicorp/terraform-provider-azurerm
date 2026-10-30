// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	parameters := servers.ServerForPatch{
		Properties: &servers.ServerPropertiesForPatch{},
	}

	requireUpdateOnLogin := false // it's required to call Create with `createMode` set to `Update` to update login name.

	woPassword, err := pluginsdk.GetWriteOnly(d, "administrator_password_wo", cty.String)
	if err != nil {
		return err
	}

	createMode := d.Get("create_mode").(string)
	if createMode == "" || servers.CreateMode(createMode) == servers.CreateModeDefault {
		_, adminLoginSet := d.GetOk("administrator_login")
		_, adminPwdSet := d.GetOk("administrator_password")

		pwdEnabled := true // it defaults to true
		if authRaw, authExist := d.GetOk("authentication"); authExist {
			authConfig := expandFlexibleServerAuthConfig(authRaw.([]any))
			if authConfig.PasswordAuth != nil {
				pwdEnabled = *authConfig.PasswordAuth == servers.PasswordBasedAuthEnabled
			}
		}

		if pwdEnabled {
			if !adminLoginSet {
				return fmt.Errorf("`administrator_login` is required when `authentication.password_auth_enabled` is set to `true`")
			}
			if !adminPwdSet && woPassword.IsNull() {
				return fmt.Errorf("`administrator_password` or `administrator_password_wo` is required when `authentication.password_auth_enabled` is set to `true`")
			}
		}

		if d.HasChange("administrator_login") {
			requireUpdateOnLogin = true
			if adminLoginSet && !pwdEnabled {
				return fmt.Errorf("when `administrator_login` is first set, `authentication.password_auth_enabled` must be set to `true`")
			}
		}
	}

	if d.HasChanges("private_dns_zone_id", "public_network_access_enabled") {
		parameters.Properties.Network = expandArmServerNetwork(d)
	}

	var requireFailover bool
	// failover is only supported when `zone` and `high_availability.0.standby_availability_zone` are exchanged with each other
	if d.HasChanges("zone", "high_availability") {
		resp, err := client.Get(ctx, *id)
		if err != nil || resp.Model == nil {
			return err
		}
		props := resp.Model.Properties

		if d.HasChange("zone") {
			if !d.HasChange("high_availability.0.standby_availability_zone") {
				return fmt.Errorf("`zone` can only be changed when exchanged with the zone specified in `high_availability.0.standby_availability_zone`")
			} else {
				// zone can only be changed when it is swapped for an existing high_availability.0.standby_availability_zone - a failover is triggered to make it the new primary availability zone
				// compare current values of zone and high_availability.0.standby_availability_zone with new values and only allow update/failover if the values of zone and an existing high_availability.0.standby_availability_zone have been swapped
				var newZone, newHAStandbyZone string
				newZone = d.Get("zone").(string)
				newHAStandbyZone = d.Get("high_availability.0.standby_availability_zone").(string)
				if props != nil && props.AvailabilityZone != nil && props.HighAvailability != nil && props.HighAvailability.StandbyAvailabilityZone != nil {
					if newZone == *props.HighAvailability.StandbyAvailabilityZone && newHAStandbyZone == *props.AvailabilityZone {
						requireFailover = true
					} else {
						return fmt.Errorf("`zone` can only be changed when exchanged with the zone specified in `high_availability.0.standby_availability_zone`")
					}
				}
			}

			// changes can occur in high_availability.0.standby_availability_zone when zone has not changed in the case where a high_availability block has been newly added or a high_availability block is removed, meaning HA is now disabled
		} else if d.HasChange("high_availability.0.standby_availability_zone") {
			if props != nil && props.HighAvailability != nil && props.HighAvailability.Mode != nil {
				// if HA Mode is currently "ZoneRedundant" and is still set to "ZoneRedundant", high_availability.0.standby_availability_zone cannot be changed
				if *props.HighAvailability.Mode == servers.PostgreSqlFlexibleServerHighAvailabilityModeZoneRedundant && !d.HasChange("high_availability.0.mode") {
					return fmt.Errorf("an existing `high_availability.0.standby_availability_zone` can only be changed when exchanged with the zone specified in `zone`")
				}
				// if high_availability.0.mode changes from "ZoneRedundant", an existing high_availability block has been removed as this is a required field
				// if high_availability.0.mode is not currently "ZoneRedundant", this must be a newly added block
			}
		}
	}

	if d.HasChange("replication_role") {
		createMode := d.Get("create_mode").(string)
		replicationRole := d.Get("replication_role").(string)
		if createMode == string(servers.CreateModeReplica) && replicationRole == string(servers.ReplicationRoleNone) {
			parameters := servers.ServerForPatch{
				Properties: &servers.ServerPropertiesForPatch{
					ReplicationRole: pointer.To(servers.ReplicationRoleNone),
				},
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating `replication_role` for %s: %+v", *id, err)
			}
		} else {
			return fmt.Errorf("`replication_role` only can be updated to `None` for replica server")
		}
	}

	if d.HasChange("administrator_password") {
		parameters.Properties.AdministratorLoginPassword = pointer.To(d.Get("administrator_password").(string))
	}

	if d.HasChange("administrator_password_wo_version") {
		if !woPassword.IsNull() {
			parameters.Properties.AdministratorLoginPassword = pointer.To(woPassword.AsString())
		}
	}

	if d.HasChange("authentication") {
		parameters.Properties.AuthConfig = expandFlexibleServerAuthConfigForPatch(d.Get("authentication").([]any))
	}

	if d.HasChanges("auto_grow_enabled", "storage_mb", "storage_tier", "storage_iops", "storage_throughput") {
		// TODO remove the additional update after https://github.com/Azure/azure-rest-api-specs/issues/22867 is fixed
		storage := expandArmServerStorage(d)

		storageUpdateParameters := servers.ServerForPatch{
			Properties: &servers.ServerPropertiesForPatch{
				Storage: storage,
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, storageUpdateParameters); err != nil {
			return fmt.Errorf("updating `auto_grow_enabled` / `storage_mb` for %s: %+v", *id, err)
		}
	}

	if d.HasChange("backup_retention_days") {
		parameters.Properties.Backup = expandArmServerBackupForPatch(d)
	}

	if d.HasChange("maintenance_window") {
		parameters.Properties.MaintenanceWindow = expandArmServerMaintenanceWindow(d.Get("maintenance_window").([]any))
	}

	if d.HasChange("sku_name") {
		sku, err := expandFlexibleServerSkuForPatch(d.Get("sku_name").(string))
		if err != nil {
			return fmt.Errorf("expanding `sku_name` for %s: %v", id, err)
		}
		parameters.Sku = sku
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("high_availability") {
		parameters.Properties.HighAvailability = expandFlexibleServerHighAvailabilityForPatch(d.Get("high_availability").([]any), false)
	}

	if d.HasChange("customer_managed_key") {
		parameters.Properties.DataEncryption = expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]any))
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity` for %s: %+v", *id, err)
		}
		parameters.Identity = identity
	}

	if d.HasChange("create_mode") {
		createMode := servers.CreateModeForPatch(d.Get("create_mode").(string))
		parameters.Properties.CreateMode = &createMode
	}

	if d.HasChange("cluster") {
		parameters.Properties.Cluster = expandFlexibleServerCluster(d.Get("cluster").([]any))
	}

	if d.HasChange("version") {
		version := servers.PostgresMajorVersion(d.Get("version").(string))
		parameters.Properties.Version = &version
	}

	if requireUpdateOnLogin {
		// Login password can be set using `administrator_password` or `administrator_password_wo`
		loginPassword := d.Get("administrator_password").(string)
		if !woPassword.IsNull() {
			loginPassword = woPassword.AsString()
		}

		loginParameters := servers.Server{
			Location: location.Normalize(d.Get("location").(string)),
			Properties: &servers.ServerProperties{
				CreateMode:                 pointer.To(servers.CreateModeUpdate),
				AuthConfig:                 expandFlexibleServerAuthConfig(d.Get("authentication").([]any)),
				AdministratorLogin:         pointer.To(d.Get("administrator_login").(string)),
				AdministratorLoginPassword: &loginPassword,
				Network:                    expandArmServerNetwork(d),
			},
		}
		if err = client.CreateOrUpdateThenPoll(ctx, *id, loginParameters); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	if requireFailover {
		restartParameters := servers.RestartParameter{
			RestartWithFailover: pointer.To(true),
			FailoverMode:        pointer.To(servers.FailoverModePlannedFailover),
		}

		if err = client.RestartThenPoll(ctx, *id, restartParameters); err != nil {
			return fmt.Errorf("failing over %s: %+v", *id, err)
		}
	}

	return resourcePostgresqlFlexibleServerRead(d, meta)
}
