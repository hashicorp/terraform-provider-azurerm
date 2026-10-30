// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package postgresql_flexible_server

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/postgresql/2025-08-01/servers"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/postgres/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePostgresqlFlexibleServerCreate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).Postgres.FlexibleServersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewFlexibleServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_postgresql_flexible_server", id.ID())
		}
	}

	createMode := d.Get("create_mode").(string)

	if _, ok := d.GetOk("replication_role"); ok {
		return fmt.Errorf("`replication_role` cannot be set while creating")
	}

	if servers.CreateMode(createMode) == servers.CreateModePointInTimeRestore || servers.CreateMode(createMode) == servers.CreateModeGeoRestore {
		if _, ok := d.GetOk("source_server_id"); !ok {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is  %s", createMode)
		}
		if _, ok := d.GetOk("point_in_time_restore_time_in_utc"); !ok {
			return fmt.Errorf("`point_in_time_restore_time_in_utc` is required when `create_mode` is  %s", createMode)
		}
	}

	if servers.CreateMode(createMode) == servers.CreateModeReplica {
		if _, ok := d.GetOk("source_server_id"); !ok {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `Replica`")
		}
	}

	woPassword, err := pluginsdk.GetWriteOnly(d, "administrator_password_wo", cty.String)
	if err != nil {
		return err
	}

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
				return fmt.Errorf("`administrator_login` is required when `create_mode` is `Default` and `authentication.password_auth_enabled` is set to `true`")
			}

			if !adminPwdSet && woPassword.IsNull() {
				return fmt.Errorf("`administrator_password` or `administrator_password_wo` is required when `create_mode` is `Default` and `authentication.password_auth_enabled` is set to `true`")
			}
		} else if adminLoginSet || adminPwdSet || !woPassword.IsNull() {
			return fmt.Errorf("`administrator_login`, `administrator_password` and `administrator_password_wo` cannot be set during creation when `authentication.password_auth_enabled` is set to `false`")
		}

		if _, ok := d.GetOk("sku_name"); !ok {
			return fmt.Errorf("`sku_name` is required when `create_mode` is `Default`")
		}

		if _, ok := d.GetOk("version"); !ok {
			return fmt.Errorf("`version` is required when `create_mode` is `Default`")
		}
	}

	sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for %s: %v", id, err)
	}

	storage := expandArmServerStorage(d)
	var storageMb int

	if storage.StorageSizeGB == nil || *storage.StorageSizeGB == 0 {
		// set the default value for storage_mb...
		storageMb = 32768
		storage.StorageSizeGB = pointer.To(int64(32))
		log.Printf("[DEBUG]: Default 'storage_mb' Set -> %d\n", storageMb)
	} else {
		storageMb = int(*storage.StorageSizeGB) * 1024
	}

	if pointer.From(storage.Type) != servers.StorageTypePremiumVTwoLRS && pointer.From(storage.Tier) == "" {
		storageTierMappings := validate.InitializeFlexibleServerStorageTierDefaults()
		storageTiers := storageTierMappings[storageMb]
		storage.Tier = pointer.To(storageTiers.DefaultTier)
		log.Printf("[DEBUG]: Default 'storage_tier' Set -> %q\n", storageTiers.DefaultTier)
	}

	parameters := servers.Server{
		Location: location.Normalize(d.Get("location").(string)),
		Properties: &servers.ServerProperties{
			Network:          expandArmServerNetwork(d),
			Storage:          storage,
			HighAvailability: expandFlexibleServerHighAvailability(d.Get("high_availability").([]any), true),
			Backup:           expandArmServerBackup(d),
			DataEncryption:   expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]any)),
			Cluster:          expandFlexibleServerCluster(d.Get("cluster").([]any)),
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

	if createMode != "" {
		parameters.Properties.CreateMode = pointer.ToEnum[servers.CreateMode](createMode)
	}

	if v, ok := d.GetOk("version"); ok && v.(string) != "" {
		parameters.Properties.Version = pointer.ToEnum[servers.PostgresMajorVersion](v.(string))
	}

	if v, ok := d.GetOk("zone"); ok && v.(string) != "" {
		parameters.Properties.AvailabilityZone = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("source_server_id"); ok && v.(string) != "" {
		// The source server will be Updating status when creating a replica
		sourceServerId, _ := servers.ParseFlexibleServerID(v.(string))
		locks.ByName(sourceServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)
		defer locks.UnlockByName(sourceServerId.FlexibleServerName, postgresqlFlexibleServerResourceName)

		parameters.Properties.SourceServerResourceId = pointer.To(v.(string))
	}

	pointInTimeUTC := d.Get("point_in_time_restore_time_in_utc").(string)
	if pointInTimeUTC != "" {
		v, err := time.Parse(time.RFC3339, pointInTimeUTC)
		if err != nil {
			return fmt.Errorf("unable to parse `point_in_time_restore_time_in_utc` value")
		}
		parameters.Properties.SetPointInTimeUTCAsTime(v)
	}

	if authRaw, ok := d.GetOk("authentication"); ok {
		parameters.Properties.AuthConfig = expandFlexibleServerAuthConfig(authRaw.([]any))
	}

	identity, err := identity.ExpandLegacySystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`")
	}
	parameters.Identity = identity

	if err = client.CreateOrUpdateCallbackThenPoll(ctx, id, parameters, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	requireAdditionalUpdate := false
	updateProperties := servers.ServerPropertiesForPatch{}
	// `maintenance_window` could only be updated with, could not be created with
	if v, ok := d.GetOk("maintenance_window"); ok {
		requireAdditionalUpdate = true
		updateProperties.MaintenanceWindow = expandArmServerMaintenanceWindow(v.([]any))
	}

	if requireAdditionalUpdate {
		update := servers.ServerForPatch{
			Properties: &updateProperties,
		}
		if err = client.UpdateThenPoll(ctx, id, update); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourcePostgresqlFlexibleServerRead(d, meta)
}
