// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Postgresql Flexibleserver %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("administrator_password_wo_version", d.Get("administrator_password_wo_version").(int))

	if model := resp.Model; model != nil {
		d.Set("location", location.Normalize(model.Location))

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin) // if pwdEnabled is set to false, then the service does not return the value of AdministratorLogin
			d.Set("zone", props.AvailabilityZone)
			d.Set("version", pointer.From(props.Version))
			d.Set("fqdn", props.FullyQualifiedDomainName)

			// According to the API spec, `sourceServerResourceId`(`source_server_id`) is only returned by the Azure REST API
			// when `create_mode` is 'Replica'. For other create modes, this field is not returned, which is intended behaviour of the API.
			// Therefore, we populate this field from the API response if present; otherwise, we read the value from the configuration.
			sourceResourceId := pointer.From(props.SourceServerResourceId)
			if sourceResourceId == "" {
				sourceResourceId = d.Get("source_server_id").(string)
			}
			d.Set("source_server_id", sourceResourceId)

			// `pointInTimeUTC` has "x-ms-mutability": ["create"], so it's only used at creation.
			// Without "read," it's not persisted in the resource model and won't appear in GET responses.
			// So read it from the configuration.
			d.Set("point_in_time_restore_time_in_utc", d.Get("point_in_time_restore_time_in_utc").(string))

			// Currently, `replicationRole` is set to `Primary` when `createMode` is `Replica` and `replicationRole` is updated to `None`. Service team confirmed it should be set to `None` for this scenario. See more details from https://github.com/Azure/azure-rest-api-specs/issues/22499
			d.Set("replication_role", d.Get("replication_role").(string))

			if network := props.Network; network != nil {
				publicNetworkAccess := false
				if network.PublicNetworkAccess != nil {
					publicNetworkAccess = *network.PublicNetworkAccess == servers.ServerPublicNetworkAccessStateEnabled
				}
				d.Set("public_network_access_enabled", publicNetworkAccess)
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				d.Set("private_dns_zone_id", network.PrivateDnsZoneArmResourceId)
			}

			if err := d.Set("maintenance_window", flattenArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
				return fmt.Errorf("setting `maintenance_window`: %+v", err)
			}

			if storage := props.Storage; storage != nil {
				if storage.AutoGrow != nil {
					d.Set("auto_grow_enabled", *storage.AutoGrow == servers.StorageAutoGrowEnabled)
				}

				if storage.StorageSizeGB != nil {
					d.Set("storage_mb", (*storage.StorageSizeGB * 1024))
				}

				if storage.Tier != nil {
					d.Set("storage_tier", string(*storage.Tier))
				}

				if storage.Type != nil {
					d.Set("storage_type", string(*storage.Type))
				}

				if storage.Iops != nil {
					d.Set("storage_iops", *storage.Iops)
				}

				if storage.Throughput != nil {
					d.Set("storage_throughput", *storage.Throughput)
				}
			}

			if backup := props.Backup; backup != nil {
				d.Set("backup_retention_days", backup.BackupRetentionDays)

				geoRedundantBackup := false
				if backup.GeoRedundantBackup != nil {
					geoRedundantBackup = *backup.GeoRedundantBackup == servers.GeographicallyRedundantBackupEnabled
				}
				d.Set("geo_redundant_backup_enabled", geoRedundantBackup)
			}

			if err := d.Set("high_availability", flattenFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
				return fmt.Errorf("setting `high_availability`: %+v", err)
			}

			if pointer.From(props.SourceServerResourceId) == "" {
				if err := d.Set("cluster", flattenFlexibleServerCluster(props.Cluster)); err != nil {
					return fmt.Errorf("setting `cluster`: %+v", err)
				}
			} else {
				if err := d.Set("cluster", []any{}); err != nil {
					return fmt.Errorf("setting `cluster`: %+v", err)
				}
			}

			if props.AuthConfig != nil {
				d.Set("authentication", flattenFlexibleServerAuthConfig(props.AuthConfig))
			}

			cmk, err := flattenFlexibleServerDataEncryption(props.DataEncryption)
			if err != nil {
				return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
			}
			if err := d.Set("customer_managed_key", cmk); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			identity, err := identity.FlattenLegacySystemAndUserAssignedMap(model.Identity)
			if err != nil {
				return fmt.Errorf("flattening `identity`: %+v", err)
			}
			if err := d.Set("identity", identity); err != nil {
				return fmt.Errorf("setting `identity`: %+v", err)
			}
		}

		sku, err := flattenFlexibleServerSku(model.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name` for %s: %v", id, err)
		}

		d.Set("sku_name", sku)

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

const (
	ServerMaintenanceWindowEnabled  = "Enabled"
	ServerMaintenanceWindowDisabled = "Disabled"
)
