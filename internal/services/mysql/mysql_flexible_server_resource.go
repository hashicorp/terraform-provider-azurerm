// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/serverfailover"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2023-12-30/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/privatedns/2024-06-01/privatezones"
	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	keyVaultValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/keyvault/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

const (
	ServerMaintenanceWindowEnabled  = "Enabled"
	ServerMaintenanceWindowDisabled = "Disabled"
)

var mysqlFlexibleServerResourceName = "azurerm_mysql_flexible_server"

func resourceMysqlFlexibleServer() *pluginsdk.Resource {
	resource := &pluginsdk.Resource{
		Create: resourceMysqlFlexibleServerCreate,
		Read:   resourceMysqlFlexibleServerRead,
		Update: resourceMysqlFlexibleServerUpdate,
		Delete: resourceMysqlFlexibleServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(2 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(2 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := servers.ParseFlexibleServerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"administrator_login": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerAdministratorLogin,
			},

			"administrator_password": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				Sensitive:     true,
				ValidateFunc:  validate.FlexibleServerAdministratorPassword,
				ConflictsWith: []string{"administrator_password_wo"},
			},

			"administrator_password_wo": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				WriteOnly:     true,
				ConflictsWith: []string{"administrator_password"},
				RequiredWith:  []string{"administrator_password_wo_version"},
				ValidateFunc:  validate.FlexibleServerAdministratorPassword,
			},

			"administrator_password_wo_version": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				RequiredWith: []string{"administrator_password_wo"},
			},

			"backup_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Default:      7,
				ValidateFunc: validation.IntBetween(1, 35),
			},

			"create_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.CreateModeDefault),
					string(servers.CreateModeGeoRestore),
					string(servers.CreateModePointInTimeRestore),
					string(servers.CreateModeReplica),
				}, false),
			},

			"customer_managed_key": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
							RequiredWith: []string{
								"identity",
								"customer_managed_key.0.primary_user_assigned_identity_id",
							},
						},
						"primary_user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
						"geo_backup_key_vault_key_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: keyVaultValidate.NestedItemIdWithOptionalVersion,
							RequiredWith: []string{
								"identity",
								"customer_managed_key.0.geo_backup_user_assigned_identity_id",
							},
						},
						"geo_backup_user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateUserAssignedIdentityID,
						},
					},
				},
			},

			"delegated_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: commonids.ValidateSubnetID,
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"high_availability": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"mode": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(servers.HighAvailabilityModeZoneRedundant),
								string(servers.HighAvailabilityModeSameZone),
							}, false),
						},

						"standby_availability_zone": commonschema.ZoneSingleOptionalComputed(),
					},
				},
			},

			"identity": commonschema.UserAssignedIdentityOptional(),

			"maintenance_window": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"day_of_week": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 6),
						},

						"start_hour": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 23),
						},

						"start_minute": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Default:      0,
							ValidateFunc: validation.IntBetween(0, 59),
						},
					},
				},
			},

			"point_in_time_restore_time_in_utc": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"private_dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: privatezones.ValidatePrivateDnsZoneID,
			},

			"public_network_access": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				// NOTE: O+C: Azure normally defaults this to `Enabled` unless values are provided for `delegated_subnet_id` and `private_dns_zone_id`
				Computed:     true,
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForEnableStatusEnum(), false),
			},

			"replication_role": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.ReplicationRoleNone),
				}, false),
			},

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.FlexibleServerSkuName,
			},

			"source_server_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: servers.ValidateFlexibleServerID,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"auto_grow_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  true,
						},

						"iops": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(360, 48000),
						},

						"log_on_disk_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},

						"size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(20, 16384),
						},
						"io_scaling_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},

			"version": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.ServerVersionFivePointSeven),
					string(servers.ServerVersionEightPointZeroPointTwoOne),
				}, false),
			},

			"zone": commonschema.ZoneSingleOptionalComputed(),

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"replica_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomDiffWithAll(
			pluginsdk.ForceNewIfChange("storage.0.size_gb", func(ctx context.Context, old, new, meta interface{}) bool {
				return new.(int) < old.(int)
			}),
		),
	}

	if !features.FivePointOh() {
		resource.Schema["public_network_access_enabled"] = &pluginsdk.Schema{
			Type:     pluginsdk.TypeBool,
			Computed: true,
		}
	}

	return resource
}

func resourceMysqlFlexibleServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := servers.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_mysql_flexible_server", id.ID())
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

	storageSettings := expandArmServerStorage(d.Get("storage").([]interface{}))
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
			Storage:          expandArmServerStorage(d.Get("storage").([]interface{})),
			Network:          expandArmServerNetwork(d),
			HighAvailability: expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{})),
			Backup:           expandArmServerBackup(d),
			DataEncryption:   expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]interface{})),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
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

	identity, err := expandFlexibleServerIdentity(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`")
	}
	parameters.Identity = identity

	if err := client.CreateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Add the state wait function until issue https://github.com/Azure/azure-rest-api-specs/issues/21178 is fixed.
	stateConf := &pluginsdk.StateChangeConf{
		Pending: []string{
			"Pending",
		},
		Target: []string{
			"OK",
		},
		Refresh:    mySqlFlexibleServerCreationRefreshFunc(ctx, client, id),
		MinTimeout: 10 * time.Second,
		Timeout:    d.Timeout(pluginsdk.TimeoutCreate),
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for creation of Mysql Flexible Server %s: %+v", id, err)
	}

	// `maintenance_window` could only be updated with, could not be created with
	if v, ok := d.GetOk("maintenance_window"); ok {
		mwParams := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				MaintenanceWindow: expandArmServerMaintenanceWindow(v.([]interface{})),
			},
		}
		if err := client.UpdateThenPoll(ctx, id, mwParams); err != nil {
			return fmt.Errorf("updating Maintenance Window for %s: %+v", id, err)
		}
	}

	d.SetId(id.ID())

	return resourceMysqlFlexibleServerRead(d, meta)
}

func mySqlFlexibleServerCreationRefreshFunc(ctx context.Context, client *servers.ServersClient, id servers.FlexibleServerId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := client.Get(ctx, id)
		if err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return resp, "Pending", nil
			}
			return resp, "Error", err
		}
		return "OK", "OK", nil
	}
}

func resourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] Mysql Flexible Server %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.FlexibleServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(&model.Location))
		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("zone", props.AvailabilityZone)
			d.Set("version", string(pointer.From(props.Version)))
			d.Set("fqdn", props.FullyQualifiedDomainName)
			d.Set("source_server_id", props.SourceServerResourceId)

			if network := props.Network; network != nil {
				d.Set("delegated_subnet_id", network.DelegatedSubnetResourceId)
				d.Set("private_dns_zone_id", network.PrivateDnsZoneResourceId)
				d.Set("public_network_access", string(pointer.From(network.PublicNetworkAccess)))

				if !features.FivePointOh() {
					d.Set("public_network_access_enabled", pointer.From(network.PublicNetworkAccess) == servers.EnableStatusEnumEnabled)
				}
			}

			cmk, err := flattenFlexibleServerDataEncryption(props.DataEncryption)
			if err != nil {
				return fmt.Errorf("flattening `customer_managed_key`: %+v", err)
			}
			if err := d.Set("customer_managed_key", cmk); err != nil {
				return fmt.Errorf("setting `customer_managed_key`: %+v", err)
			}

			identity, err := flattenFlexibleServerIdentity(model.Identity)
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
			d.Set("replication_role", string(pointer.From(props.ReplicationRole)))
			d.Set("replica_capacity", props.ReplicaCapacity)
		}
		sku, err := flattenFlexibleServerSku(model.Sku)
		if err != nil {
			return fmt.Errorf("flattening `sku_name`: %+v", err)
		}
		d.Set("sku_name", sku)

		d.Set("administrator_password_wo_version", d.Get("administrator_password_wo_version").(int))

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceMysqlFlexibleServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			replicationRole := servers.ReplicationRoleNone
			parameters := servers.ServerForUpdate{
				Properties: &servers.ServerPropertiesForUpdate{
					ReplicationRole: &replicationRole,
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
				Storage: expandArmServerStorage(d.Get("storage").([]interface{})),
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
		mode := servers.HighAvailabilityModeDisabled
		parameters := servers.ServerForUpdate{
			Properties: &servers.ServerPropertiesForUpdate{
				HighAvailability: &servers.HighAvailability{
					Mode: &mode,
				},
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("disabling `high_availability` for %s: %+v", *id, err)
		}

		parameters.Properties.HighAvailability = expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{}))

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

	if d.HasChange("backup_retention_days") || d.HasChange("geo_redundant_backup_enabled") {
		parameters.Properties.Backup = expandArmServerBackup(d)
	}

	if d.HasChange("customer_managed_key") {
		parameters.Properties.DataEncryption = expandFlexibleServerDataEncryption(d.Get("customer_managed_key").([]interface{}))
	}

	if d.HasChange("identity") {
		identity, err := expandFlexibleServerIdentity(d.Get("identity").([]interface{}))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		parameters.Identity = identity
	}

	if d.HasChange("maintenance_window") {
		parameters.Properties.MaintenanceWindow = expandArmServerMaintenanceWindow(d.Get("maintenance_window").([]interface{}))
	}

	if d.HasChange("sku_name") {
		sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
		if err != nil {
			return fmt.Errorf("expanding `sku_name`: %+v", err)
		}
		parameters.Sku = sku
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if d.HasChange("public_network_access") {
		if parameters.Properties.Network == nil {
			parameters.Properties.Network = &servers.Network{}
		}
		parameters.Properties.Network.PublicNetworkAccess = pointer.To(servers.EnableStatusEnum(d.Get("public_network_access").(string)))
	}

	if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if d.HasChange("storage") && !d.Get("storage.0.auto_grow_enabled").(bool) {
		// log_on_disk_enabled must be updated first when auto_grow_enabled and log_on_disk_enabled are updated from true to false in one request
		if oldLogOnDiskEnabled, newLogOnDiskEnabled := d.GetChange("storage.0.log_on_disk_enabled"); oldLogOnDiskEnabled.(bool) && !newLogOnDiskEnabled.(bool) {
			if oldAutoGrowEnabled, newAutoGrowEnabled := d.GetChange("storage.0.auto_grow_enabled"); oldAutoGrowEnabled.(bool) && !newAutoGrowEnabled.(bool) {
				logOnDiskDisabled := servers.EnableStatusEnumDisabled
				parameters := servers.ServerForUpdate{
					Properties: &servers.ServerPropertiesForUpdate{
						Storage: &servers.Storage{
							LogOnDisk: &logOnDiskDisabled,
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
				Storage: expandArmServerStorage(d.Get("storage").([]interface{})),
			},
		}

		if err := client.UpdateThenPoll(ctx, *id, parameters); err != nil {
			return fmt.Errorf("disabling `auto_grow_enabled` for %s: %+v", *id, err)
		}
	}

	return resourceMysqlFlexibleServerRead(d, meta)
}

func resourceMysqlFlexibleServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServers.Servers
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseFlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandArmServerNetwork(d *pluginsdk.ResourceData) *servers.Network {
	network := servers.Network{}

	if v, ok := d.GetOk("delegated_subnet_id"); ok {
		network.DelegatedSubnetResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		network.PrivateDnsZoneResourceId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("public_network_access"); ok {
		network.PublicNetworkAccess = pointer.To(servers.EnableStatusEnum(v.(string)))
	}

	return &network
}

func expandArmServerMaintenanceWindow(input []interface{}) *servers.MaintenanceWindow {
	if len(input) == 0 {
		return &servers.MaintenanceWindow{
			CustomWindow: pointer.To(ServerMaintenanceWindowDisabled),
		}
	}
	v := input[0].(map[string]interface{})

	maintenanceWindow := servers.MaintenanceWindow{
		CustomWindow: pointer.To(ServerMaintenanceWindowEnabled),
		StartHour:    pointer.To(int64(v["start_hour"].(int))),
		StartMinute:  pointer.To(int64(v["start_minute"].(int))),
		DayOfWeek:    pointer.To(int64(v["day_of_week"].(int))),
	}

	return &maintenanceWindow
}

func expandArmServerStorage(inputs []interface{}) *servers.Storage {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
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

func flattenArmServerStorage(storage *servers.Storage) []interface{} {
	if storage == nil {
		return []interface{}{}
	}

	var size, iops int64
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []interface{}{
		map[string]interface{}{
			"size_gb":             size,
			"iops":                iops,
			"auto_grow_enabled":   *storage.AutoGrow == servers.EnableStatusEnumEnabled,
			"io_scaling_enabled":  *storage.AutoIoScaling == servers.EnableStatusEnumEnabled,
			"log_on_disk_enabled": *storage.LogOnDisk == servers.EnableStatusEnumEnabled,
		},
	}
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

func flattenArmServerMaintenanceWindow(input *servers.MaintenanceWindow) []interface{} {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]interface{}, 0)
	}

	var dayOfWeek int64
	if input.DayOfWeek != nil {
		dayOfWeek = *input.DayOfWeek
	}
	var startHour int64
	if input.StartHour != nil {
		startHour = *input.StartHour
	}
	var startMinute int64
	if input.StartMinute != nil {
		startMinute = *input.StartMinute
	}
	return []interface{}{
		map[string]interface{}{
			"day_of_week":  dayOfWeek,
			"start_hour":   startHour,
			"start_minute": startMinute,
		},
	}
}

func expandFlexibleServerHighAvailability(inputs []interface{}) *servers.HighAvailability {
	if len(inputs) == 0 || inputs[0] == nil {
		highAvailability := servers.HighAvailabilityModeDisabled
		return &servers.HighAvailability{
			Mode: &highAvailability,
		}
	}

	input := inputs[0].(map[string]interface{})

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

func flattenFlexibleServerHighAvailability(ha *servers.HighAvailability) []interface{} {
	if ha == nil || *ha.Mode == servers.HighAvailabilityModeDisabled {
		return []interface{}{}
	}

	var zone string
	if ha.StandbyAvailabilityZone != nil {
		zone = *ha.StandbyAvailabilityZone
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                      string(*ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}

func expandFlexibleServerDataEncryption(input []interface{}) *servers.DataEncryption {
	if len(input) == 0 || input[0] == nil {
		det := servers.DataEncryptionTypeSystemManaged
		return &servers.DataEncryption{
			Type: &det,
		}
	}
	v := input[0].(map[string]interface{})

	det := servers.DataEncryptionTypeAzureKeyVault
	dataEncryption := servers.DataEncryption{
		Type: &det,
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

func flattenFlexibleServerDataEncryption(de *servers.DataEncryption) ([]interface{}, error) {
	if de == nil || *de.Type == servers.DataEncryptionTypeSystemManaged {
		return []interface{}{}, nil
	}

	item := map[string]interface{}{}
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

	return []interface{}{item}, nil
}

func expandFlexibleServerIdentity(input []interface{}) (*servers.MySQLServerIdentity, error) {
	expanded, err := identity.ExpandUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	identityType := servers.ManagedServiceIdentityType(string(expanded.Type))
	out := servers.MySQLServerIdentity{
		Type: &identityType,
	}
	if expanded.Type == identity.TypeUserAssigned {
		ids := make(map[string]interface{})
		for k := range expanded.IdentityIds {
			ids[k] = struct{}{}
		}
		out.UserAssignedIdentities = &ids
	}

	return &out, nil
}

func flattenFlexibleServerIdentity(input *servers.MySQLServerIdentity) (*[]interface{}, error) {
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
