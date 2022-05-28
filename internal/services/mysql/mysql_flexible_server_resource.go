package mysql

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2021-05-01/mysqlflexibleservers"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/sdk/2018-09-01/privatezones"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	ServerMaintenanceWindowEnabled  = "Enabled"
	ServerMaintenanceWindowDisabled = "Disabled"
)

func resourceMysqlFlexibleServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMysqlFlexibleServerCreate,
		Read:   resourceMysqlFlexibleServerRead,
		Update: resourceMysqlFlexibleServerUpdate,
		Delete: resourceMysqlFlexibleServerDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(1 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(1 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(1 * time.Hour),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.FlexibleServerID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"administrator_login": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerAdministratorLogin,
			},

			"administrator_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validate.FlexibleServerAdministratorPassword,
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
					string(mysqlflexibleservers.CreateModeDefault),
					string(mysqlflexibleservers.CreateModeGeoRestore),
					string(mysqlflexibleservers.CreateModePointInTimeRestore),
					string(mysqlflexibleservers.CreateModeReplica),
				}, false),
			},

			"delegated_subnet_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.SubnetID,
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
								string(mysqlflexibleservers.HighAvailabilityModeZoneRedundant),
								string(mysqlflexibleservers.HighAvailabilityModeSameZone),
							}, false),
						},

						"standby_availability_zone": commonschema.ZoneSingleOptional(),
					},
				},
			},

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

			"replication_role": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mysqlflexibleservers.ReplicationRoleNone),
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
				ValidateFunc: validate.FlexibleServerID,
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
							ValidateFunc: validation.IntBetween(360, 20000),
						},

						"size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(20, 16384),
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
					string(mysqlflexibleservers.ServerVersionFiveFullStopSeven),
					string(mysqlflexibleservers.ServerVersionEightFullStopZeroFullStopTwoOne),
				}, false),
			},

			"zone": commonschema.ZoneSingleOptional(),

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"replica_capacity": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceMysqlFlexibleServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).MySQL.FlexibleServerClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	id := parse.NewFlexibleServerID(subscriptionId, resourceGroup, name)

	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for present of existing Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}
	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_mysql_flexible_server", id.ID())
	}

	createMode := mysqlflexibleservers.CreateMode(d.Get("create_mode").(string))

	if _, ok := d.GetOk("replication_role"); ok {
		return fmt.Errorf("`replication_role` cannot be set while creating")
	}

	if _, ok := d.GetOk("source_server_id"); !ok {
		if createMode == mysqlflexibleservers.CreateModePointInTimeRestore || createMode == mysqlflexibleservers.CreateModeReplica || createMode == mysqlflexibleservers.CreateModeGeoRestore {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `PointInTimeRestore`, `GeoRestore`, or `Replica`")
		}
	}

	if createMode == mysqlflexibleservers.CreateModePointInTimeRestore {
		if _, ok := d.GetOk("point_in_time_restore_time_in_utc"); !ok {
			return fmt.Errorf("`point_in_time_restore_time_in_utc` is required when `create_mode` is `PointInTimeRestore`")
		}
	}

	if createMode == "" || createMode == mysqlflexibleservers.CreateModeDefault {
		if _, ok := d.GetOk("administrator_login"); !ok {
			return fmt.Errorf("`administrator_login` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("administrator_password"); !ok {
			return fmt.Errorf("`administrator_password` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("sku_name"); !ok {
			return fmt.Errorf("`sku_name` is required when `create_mode` is `Default`")
		}
	}

	sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for MySql Flexible Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	parameters := mysqlflexibleservers.Server{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ServerProperties: &mysqlflexibleservers.ServerProperties{
			CreateMode:       createMode,
			Version:          mysqlflexibleservers.ServerVersion(d.Get("version").(string)),
			Storage:          expandArmServerStorage(d.Get("storage").([]interface{})),
			Network:          expandArmServerNetwork(d),
			HighAvailability: expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{})),
			Backup:           expandArmServerBackup(d),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if v, ok := d.GetOk("administrator_login"); ok && v.(string) != "" {
		parameters.ServerProperties.AdministratorLogin = utils.String(v.(string))
	}

	if v, ok := d.GetOk("administrator_password"); ok && v.(string) != "" {
		parameters.ServerProperties.AdministratorLoginPassword = utils.String(v.(string))
	}

	if v, ok := d.GetOk("zone"); ok && v.(string) != "" {
		parameters.ServerProperties.AvailabilityZone = utils.String(v.(string))
	}

	if v, ok := d.GetOk("source_server_id"); ok && v.(string) != "" {
		parameters.SourceServerResourceID = utils.String(v.(string))
	}

	pointInTimeUTC := d.Get("point_in_time_restore_time_in_utc").(string)
	if pointInTimeUTC != "" {
		v, err := time.Parse(time.RFC3339, pointInTimeUTC)
		if err != nil {
			return fmt.Errorf("unable to parse `point_in_time_restore_time_in_utc` value")
		}
		parameters.ServerProperties.RestorePointInTime = &date.Time{Time: v}
	}

	future, err := client.Create(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	// `maintenance_window` could only be updated with, could not be created with
	if v, ok := d.GetOk("maintenance_window"); ok {
		mwParams := mysqlflexibleservers.ServerForUpdate{
			ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{
				MaintenanceWindow: expandArmServerMaintenanceWindow(v.([]interface{})),
			},
		}
		mwFuture, err := client.Update(ctx, id.ResourceGroup, id.Name, mwParams)
		if err != nil {
			return fmt.Errorf("updating Mysql Flexible Server %q maintenance window (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}

		if err := mwFuture.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q maintenance window (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
		}
	}

	d.SetId(id.ID())

	return resourceMysqlFlexibleServerRead(d, meta)
}

func resourceMysqlFlexibleServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServerClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Mysql Flexible Server %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("zone", props.AvailabilityZone)
		d.Set("version", props.Version)
		d.Set("fqdn", props.FullyQualifiedDomainName)

		if network := props.Network; network != nil {
			d.Set("public_network_access_enabled", network.PublicNetworkAccess == mysqlflexibleservers.EnableStatusEnumEnabled)
			d.Set("delegated_subnet_id", network.DelegatedSubnetResourceID)
			d.Set("private_dns_zone_id", network.PrivateDNSZoneResourceID)
		}

		if err := d.Set("maintenance_window", flattenArmServerMaintenanceWindow(props.MaintenanceWindow)); err != nil {
			return fmt.Errorf("setting `maintenance_window`: %+v", err)
		}

		if err := d.Set("storage", flattenArmServerStorage(props.Storage)); err != nil {
			return fmt.Errorf("setting `storage`: %+v", err)
		}

		if backup := props.Backup; backup != nil {
			d.Set("backup_retention_days", backup.BackupRetentionDays)
			d.Set("geo_redundant_backup_enabled", backup.GeoRedundantBackup == mysqlflexibleservers.EnableStatusEnumEnabled)
		}

		if err := d.Set("high_availability", flattenFlexibleServerHighAvailability(props.HighAvailability)); err != nil {
			return fmt.Errorf("setting `high_availability`: %+v", err)
		}
		d.Set("replication_role", props.ReplicationRole)
		d.Set("replica_capacity", props.ReplicaCapacity)
	}

	sku, err := flattenFlexibleServerSku(resp.Sku)
	if err != nil {
		return fmt.Errorf("flattening `sku_name` for Mysql Flexible Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	d.Set("sku_name", sku)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceMysqlFlexibleServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServerClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	// failover is only supported when `zone` and `standby_availability_zone` is exchanged
	var requireFailover bool
	switch {
	case d.HasChange("zone") && d.HasChange("high_availability.0.standby_availability_zone"):
		resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return err
		}

		if props := resp.ServerProperties; props != nil {
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
	case d.HasChange("zone") && !d.HasChange("high_availability.0.standby_availability_zone"):
		return fmt.Errorf("`zone` cannot be changed independently")
	default:
		// No need failover when only `standby_availability_zone` is changed and both `zone` and `standby_availability_zone` aren't changed
		requireFailover = false
	}

	if d.HasChange("replication_role") {
		oldReplicationRole, newReplicationRole := d.GetChange("replication_role")
		if oldReplicationRole == "Replica" && newReplicationRole == "None" {
			parameters := mysqlflexibleservers.ServerForUpdate{
				ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{
					ReplicationRole: mysqlflexibleservers.ReplicationRoleNone,
				},
			}

			future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
			if err != nil {
				return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q) to update `replication_role`: %+v", id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q) to update `replication_role`: %+v", id.Name, id.ResourceGroup, err)
			}
		} else {
			return fmt.Errorf("`replication_role` only can be updated from `Replica` to `None`")
		}
	}

	// ha Enabled is dependent on storage auto grow Enabled. But when we enabled this two features in one request, it returns bad request.
	// Thus we need to separate these two updates in two requests.
	if d.HasChange("storage") && d.Get("storage.0.auto_grow_enabled").(bool) {
		parameters := mysqlflexibleservers.ServerForUpdate{
			ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{
				Storage: expandArmServerStorage(d.Get("storage").([]interface{})),
			},
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
		if err != nil {
			return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q) to enable `auto_grow_enabled`: %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q) to enable `auto_grow_enabled`: %+v", id.Name, id.ResourceGroup, err)
		}
	}

	if requireFailover {
		future, err := client.Failover(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			return fmt.Errorf("failing over %s: %+v", *id, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for failover of %s: %+v", *id, err)
		}
	} else if d.HasChange("high_availability") {
		parameters := mysqlflexibleservers.ServerForUpdate{
			ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{
				HighAvailability: &mysqlflexibleservers.HighAvailability{
					Mode: mysqlflexibleservers.HighAvailabilityModeDisabled,
				},
			},
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
		if err != nil {
			return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q) to disable `high_availability`: %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q) to disable `high_availability`: %+v", id.Name, id.ResourceGroup, err)
		}

		parameters.HighAvailability = expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{}))

		if parameters.HighAvailability.Mode != mysqlflexibleservers.HighAvailabilityModeDisabled {
			future, err = client.Update(ctx, id.ResourceGroup, id.Name, parameters)
			if err != nil {
				return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q) to update `high_availability`: %+v", id.Name, id.ResourceGroup, err)
			}

			if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
				return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q) to update `high_availability`: %+v", id.Name, id.ResourceGroup, err)
			}
		}
	}

	parameters := mysqlflexibleservers.ServerForUpdate{
		ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{},
	}

	if d.HasChange("administrator_password") {
		parameters.ServerPropertiesForUpdate.AdministratorLoginPassword = utils.String(d.Get("administrator_password").(string))
	}

	if d.HasChange("backup_retention_days") || d.HasChange("geo_redundant_backup_enabled") {
		parameters.ServerPropertiesForUpdate.Backup = expandArmServerBackup(d)
	}

	if d.HasChange("maintenance_window") {
		parameters.ServerPropertiesForUpdate.MaintenanceWindow = expandArmServerMaintenanceWindow(d.Get("maintenance_window").([]interface{}))
	}

	if d.HasChange("sku_name") {
		sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
		if err != nil {
			return fmt.Errorf("expanding `sku_name` for Mysql Flexible Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
		}
		parameters.Sku = sku
	}

	if d.HasChange("tags") {
		parameters.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if d.HasChange("storage") && !d.Get("storage.0.auto_grow_enabled").(bool) {
		parameters := mysqlflexibleservers.ServerForUpdate{
			ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{
				Storage: expandArmServerStorage(d.Get("storage").([]interface{})),
			},
		}

		future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
		if err != nil {
			return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q) to disable `auto_grow_enabled`: %+v", id.Name, id.ResourceGroup, err)
		}

		if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
			return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q) to disable `auto_grow_enabled`: %+v", id.Name, id.ResourceGroup, err)
		}
	}

	return resourceMysqlFlexibleServerRead(d, meta)
}

func resourceMysqlFlexibleServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.FlexibleServerClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.FlexibleServerID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of the Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	return nil
}

func expandArmServerNetwork(d *pluginsdk.ResourceData) *mysqlflexibleservers.Network {
	network := mysqlflexibleservers.Network{}

	if v, ok := d.GetOk("delegated_subnet_id"); ok {
		network.DelegatedSubnetResourceID = utils.String(v.(string))
	}

	if v, ok := d.GetOk("private_dns_zone_id"); ok {
		network.PrivateDNSZoneResourceID = utils.String(v.(string))
	}

	return &network
}

func expandArmServerMaintenanceWindow(input []interface{}) *mysqlflexibleservers.MaintenanceWindow {
	if len(input) == 0 {
		return &mysqlflexibleservers.MaintenanceWindow{
			CustomWindow: utils.String(ServerMaintenanceWindowDisabled),
		}
	}
	v := input[0].(map[string]interface{})

	maintenanceWindow := mysqlflexibleservers.MaintenanceWindow{
		CustomWindow: utils.String(ServerMaintenanceWindowEnabled),
		StartHour:    utils.Int32(int32(v["start_hour"].(int))),
		StartMinute:  utils.Int32(int32(v["start_minute"].(int))),
		DayOfWeek:    utils.Int32(int32(v["day_of_week"].(int))),
	}

	return &maintenanceWindow
}

func expandArmServerStorage(inputs []interface{}) *mysqlflexibleservers.Storage {
	if len(inputs) == 0 || inputs[0] == nil {
		return nil
	}

	input := inputs[0].(map[string]interface{})
	autoGrow := mysqlflexibleservers.EnableStatusEnumDisabled
	if v := input["auto_grow_enabled"].(bool); v {
		autoGrow = mysqlflexibleservers.EnableStatusEnumEnabled
	}

	storage := mysqlflexibleservers.Storage{
		AutoGrow: autoGrow,
	}

	if v := input["size_gb"].(int); v != 0 {
		storage.StorageSizeGB = utils.Int32(int32(v))
	}

	if v := input["iops"].(int); v != 0 {
		storage.Iops = utils.Int32(int32(v))
	}

	return &storage
}

func flattenArmServerStorage(storage *mysqlflexibleservers.Storage) []interface{} {
	if storage == nil {
		return []interface{}{}
	}

	var size, iops int32
	if storage.StorageSizeGB != nil {
		size = *storage.StorageSizeGB
	}

	if storage.Iops != nil {
		iops = *storage.Iops
	}

	return []interface{}{
		map[string]interface{}{
			"size_gb":           size,
			"iops":              iops,
			"auto_grow_enabled": storage.AutoGrow == mysqlflexibleservers.EnableStatusEnumEnabled,
		},
	}
}

func expandArmServerBackup(d *pluginsdk.ResourceData) *mysqlflexibleservers.Backup {
	geoRedundantBackup := mysqlflexibleservers.EnableStatusEnumDisabled
	if d.Get("geo_redundant_backup_enabled").(bool) {
		geoRedundantBackup = mysqlflexibleservers.EnableStatusEnumEnabled
	}

	backup := mysqlflexibleservers.Backup{
		GeoRedundantBackup: geoRedundantBackup,
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		backup.BackupRetentionDays = utils.Int32(int32(v.(int)))
	}

	return &backup
}

func expandFlexibleServerSku(name string) (*mysqlflexibleservers.Sku, error) {
	if name == "" {
		return nil, nil
	}
	parts := strings.SplitAfterN(name, "_", 2)

	var tier mysqlflexibleservers.SkuTier
	switch strings.TrimSuffix(parts[0], "_") {
	case "B":
		tier = mysqlflexibleservers.SkuTierBurstable
	case "GP":
		tier = mysqlflexibleservers.SkuTierGeneralPurpose
	case "MO":
		tier = mysqlflexibleservers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", name, parts[0])
	}

	return &mysqlflexibleservers.Sku{
		Name: utils.String(parts[1]),
		Tier: tier,
	}, nil
}

func flattenFlexibleServerSku(sku *mysqlflexibleservers.Sku) (string, error) {
	if sku == nil || sku.Name == nil || sku.Tier == "" {
		return "", nil
	}

	var tier string
	switch sku.Tier {
	case mysqlflexibleservers.SkuTierBurstable:
		tier = "B"
	case mysqlflexibleservers.SkuTierGeneralPurpose:
		tier = "GP"
	case mysqlflexibleservers.SkuTierMemoryOptimized:
		tier = "MO"
	default:
		return "", fmt.Errorf("sku_name has unknown sku tier %s", sku.Tier)
	}

	return strings.Join([]string{tier, *sku.Name}, "_"), nil
}

func flattenArmServerMaintenanceWindow(input *mysqlflexibleservers.MaintenanceWindow) []interface{} {
	if input == nil || input.CustomWindow == nil || *input.CustomWindow == string(ServerMaintenanceWindowDisabled) {
		return make([]interface{}, 0)
	}

	var dayOfWeek int32
	if input.DayOfWeek != nil {
		dayOfWeek = *input.DayOfWeek
	}
	var startHour int32
	if input.StartHour != nil {
		startHour = *input.StartHour
	}
	var startMinute int32
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

func expandFlexibleServerHighAvailability(inputs []interface{}) *mysqlflexibleservers.HighAvailability {
	if len(inputs) == 0 || inputs[0] == nil {
		return &mysqlflexibleservers.HighAvailability{
			Mode: mysqlflexibleservers.HighAvailabilityModeDisabled,
		}
	}

	input := inputs[0].(map[string]interface{})

	mode := mysqlflexibleservers.HighAvailabilityMode(input["mode"].(string))

	result := mysqlflexibleservers.HighAvailability{
		Mode: mode,
	}

	// for updating mode from ZoneRedundant to SameZone, the standby az will be changed
	// if we keep setting the standby az of ZoneRedundant, ha could not be changed to SameZone
	if mode == mysqlflexibleservers.HighAvailabilityModeSameZone {
		return &result
	}

	if v, ok := input["standby_availability_zone"]; ok && v.(string) != "" {
		result.StandbyAvailabilityZone = utils.String(v.(string))
	}

	return &result
}

func flattenFlexibleServerHighAvailability(ha *mysqlflexibleservers.HighAvailability) []interface{} {
	if ha == nil || ha.Mode == mysqlflexibleservers.HighAvailabilityModeDisabled {
		return []interface{}{}
	}

	var zone string
	if ha.StandbyAvailabilityZone != nil {
		zone = *ha.StandbyAvailabilityZone
	}

	return []interface{}{
		map[string]interface{}{
			"mode":                      string(ha.Mode),
			"standby_availability_zone": zone,
		},
	}
}
