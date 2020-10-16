package mssql

import (
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/sql/mgmt/v3.0/sql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMSSQLManagedInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMSSQLManagedInstanceCreateUpdate,
		Read:   resourceArmMSSQLManagedInstanceRead,
		Update: resourceArmMSSQLManagedInstanceCreateUpdate,
		Delete: resourceArmMSSQLManagedInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(600 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(1200 * time.Minute),
			Delete: schema.DefaultTimeout(600 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateManagedInstanceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SystemAssigned",
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
							Optional: true,
						},
					},
				},
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"collation": {
				Type:             schema.TypeString,
				DiffSuppressFunc: suppress.CaseDifference,
				Optional:         true,
				ForceNew:         true,
			},

			"dns_zone_partner": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"instance_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"license_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.LicenseIncluded),
					string(sql.BasePrice),
				}, false),
			},

			"maintenance_configuration_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ManagedServerCreateModeDefault),
					string(sql.ManagedServerCreateModePointInTimeRestore),
				}, false),
			},

			"minimal_tls_version": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"None",
					"1.0",
					"1.1",
					"1.2",
				}, false),
			},

			"proxy_override": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(sql.ManagedInstanceProxyOverrideDefault),
					string(sql.ManagedInstanceProxyOverrideProxy),
					string(sql.ManagedInstanceProxyOverrideRedirect),
				}, false),
			},

			"public_data_endpoint_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"restore_point_in_time": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppress.RFC3339Time,
				ValidateFunc:     validation.IsRFC3339Time,
				RequiredWith:     []string{"restore_point_in_time", "source_managed_instance_id"},
			},

			"source_managed_instance_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
				RequiredWith: []string{"restore_point_in_time", "source_managed_instance_id"},
			},

			"storage_size_gb": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(32, 16384),
			},

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"timezone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "UTC",
				ForceNew:     true,
				ValidateFunc: azure.ValidateManagedInstanceTimeZones(),
			},

			"vcores": {
				Type:     schema.TypeInt,
				Optional: true,
				ValidateFunc: validation.IntInSlice([]int{
					4,
					8,
					16,
					24,
					32,
					40,
					64,
					80,
				}),
			},

			"sku": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"capacity": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"family": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"GP_Gen4",
								"GP_Gen5",
								"BC_Gen4",
								"BC_Gen5",
							}, false),
						},

						"size": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"tier": {
							Type:             schema.TypeString,
							Optional:         true,
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"fully_qualified_domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_zone": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmMSSQLManagedInstanceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	createMode := d.Get("create_mode").(string)
	adminName := d.Get("administrator_login").(string)
	restorePoint := d.Get("restore_point_in_time").(string)
	sourceManagedInstanceID := d.Get("source_managed_instance_id").(string)
	location := d.Get("location").(string)
	t := d.Get("tags").(map[string]interface{})

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing managed sql instance %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_sql_managed_instance", *existing.ID)
		}
	}

	parameters := sql.ManagedInstance{
		Location: utils.String(location),
		Tags:     tags.Expand(t),
		ManagedInstanceProperties: &sql.ManagedInstanceProperties{
			AdministratorLogin: utils.String(adminName),
		},
	}

	if d.HasChange("create_mode") {
		if createMode == string(sql.ManagedServerCreateModePointInTimeRestore) && (len(restorePoint) == 0 || len(sourceManagedInstanceID) == 0) {
			return fmt.Errorf("could not configure managed SQL instance %q (Resource Group %q) in restore in point create mode. The restore point in time value must be supplied", name, resourceGroup)
		}
		parameters.ManagedInstanceProperties.ManagedInstanceCreateMode = sql.ManagedServerCreateMode(createMode)
	}

	if _, ok := d.GetOk("identity"); ok {
		sqlServerIdentity := expandManagedInstanceIdentity(d)
		parameters.Identity = sqlServerIdentity
	}

	if _, ok := d.GetOk("sku"); ok {
		sku := expandManagedInstanceSku(d)
		parameters.Sku = sku
	}

	if v, exists := d.GetOk("license_type"); exists {
		licenseType := v.(string)
		parameters.ManagedInstanceProperties.LicenseType = sql.ManagedInstanceLicenseType(licenseType)
	}

	if v, exists := d.GetOk("collation"); exists {
		collation := v.(string)
		parameters.ManagedInstanceProperties.Collation = utils.String(collation)
	}

	if v, exists := d.GetOk("storage_size_gb"); exists {
		storageSize := v.(int)
		if storageSize%32 != 0 {
			return fmt.Errorf("Could not create managed sql instance %q (Resource Group %q). The storage size in db should be in increments of 32", name, resourceGroup)
		}

		parameters.ManagedInstanceProperties.StorageSizeInGB = utils.Int32(int32(storageSize))
	}

	if v, exists := d.GetOk("vcores"); exists {
		vcores := v.(int)
		parameters.ManagedInstanceProperties.VCores = utils.Int32(int32(vcores))
	}

	if v, exists := d.GetOk("dns_zone_partner"); exists {
		dnsZonePartner := v.(string)
		parameters.ManagedInstanceProperties.DNSZonePartner = utils.String(dnsZonePartner)
	}

	if v, ok := d.GetOk("restore_point_in_time"); ok {
		if createMode != string(sql.ManagedServerCreateModePointInTimeRestore) {
			return fmt.Errorf("'restore_point_in_time' is supported only for create_mode %s", string(sql.ManagedServerCreateModePointInTimeRestore))
		}
		restorePointInTime := v.(string)
		restorePointInTimeDate, err2 := date.ParseTime(time.RFC3339, restorePointInTime)
		if err2 != nil {
			return fmt.Errorf("`restore_point_in_time` wasn't a valid RFC3339 date %q: %+v", restorePointInTime, err2)
		}

		parameters.ManagedInstanceProperties.RestorePointInTime = &date.Time{
			Time: restorePointInTimeDate,
		}
	}

	if v, ok := d.GetOk("source_managed_instance_id"); ok {
		if createMode != string(sql.ManagedServerCreateModePointInTimeRestore) {
			return fmt.Errorf("'source_managed_instance_id' is supported only for create_mode %s", string(sql.ManagedServerCreateModePointInTimeRestore))
		}
		sourceManagedInstance := v.(string)
		parameters.ManagedInstanceProperties.SourceManagedInstanceID = utils.String(sourceManagedInstance)
	}

	if d.HasChange("administrator_login_password") {
		adminPassword := d.Get("administrator_login_password").(string)
		parameters.ManagedInstanceProperties.AdministratorLoginPassword = utils.String(adminPassword)
	}

	if v, ok := d.GetOk("public_data_endpoint_enabled"); ok {
		publicDataEndpointEnabled := v.(bool)
		parameters.ManagedInstanceProperties.PublicDataEndpointEnabled = utils.Bool(publicDataEndpointEnabled)
	}

	if v, ok := d.GetOk("subnet_id"); ok {
		subnetID := v.(string)
		parameters.ManagedInstanceProperties.SubnetID = utils.String(subnetID)
	}

	if v, ok := d.GetOk("proxy_override"); ok {
		proxyOverride := v.(string)
		parameters.ManagedInstanceProperties.ProxyOverride = sql.ManagedInstanceProxyOverride(proxyOverride)
	}

	if v, ok := d.GetOk("timezone_id"); ok {
		timezoneID := v.(string)
		parameters.ManagedInstanceProperties.TimezoneID = utils.String(timezoneID)
	}

	if v, ok := d.GetOk("instance_pool_id"); ok {
		instancePoolID := v.(string)
		parameters.ManagedInstanceProperties.InstancePoolID = utils.String(instancePoolID)
	}

	if v, ok := d.GetOk("maintenance_configuration_id"); ok {
		maintenanceConfigurationID := v.(string)
		parameters.ManagedInstanceProperties.MaintenanceConfigurationID = utils.String(maintenanceConfigurationID)
	}

	if v, ok := d.GetOk("minimal_tls_version"); ok {
		minimalTLSVersion := v.(string)
		parameters.ManagedInstanceProperties.MinimalTLSVersion = utils.String(minimalTLSVersion)
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error issuing create/update request for Managed Instance %q (Resource group: %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting on create/update future for managed SQL instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	result, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error making get request for managed SQL instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if result.ID == nil {
		return fmt.Errorf("Error getting ID from managed SQL instance %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*result.ID)

	return resourceArmMSSQLManagedInstanceRead(d, meta)
}

func resourceArmMSSQLManagedInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	resp, err := client.Get(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error reading managed SQL instance %s: %v", name, err)
	}

	d.Set("name", name)
	d.Set("resource_group_name", resGroup)
	d.Set("type", (resp.Type))

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if err := d.Set("identity", flattenManagedInstanceIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if props := resp.ManagedInstanceProperties; props != nil {
		d.Set("fully_qualified_domain_name", props.FullyQualifiedDomainName)
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("subnet_id", props.SubnetID)
		d.Set("state", props.State)
		d.Set("license_type", props.LicenseType)
		d.Set("vcores", props.VCores)
		d.Set("storage_size_gb", props.StorageSizeInGB)
		d.Set("collation", props.Collation)
		d.Set("dns_zone", props.DNSZone)
		d.Set("public_data_endpoint_enabled", props.PublicDataEndpointEnabled)
		d.Set("proxy_override", props.ProxyOverride)
		d.Set("timezone_id", props.TimezoneID)
		d.Set("instance_pool_id", props.InstancePoolID)
		d.Set("maintenance_configuration_id", props.MaintenanceConfigurationID)
		d.Set("minimal_tls_version", props.MinimalTLSVersion)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMSSQLManagedInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MSSQL.ManagedInstancesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resGroup := id.ResourceGroup
	name := id.Path["managedInstances"]

	future, err := client.Delete(ctx, resGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting managed SQL instance %s: %+v", name, err)
	}

	return future.WaitForCompletionRef(ctx, client.Client)
}

func expandManagedInstanceIdentity(d *schema.ResourceData) *sql.ResourceIdentity {
	identities := d.Get("identity").([]interface{})
	if len(identities) == 0 {
		return &sql.ResourceIdentity{}
	}
	identity := identities[0].(map[string]interface{})
	identityType := sql.IdentityType(identity["type"].(string))
	return &sql.ResourceIdentity{
		Type: identityType,
	}
}

func expandManagedInstanceSku(d *schema.ResourceData) *sql.Sku {
	skus := d.Get("sku").([]interface{})
	if len(skus) == 0 {
		return &sql.Sku{}
	}
	sku := skus[0].(map[string]interface{})
	skuName := sku["name"].(string)
	skuCapacity := sku["capacity"].(int)
	skuFamily := sku["family"].(string)
	skuSize := sku["size"].(string)
	skuTier := sku["tier"].(string)

	managedInstanceSku := sql.Sku{
		Name: utils.String(skuName),
	}
	if skuCapacity > 0 {
		managedInstanceSku.Capacity = utils.Int32(int32(skuCapacity))
	}

	if len(strings.TrimSpace(skuFamily)) != 0 {
		managedInstanceSku.Family = utils.String(skuFamily)
	}

	if len(strings.TrimSpace(skuSize)) != 0 {
		managedInstanceSku.Size = utils.String(skuSize)
	}

	if len(strings.TrimSpace(skuTier)) != 0 {
		managedInstanceSku.Tier = utils.String(skuTier)
	}

	return &managedInstanceSku
}

func flattenManagedInstanceIdentity(identity *sql.ResourceIdentity) []interface{} {
	if identity == nil {
		return []interface{}{}
	}
	result := make(map[string]interface{})
	result["type"] = identity.Type
	if identity.PrincipalID != nil {
		result["principal_id"] = identity.PrincipalID.String()
	}
	if identity.TenantID != nil {
		result["tenant_id"] = identity.TenantID.String()
	}

	return []interface{}{result}
}
