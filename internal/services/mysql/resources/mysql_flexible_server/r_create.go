// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package mysql_flexible_server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceMysqlFlexibleServerCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := servers.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	{
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mysql_flexible_server", id.ID())
		}
	}

	woPassword, err := pluginsdk.GetWriteOnly(d, "administrator_password_wo", cty.String)
	if err != nil {
		return err
	}

	createMode := servers.CreateMode(d.Get("create_mode").(string))

	if _, ok := d.GetOk("replication_role"); ok {
		return fmt.Errorf("`replication_role` cannot be set while creating")
	}

	if _, ok := d.GetOk("source_server_id"); !ok {
		if createMode == servers.CreateModePointInTimeRestore || createMode == servers.CreateModeReplica || createMode == servers.CreateModeGeoRestore {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `PointInTimeRestore`, `GeoRestore`, or `Replica`")
		}
	}

	if createMode == servers.CreateModePointInTimeRestore {
		if _, ok := d.GetOk("point_in_time_restore_time_in_utc"); !ok {
			return fmt.Errorf("`point_in_time_restore_time_in_utc` is required when `create_mode` is `PointInTimeRestore`")
		}
	}

	if createMode == "" || createMode == servers.CreateModeDefault {
		if _, ok := d.GetOk("administrator_login"); !ok {
			return fmt.Errorf("`administrator_login` is required when `create_mode` is `Default`")
		}

		if _, ok := d.GetOk("administrator_password"); !ok && woPassword.IsNull() {
			return fmt.Errorf("`administrator_password_wo` or `administrator_password` is required when `create_mode` is `Default`")
		}

		if _, ok := d.GetOk("sku_name"); !ok {
			return fmt.Errorf("`sku_name` is required when `create_mode` is `Default`")
		}
	}

	storageSettings := expandArmServerStorage(d.Get("storage").([]any))
	if storageSettings != nil {
		if storageSettings.Iops != nil && *storageSettings.AutoIoScaling == servers.EnableStatusEnumEnabled {
			return fmt.Errorf("`iops` can not be set if `io_scaling_enabled` is set to true")
		}
	}

	sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for %s: %+v", id, err)
	}

	version := servers.ServerVersion(d.Get("version").(string))
	parameters := servers.Server{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &servers.ServerProperties{
			CreateMode:       &createMode,
			Version:          &version,
			Storage:          expandArmServerStorage(d.Get("storage").([]any)),
			Network:          expandArmServerNetwork(d),
			HighAvailability: expandFlexibleServerHighAvailability(d.Get("high_availability").([]any)),
			Backup:           expandArmServerBackup(d),
			DataEncryption:   expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]any)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if v, ok := d.GetOk("administrator_login"); ok && v.(string) != "" {
		parameters.Properties.AdministratorLogin = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("administrator_password"); ok && v.(string) != "" {
		parameters.Properties.AdministratorLoginPassword = pointer.To(v.(string))
	}

	if !woPassword.IsNull() {
		parameters.Properties.AdministratorLoginPassword = pointer.To(woPassword.AsString())
	}

	if v, ok := d.GetOk("zone"); ok && v.(string) != "" {
		parameters.Properties.AvailabilityZone = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("source_server_id"); ok && v.(string) != "" {
		parameters.Properties.SourceServerResourceId = pointer.To(v.(string))
	}

	pointInTimeUTC := d.Get("point_in_time_restore_time_in_utc").(string)
	if pointInTimeUTC != "" {
		v, err := time.Parse(time.RFC3339, pointInTimeUTC)
		if err != nil {
			return fmt.Errorf("unable to parse `point_in_time_restore_time_in_utc` value")
		}
		parameters.Properties.SetRestorePointInTimeAsTime(v)
	}

	identity, err := expandFlexibleServerIdentity(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`")
	}
	parameters.Identity = identity

	if err := client.CreateCallbackThenPoll(ctx, id, parameters, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	// Account for eventual consistency until issue https://github.com/Azure/azure-rest-api-specs/issues/21178 is fixed.
	poller := custompollers.NewEventualConsistencyPoller(1, func(pollerCtx context.Context) (*http.Response, error) {
		resp, err := client.Get(pollerCtx, id)
		return resp.HttpResponse, err
	}, custompollers.DefaultCreationEventualConsistencyPollerOptions())
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Mysql Flexible Server %s: %+v", id, err)
	}

	// `maintenance_window` could only be updated with, could not be created with
	if v, ok := d.GetOk("maintenance_window"); ok {
		mwParams := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				MaintenanceWindow: expandArmServerMaintenanceWindow(v.([]any)),
			},
		}
		if err := client.UpdateThenPoll(ctx, id, mwParams); err != nil {
			return fmt.Errorf("updating Maintenance Window for %s: %+v", id, err)
		}
	}

	return resourceMysqlFlexibleServerRead(d, meta)
}
