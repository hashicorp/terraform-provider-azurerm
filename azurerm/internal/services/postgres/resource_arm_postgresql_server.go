package postgres

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/postgresql/mgmt/2017-12-01/postgresql"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ValidatePSQLServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z]{2}[-0-9a-z]{0,60}[0-9a-z]$`); !m {
		errors = append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-' or have more than 63 characters.", k))
	}

	return nil, errors
}

func resourceArmPostgreSQLServer() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPostgreSQLServerCreate,
		Read:   resourceArmPostgreSQLServerRead,
		Update: resourceArmPostgreSQLServerUpdate,
		Delete: resourceArmPostgreSQLServerDelete,
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
				ValidateFunc: ValidatePSQLServerName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
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
				}, false),
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
					string(postgresql.OneOne),
					string(postgresql.OneZero),
					string(postgresql.OneZeroFullStopZero),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
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
							Default:  string(postgresql.StorageAutogrowEnabled),
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.StorageAutogrowEnabled),
								string(postgresql.StorageAutogrowDisabled),
							}, false),
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
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPostgreSQLServerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server creation.")

	name := d.Get("name").(string)
	location := azure.NormalizeLocation(d.Get("location").(string))
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() {
		existing, err := client.Get(ctx, resourceGroup, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_postgresql_server", *existing.ID)
		}
	}

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	properties := postgresql.ServerForCreate{
		Location: &location,
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(d.Get("administrator_login").(string)),
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             postgresql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
			StorageProfile:             expandAzureRmPostgreSQLStorageProfile(d),
			CreateMode:                 postgresql.CreateMode("Default"),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error creating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM PostgreSQL Server update.")

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", name, resourceGroup, err)
	}

	properties := postgresql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
			SslEnforcement:             postgresql.SslEnforcementEnum(d.Get("ssl_enforcement").(string)),
			StorageProfile:             expandAzureRmPostgreSQLStorageProfile(d),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, resourceGroup, name, properties)
	if err != nil {
		return fmt.Errorf("Error updating PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
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
	client := meta.(*clients.Client).Postgres.ServersClient
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
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", name, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	d.Set("administrator_login", resp.AdministratorLogin)
	d.Set("version", string(resp.Version))
	d.Set("ssl_enforcement", string(resp.SslEnforcement))

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if err := d.Set("storage_profile", flattenPostgreSQLStorageProfile(resp.StorageProfile)); err != nil {
		return fmt.Errorf("Error setting `storage_profile`: %+v", err)
	}

	// Computed
	d.Set("fqdn", resp.FullyQualifiedDomainName)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
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
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of PostgreSQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	return nil
}

func expandServerSkuName(skuName string) (*postgresql.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier postgresql.SkuTier
	switch parts[0] {
	case "B":
		tier = postgresql.Basic
	case "GP":
		tier = postgresql.GeneralPurpose
	case "MO":
		tier = postgresql.MemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &postgresql.Sku{
		Name:     utils.String(skuName),
		Tier:     tier,
		Capacity: utils.Int32(int32(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandAzureRmPostgreSQLStorageProfile(d *schema.ResourceData) *postgresql.StorageProfile {
	storageprofiles := d.Get("storage_profile").([]interface{})
	storageprofile := storageprofiles[0].(map[string]interface{})

	backupRetentionDays := storageprofile["backup_retention_days"].(int)
	geoRedundantBackup := storageprofile["geo_redundant_backup"].(string)
	storageMB := storageprofile["storage_mb"].(int)
	autoGrow := storageprofile["auto_grow"].(string)

	return &postgresql.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  postgresql.GeoRedundantBackup(geoRedundantBackup),
		StorageMB:           utils.Int32(int32(storageMB)),
		StorageAutogrow:     postgresql.StorageAutogrow(autoGrow),
	}
}

func flattenPostgreSQLStorageProfile(resp *postgresql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	values["auto_grow"] = string(resp.StorageAutogrow)

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	return []interface{}{values}
}
