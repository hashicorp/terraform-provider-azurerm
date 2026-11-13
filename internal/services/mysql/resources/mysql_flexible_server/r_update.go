// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMysqlFlexibleServerUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	// failover is only supported when `zone` and `standby_availability_zone` is exchanged
	var requireFailover bool
	switch {
	case d.HasChange("zone") && d.HasChange("high_availability.0.standby_availability_zone"):
		resp, err := client.Get(ctx, *id)
		if err != nil {
			return err
		}
		if model := resp.Model; model != nil {
			if props := model.Properties; props != nil {
				zone := d.Get("zone").(string)
				standbyZone := d.Get("high_availability.0.standby_availability_zone").(string)

				if props.AvailabilityZone != nil && props.HighAvailability != nil && props.HighAvailability.StandbyAvailabilityZone != nil {
					if zone == *props.HighAvailability.StandbyAvailabilityZone && standbyZone == *props.AvailabilityZone {
						requireFailover = true
					} else {
						return fmt.Errorf("failover only supports exchange between `zone` and `standby_availability_zone`")
					}
				} else {
					return fmt.Errorf("`standby_availability_zone` cannot be added while changing `zone`")
				}
			}
		}
	case d.HasChange("zone") && !d.HasChange("high_availability.0.standby_availability_zone"):
		return fmt.Errorf("`zone` cannot be changed independently")
	default:
		// No need failover when only `standby_availability_zone` is changed and both `zone` and `standby_availability_zone` aren't changed
		requireFailover = false
	}

	if d.HasChange("replication_role") {
		oldReplicationRole, newReplicationRole := d.GetChange("replication_role")
		if oldReplicationRole == "Replica" && newReplicationRole == "None" {
			parameters := servers.ServerForUpdate{
				Properties: &servers.ServerPropertiesForUpdate{
					ReplicationRole: pointer.To(servers.ReplicationRoleNone),
				},
			}

			if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating `replication_role` for %s: %+v", *id, err)
			}
		} else {
			return fmt.Errorf("`replication_role` only can be updated from `Replica` to `None`")
		}
	}

	// ha Enabled is dependent on storage auto grow Enabled. But when we enabled this two features in one request, it returns bad request.
	// Thus we need to separate these two updates in two requests.
	if d.HasChange("storage") && d.Get("storage.0.auto_grow_enabled").(bool) {
		parameters := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				Storage: expandArmServerStorage(d.Get("storage").([]any)),
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("enabling `auto_grow_enabled` for %s: %+v", *id, err)
		}
	}

	if requireFailover {
		failoverClient := meta.(*clients.Client).MySQL.FlexibleServers.ServerFailover
		failoverID := serverfailover.NewFlexibleServerID(id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)

		if err := failoverClient.ServersFailoverThenPoll(ctx, failoverID); err != nil {
			return fmt.Errorf("failing over %s: %+v", *id, err)
		}
	} else if d.HasChange("high_availability") {
		parameters := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				HighAvailability: &servers.HighAvailability{
					Mode: pointer.To(servers.HighAvailabilityModeDisabled),
				},
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("disabling `high_availability` for %s: %+v", *id, err)
		}

		parameters.Properties.HighAvailability = expandFlexibleServerHighAvailability(d.Get("high_availability").([]any))

		if *parameters.Properties.HighAvailability.Mode != servers.HighAvailabilityModeDisabled {
			if err = client.UpdateThenPoll(ctx, *id, parameters); err != nil {
				return fmt.Errorf("updating `high_availability` for %s: %+v", *id, err)
			}
		}
	}

	parameters := servers.ServerForUpdate{
		Properties: &servers.ServerPropertiesForUpdate{},
	}

	if d.HasChange("administrator_password") {
		parameters.Properties.AdministratorLoginPassword = pointer.To(d.Get("administrator_password").(string))
	}

	if d.HasChange("administrator_password_wo_version") {
		woPassword, err := pluginsdk.GetWriteOnly(d, "administrator_password_wo", cty.String)
		if err != nil {
			return err
		}
		if !woPassword.IsNull() {
			parameters.Properties.AdministratorLoginPassword = pointer.To(woPassword.AsString())
		}
	}

	if d.HasChanges("backup_retention_days", "geo_redundant_backup_enabled") {
		parameters.Properties.Backup = expandArmServerBackup(d)
	}

	if d.HasChange("customer_managed_key") {
		parameters.Properties.DataEncryption = expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]any))
	}

	if d.HasChange("identity") {
		identity, err := expandFlexibleServerIdentity(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = identity
	}

	if d.HasChange("maintenance_window") {
		parameters.Properties.MaintenanceWindow = expandArmServerMaintenanceWindow(d.Get("maintenance_window").([]any))
	}

	if d.HasChange("sku_name") {
		sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
		if err != nil {
			return fmt.Errorf("expanding `sku_name`: %+v", err)
		}
		parameters.Sku = sku
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]any))
	}

	if d.HasChange("public_network_access") {
		if parameters.Properties.Network == nil {
			parameters.Properties.Network = &servers.Network{}
		}
		parameters.Properties.Network.PublicNetworkAccess = pointer.ToEnum[servers.EnableStatusEnum](d.Get("public_network_access").(string))
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("storage") && !d.Get("storage.0.auto_grow_enabled").(bool) {
		// log_on_disk_enabled must be updated first when auto_grow_enabled and log_on_disk_enabled are updated from true to false in one request
		if oldLogOnDiskEnabled, newLogOnDiskEnabled := d.GetChange("storage.0.log_on_disk_enabled"); oldLogOnDiskEnabled.(bool) && !newLogOnDiskEnabled.(bool) {
			if oldAutoGrowEnabled, newAutoGrowEnabled := d.GetChange("storage.0.auto_grow_enabled"); oldAutoGrowEnabled.(bool) && !newAutoGrowEnabled.(bool) {
				parameters := servers.ServerForUpdate{
					Properties: &servers.ServerPropertiesForUpdate{
						Storage: &servers.Storage{
							LogOnDisk: pointer.To(servers.EnableStatusEnumDisabled),
						},
					},
				}
				if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
					return fmt.Errorf("disabling `log_on_disk_enabled` for %s: %+v", *id, err)
				}
			}
		}

		parameters := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				Storage: expandArmServerStorage(d.Get("storage").([]any)),
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("disabling `auto_grow_enabled` for %s: %+v", *id, err)
		}
	}

	if d.HasChange("version") {
		parameters := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				Version: pointer.ToEnum[servers.ServerVersion](d.Get("version").(string)),
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("updating `version` for %s: %+v", *id, err)
		}
	}

	return resourceMysqlFlexibleServerRead(d, meta)
}
