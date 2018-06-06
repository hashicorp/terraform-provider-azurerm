package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPostgreSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerCreate,
		Read:   resourceArmPostgreSQLServerRead,
		Update: resourceArmPostgreSQLServerUpdate,
		Delete: resourceArmPostgreSQLServerDelete,
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
								"MO_Gen5_2",
								"MO_Gen5_4",
								"MO_Gen5_8",
								"MO_Gen5_16",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"capacity": {
							Type:     schema.TypeInt,
							Required: true,
							ValidateFunc: validateIntInSlice([]int{
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
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.Basic),
								string(postgresql.GeneralPurpose),
								string(postgresql.MemoryOptimized),
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},

						"family": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"Gen4",
								"Gen5",
							}, true),
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
					string(postgresql.NineFullStopFive),
					string(postgresql.NineFullStopSix),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
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
							ForceNew:     true,
							ValidateFunc: validateIntBetweenDivisibleBy(5120, 1048576, 1024),
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
							DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
						},
					},
				},
			},

			"ssl_enforcement": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.SslEnforcementEnumDisabled),
					string(postgresql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: ignoreCaseDiffSuppressFunc,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tagsSchema(),
		},
	}
}

func resourceArmPostgreSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	name := d.Get("name").(string)
	location := azureRMNormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	adminLogin := d.Get("administrator_login").(string)
	adminLoginPassword := d.Get("administrator_login_password").(string)
	sslEnforcement := d.Get("ssl_enforcement").(string)
	version := d.Get("version").(string)
	createMode := "Default"
	tags := d.Get("tags").(map[string]interface{})

	sku := expandAzureRmPostgreSQLServerSku(d)
	storageProfile := expandAzureRmPostgreSQLStorageProfile(d)

	properties := postgresql.ServerForCreate{
		Location: &location,
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(adminLogin),
			AdministratorLoginPassword: utils.String(adminLoginPassword),
			Version:                    postgresql.ServerVersion(version),
			SslEnforcement:             postgresql.SslEnforcementEnum(sslEnforcement),
			StorageProfile:             storageProfile,
			CreateMode:                 postgresql.CreateMode(createMode),
		},
		Sku:  sku,
		Tags: expandTags(tags),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for creation of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient
	ctx := meta.(*ArmClient).StopContext

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	adminLoginPassword := d.Get("administrator_login_password").(string)
	sslEnforcement := d.Get("ssl_enforcement").(string)
	version := d.Get("version").(string)
	sku := expandAzureRmPostgreSQLServerSku(d)
	storageProfile := expandAzureRmPostgreSQLStorageProfile(d)
	tags := d.Get("tags").(map[string]interface{})

	properties := postgresql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			StorageProfile:             storageProfile,
			AdministratorLoginPassword: utils.String(adminLoginPassword),
			Version:                    postgresql.ServerVersion(version),
			SslEnforcement:             postgresql.SslEnforcementEnum(sslEnforcement),
		},
		Sku:  sku,
		Tags: expandTags(tags),
	}

	future, err := client.Update(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error updating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("Error waiting for update of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("Error retrieving PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %s (resource group %s) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient
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
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))

	if err := d.Set("sku", flattenPostgreSQLServerSku(resp.Sku)); err != nil {
		return fmt.Errorf("Error flattening `sku`: %+v", err)
	}

	if err := d.Set("storage_profile", flattenPostgreSQLStorageProfile(resp.StorageProfile)); err != nil {
		return fmt.Errorf("Error flattening `storage_profile`: %+v", err)
	}

	flattenAndSetTags(d, resp.Tags)

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

	return nil
}

func resourceArmPostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).postgresqlServersClient
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

		return fmt.Errorf("Error deleting PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	err = future.WaitForCompletion(ctx, client.Client)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandAzureRmPostgreSQLServerSku(d *schema.ResourceData) *postgresql.Sku {
	skus := d.Get("sku").([]interface{})
	sku := skus[0].(map[string]interface{})

	name := sku["name"].(string)
	capacity := sku["capacity"].(int)
	tier := sku["tier"].(string)
	family := sku["family"].(string)

	return &postgresql.Sku{
		Name:     utils.String(name),
		Tier:     postgresql.SkuTier(tier),
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(family),
	}
}

func expandAzureRmPostgreSQLStorageProfile(d *schema.ResourceData) *postgresql.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)

	return &postgresql.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  postgresql.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
	}
}

func flattenPostgreSQLServerSku(resp *postgresql.Sku) []interface{} {
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

func flattenPostgreSQLStorageProfile(resp *postgresql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = postgresql.GeoRedundantBackup(resp.GeoRedundantBackup)

	return []interface{}{values}
}
