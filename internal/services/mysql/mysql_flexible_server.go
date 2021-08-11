package mysql

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/services/preview/mysql/mgmt/2021-05-01-preview/mysqlflexibleservers"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	privateDnsValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/privatedns/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
	"log"
	"strings"
	"time"
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": azure.SchemaLocation(),

			"administrator_login": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"administrator_password": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Sensitive:    true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validate.FlexibleServerSkuName,
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

			"zone": azure.SchemaZoneComputed(),

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

			"private_dns_zone_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: privateDnsValidate.PrivateDnsZoneID,
			},

			"point_in_time_restore_time_in_utc": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"source_server_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validate.FlexibleServerID,
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

			"backup_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(7, 35),
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
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
							}, false),
						},
						"standby_availability_zone": azure.SchemaZoneComputed(),
					},
				},
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(mysqlflexibleservers.ResourceIdentityTypeSystemAssigned),
							}, false),
						},

						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"replication_role": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
			},

			"storage": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"size_gb": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntBetween(20, 16384),
						},

						"iops": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							Computed:     true,
							ValidateFunc: validation.IntInSlice([]int{400, 640, 1280, 3200, 6400, 12800, 20000}),
						},

						"auto_grow_enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},

			"cmk_enabled": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

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

	createMode := d.Get("create_mode").(string)

	if mysqlflexibleservers.CreateMode(createMode) == mysqlflexibleservers.CreateModePointInTimeRestore {
		if _, ok := d.GetOk("source_server_id"); !ok {
			return fmt.Errorf("`source_server_id` is required when `create_mode` is `PointInTimeRestore`")
		}
		if _, ok := d.GetOk("point_in_time_restore_time_in_utc"); !ok {
			return fmt.Errorf("`point_in_time_restore_time_in_utc` is required when `create_mode` is `PointInTimeRestore`")
		}
	}

	if createMode == "" || mysqlflexibleservers.CreateMode(createMode) == mysqlflexibleservers.CreateModeDefault {
		if _, ok := d.GetOk("administrator_login"); !ok {
			return fmt.Errorf("`administrator_login` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("administrator_password"); !ok {
			return fmt.Errorf("`administrator_password` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("sku_name"); !ok {
			return fmt.Errorf("`sku_name` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("version"); !ok {
			return fmt.Errorf("`version` is required when `create_mode` is `Default`")
		}
		if _, ok := d.GetOk("storage"); !ok {
			return fmt.Errorf("`storage` is required when `create_mode` is `Default`")
		}
	}

	sku, err := expandFlexibleServerSku(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name` for MySql Flexible Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	parameters := mysqlflexibleservers.Server{
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		ServerProperties: &mysqlflexibleservers.ServerProperties{
			CreateMode:       mysqlflexibleservers.CreateMode(d.Get("create_mode").(string)),
			Network:          expandArmServerNetwork(d),
			Version:          mysqlflexibleservers.ServerVersion(d.Get("version").(string)),
			Storage:          expandArmServerStorage(d),
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

	if err := d.Set("identity", flattenArmServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

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

		if storage := props.Storage; storage != nil && storage.StorageSizeGB != nil {
			d.Set("storage_mb", (*storage.StorageSizeGB * 1024))
		}

		if backup := props.Backup; backup != nil {
			d.Set("backup_retention_days", backup.BackupRetentionDays)
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

	parameters := mysqlflexibleservers.ServerForUpdate{
		Location:                  utils.String(location.Normalize(d.Get("location").(string))),
		ServerPropertiesForUpdate: &mysqlflexibleservers.ServerPropertiesForUpdate{},
	}

	if d.HasChange("administrator_password") {
		parameters.ServerPropertiesForUpdate.AdministratorLoginPassword = utils.String(d.Get("administrator_password").(string))
	}

	if d.HasChange("storage_mb") {
		parameters.ServerPropertiesForUpdate.Storage = expandArmServerStorage(d)
	}

	if d.HasChange("backup_retention_days") {
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

	if d.HasChange("high_availability") {
		parameters.HighAvailability = expandFlexibleServerHighAvailability(d.Get("high_availability").([]interface{}))
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("updating Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the update of the Mysql Flexible Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
		network.PrivateDNSZoneArmResourceID = utils.String(v.(string))
	}

	return &network
}

func expandArmServerIdentity(input []interface{}) *mysqlflexibleservers.Identity {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})
	return &mysqlflexibleservers.Identity{
		Type: mysqlflexibleservers.ResourceIdentityType(v["type"].(string)),
	}
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

func expandArmServerStorage(d *pluginsdk.ResourceData) *mysqlflexibleservers.Storage {
	storage := mysqlflexibleservers.Storage{}

	if v, ok := d.GetOk("size_gb"); ok {
		storage.StorageSizeGB = utils.Int32(int32(v.(int)))
	}

	return &storage
}

func expandArmServerBackup(d *pluginsdk.ResourceData) *mysqlflexibleservers.Backup {
	backup := mysqlflexibleservers.Backup{}

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

func flattenArmServerIdentity(input *mysqlflexibleservers.Identity) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	var principalId string
	if input.PrincipalID != nil {
		principalId = *input.PrincipalID
	}
	var tenantId string
	if input.TenantID != nil {
		tenantId = *input.TenantID
	}
	return []interface{}{
		map[string]interface{}{
			"type":         string(input.Type),
			"principal_id": principalId,
			"tenant_id":    tenantId,
		},
	}
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

	result := mysqlflexibleservers.HighAvailability{
		Mode: mysqlflexibleservers.HighAvailabilityMode(input["mode"].(string)),
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
