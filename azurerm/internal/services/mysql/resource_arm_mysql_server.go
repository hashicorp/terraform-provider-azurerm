package mysql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMySqlServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMySqlServerCreate,
		Read:   resourceArmMySqlServerRead,
		Update: resourceArmMySqlServerUpdate,
		Delete: resourceArmMySqlServerDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateMySqlServerName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:          schema.TypeString,
				Optional:      true, // required in 2.0
				Computed:      true, // remove in 2.0
				ConflictsWith: []string{"sku"},
				ValidateFunc: validation.StringInSlice([]string{
					"B_Gen4_1",
					"B_Gen4_2",
					"B_Gen5_1",
					"B_Gen5_2",
					"GP_Gen4_2",
					"GP_Gen4_4",
					"GP_Gen4_8",
					"GP_Gen4_16",
					"GP_Gen4_32",
					"GP_Gen5_2",
					"GP_Gen5_4",
					"GP_Gen5_8",
					"GP_Gen5_16",
					"GP_Gen5_32",
					"GP_Gen5_64",
					"MO_Gen5_2",
					"MO_Gen5_4",
					"MO_Gen5_8",
					"MO_Gen5_16",
					"MO_Gen5_32",
				}, false),
			},

			// remove in 2.0
			"sku": {
				Type:          schema.TypeList,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"sku_name"},
				Deprecated:    "This property has been deprecated in favour of the 'sku_name' property and will be removed in version 2.0 of the provider",
				MaxItems:      1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"B_Gen4_1",
								"B_Gen4_2",
								"B_Gen5_1",
								"B_Gen5_2",
								"GP_Gen4_2",
								"GP_Gen4_4",
								"GP_Gen4_8",
								"GP_Gen4_16",
								"GP_Gen4_32",
								"GP_Gen5_2",
								"GP_Gen5_4",
								"GP_Gen5_8",
								"GP_Gen5_16",
								"GP_Gen5_32",
								"GP_Gen5_64",
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
								"MO_Gen5_32",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.IntInSlice([]int{
								1,
								2,
								4,
								8,
								16,
								32,
								64,
							}),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(mysql.Basic),
								string(mysql.GeneralPurpose),
								string(mysql.MemoryOptimized),
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
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

			"version": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.FiveFullStopSix),
					string(mysql.FiveFullStopSeven),
					string(mysql.EightFullStopZero),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
				ForceNew:         true,
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validation.All(
								validation.IntBetween(5120, 4194304),
								validation.IntDivisibleBy(1024),
							),
						},
						"backup_retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntBetween(7, 35),
						},
						"geo_redundant_backup": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
						"auto_grow": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  string(mysql.StorageAutogrowEnabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(mysql.StorageAutogrowEnabled),
								string(mysql.StorageAutogrowDisabled),
							}, false),
						},
					},
				},
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.SslEnforcementEnumDisabled),
					string(mysql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			tier, _ := diff.GetOk("sku.0.tier")
			storageMB, _ := diff.GetOk("storage_profile.0.storage_mb")

			if strings.ToLower(tier.(string)) == "basic" && storageMB.(int) > 1048576 {
				return fmt.Errorf("basic pricing tier only supports upto 1,048,576 MB (1TB) of storage")
			}

			return nil
		},
	}
}

func resourceArmMySqlServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Server creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mysql_server", *existing.ID)
		}
	}

	var sku *mysql.Sku
	if b, ok := d.GetOk("sku_name"); ok {
		var err error
		sku, err = expandServerSkuName(b.(string))
		if err != nil {
			return fmt.Errorf("error expanding sku_name for MySQL Server %q (Resource Group %q): %v", name, resourceGroup, err)
		}
	} else if _, ok := d.GetOk("sku"); ok {
		sku = expandMySQLServerSku(d)
	} else {
		return fmt.Errorf("One of `sku` or `sku_name` must be set for MySQL Server %q (Resource Group %q)", name, resourceGroup)
	}

	properties := mysql.ServerForCreate{
		Location: &location,
		Properties: &mysql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(d.Get("administrator_login").(string)),
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			Version:                    mysql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             mysql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
			StorageProfile:             expandMySQLStorageProfile(d),
			CreateMode:                 mysql.CreateMode("Default"),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read MySQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySqlServerRead(d, meta)
}

func resourceArmMySqlServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Server update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	var sku *mysql.Sku
	if b, ok := d.GetOk("sku_name"); ok {
		var err error
		sku, err = expandServerSkuName(b.(string))
		if err != nil {
			return fmt.Errorf("error expanding sku_name for MySQL Server %q (Resource Group %q): %v", name, resourceGroup, err)
		}
	} else if _, ok := d.GetOk("sku"); ok {
		sku = expandMySQLServerSku(d)
	} else {
		return fmt.Errorf("One of `sku` or `sku_name` must be set for MySQL Server %q (Resource Group %q)", name, resourceGroup)
	}

	properties := mysql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &mysql.ServerUpdateParametersProperties{
			StorageProfile:             expandMySQLStorageProfile(d),
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			Version:                    mysql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             mysql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error updating MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for MySQL Server %q (Resource Group %q) to finish updating: %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read MySQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySqlServerRead(d, meta)
}

func resourceArmMySqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))

	if err := d.Set("sku", flattenMySQLServerSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if err := d.Set("storage_profile", flattenMySQLStorageProfile(resp.StorageProfile)); err != nil {
		return fmt.Errorf("Error setting `storage_profile`: %+v", err)
	}

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMySqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error deleting MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for deletion of MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandServerSkuName(skuName string) (*mysql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the worng numberof parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier mysql.SkuTier
	switch parts[0] {
	case "B":
		tier = mysql.Basic
	case "GP":
		tier = mysql.GeneralPurpose
	case "MO":
		tier = mysql.MemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &mysql.Sku{
		Name:     utils.String(skuName),
		Tier:     tier,
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandMySQLServerSku(d *schema.ResourceData) *mysql.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)
	family := sku["family"].(string)

	return &mysql.Sku{
		Name:     utils.String(name),
		Tier:     mysql.SkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(family),
	}
}

func expandMySQLStorageProfile(d *schema.ResourceData) *mysql.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)
	autoGrow := storageprofile["auto_grow"].(string)

	return &mysql.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  mysql.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
		StorageAutogrow:     mysql.StorageAutogrow(autoGrow),
	}
}

func flattenMySQLServerSku(resp *mysql.Sku) []interface{} {
	values := map[string]interface{}{}

	if name := resp.Name; name != nil {
		values["name"] = *name
	}

	if capacity := resp.Capacity; capacity != nil {
		values["capacity"] = *capacity
	}

	values["tier"] = string(resp.Tier)

	if family := resp.Family; family != nil {
		values["family"] = *family
	}

	return []interface{}{values}
}

func flattenMySQLStorageProfile(resp *mysql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	values["auto_grow"] = string(resp.StorageAutogrow)

	return []interface{}{values}
}
