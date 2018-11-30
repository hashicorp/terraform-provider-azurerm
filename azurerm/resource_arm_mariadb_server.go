package azurerm

import (
	"fmt"
	"log"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"

	"github.com/Azure/azure-sdk-for-go/services/preview/mariadb/mgmt/2018-06-01-preview/mariadb"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmMariaDbServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmMariaDbServerCreateOrUpdate,
		Read:   resourceArmMariaDbServerRead,
		Update: resourceArmMariaDbServerCreateOrUpdate,
		Delete: resourceArmMariaDbServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"location": locationSchema(),

			"resource_group_name": resourceGroupNameSchema(),

			"sku": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"B_Gen5_1",
								"B_Gen5_2",
								"GP_Gen5_2",
								"GP_Gen5_4",
								"GP_Gen5_8",
								"GP_Gen5_16",
								"GP_Gen5_32",
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
							}, false),
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validate.IntInSlice([]int{
								1,
								2,
								4,
								8,
								16,
								32,
							}),
						},

						"tier": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(mariadb.Basic),
								string(mariadb.GeneralPurpose),
								string(mariadb.MemoryOptimized),
							}, true),
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen5",
							}, true),
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
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10.2",
				}, true),
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validate.IntBetweenAndDivisibleBy(5120, 4096000, 1024),
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
								string(mariadb.Enabled),
								string(mariadb.Disabled),
							}, true),
						},
					},
				},
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(mariadb.SslEnforcementEnumDisabled),
					string(mariadb.SslEnforcementEnumEnabled),
				}, true),
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmMariaDbServerCreateOrUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbServersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM MariaDB Server creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	adminLogin := d.Get("administrator_login").(string)
	adminLoginPassword := d.Get("administrator_login_password").(string)
	sslEnforcement := d.Get("ssl_enforcement").(string)
	version := d.Get("version").(string)
	createMode := "Default"
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmMariaDbServerSku(d)
	storageProfile := expandAzureRmMariaDbStorageProfile(d)

	capacity := sku.Capacity
	tier := sku.Tier
	storageMB := storageProfile.StorageMB

	if strings.ToLower(string(tier)) == "basic" {
		if *storageMB > 1024000 {
			return fmt.Errorf("basic pricing tier only supports upto 1,024,000 MB (1TB) of storage")
		}

		if *capacity > 2 {
			return fmt.Errorf("basic pricing tier only supports upto 2 vCores")
		}
	}

	if strings.ToLower(string(tier)) == "generalpurpose" && *capacity < 2 {
		return fmt.Errorf("general purpose pricing tier must have at least 2 vCores")
	}

	if strings.ToLower(string(tier)) == "memoryoptimized" {

		if *capacity < 2 {
			return fmt.Errorf("memory optimized pricing tier must have at least 2 vCores")
		}

		if *capacity > 16 {
			return fmt.Errorf("memory optimized pricing tier only supports upto 16 vCores")
		}
	}

	properties := mariadb.ServerForCreate{
		Location: &location,
		Properties: &mariadb.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(adminLogin),
			AdministratorLoginPassword: utils.String(adminLoginPassword),
			Version:                    mariadb.ServerVersion(version),
			SslEnforcement:             mariadb.SslEnforcementEnum(sslEnforcement),
			StorageProfile:             storageProfile,
			CreateMode:                 mariadb.CreateMode(createMode),
		},
		Sku:  sku,
		Tags: expandTags(tags),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read MariaDB Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMariaDbServerRead(d, meta)
}

func resourceArmMariaDbServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MariaDB Server %q was not found (Resource Group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))

	if err := d.Set("sku", flattenMariaDbServerSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error setting `sku`: %+v", err)
	}

	if err := d.Set("storage_profile", flattenMariaDbStorageProfile(resp.StorageProfile)); err != nil {
		return fmt.Errorf("Error setting `storage_profile`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

	return nil
}

func resourceArmMariaDbServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).mariadbServersClient
	ctx := meta.(*ArmClient).StopContext

	id, err := parseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["servers"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of MariaDB Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandAzureRmMariaDbServerSku(d *schema.ResourceData) *mariadb.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)
	family := sku["family"].(string)

	return &mariadb.Sku{
		Name:     utils.String(name),
		Tier:     mariadb.SkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(family),
	}
}

func expandAzureRmMariaDbStorageProfile(d *schema.ResourceData) *mariadb.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)

	return &mariadb.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  mariadb.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
	}
}

func flattenMariaDbServerSku(resp *mariadb.Sku) []interface{} {
	values := map[string]interface{}{}

	if resp == nil {
		return []interface{}{}
	}

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

func flattenMariaDbStorageProfile(resp *mariadb.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if resp == nil {
		return []interface{}{}
	}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	return []interface{}{values}
}
