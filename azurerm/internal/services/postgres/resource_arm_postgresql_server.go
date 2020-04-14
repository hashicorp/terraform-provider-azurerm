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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/postgres/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func ValidatePSQLServerName(i interface{}, k string) (_ []string, errors []error) {
	if m, regexErrs := validate.RegExHelper(i, k, `^[0-9a-z][-0-9a-z]{1,61}[0-9a-z]$`); !m {
		return nil, append(regexErrs, fmt.Errorf("%q can contain only lowercase letters, numbers, and '-', but can't start or end with '-'. And must be at least 3 characters and at most 63 characters", k))
	}

	return nil, nil
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
				DiffSuppressFunc: suppress.CaseDifference, // make case sensitive in 3.0
			},

			"storage_profile": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"storage_mb": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
							//ExactlyOneOf: []string{"storage_profile.0.storage_mb", "storage_profile.0.auto_grow_enabled", "storage_profile.0.auto_grow"},
							ValidateFunc: validation.All(
								validation.IntBetween(5120, 4194304),
								validation.IntDivisibleBy(1024),
							),
						},

						"auto_grow_enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true, // remove in 3.0 and default to true
							ConflictsWith: []string{"storage_profile.0.auto_grow"},
						},

						"auto_grow": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_profile.0.auto_grow_enabled"},
							Deprecated:    "this has been renamed to the boolean `auto_grow_enabled` and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.StringInSlice([]string{
								string(postgresql.StorageAutogrowEnabled),
								string(postgresql.StorageAutogrowDisabled),
							}, false),
						},

						"backup_retention_days": {
							Type:         schema.TypeInt,
							Optional:     true,
							Default:      7,
							ValidateFunc: validation.IntBetween(7, 35),
						},

						"geo_redundant_backup_enabled": {
							Type:          schema.TypeBool,
							Optional:      true,
							Computed:      true, // remove in 2.0 and default to false
							ConflictsWith: []string{"storage_profile.0.geo_redundant_backup"},
						},

						"geo_redundant_backup": {
							Type:          schema.TypeString,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"storage_profile.0.geo_redundant_backup_enabled"},
							Deprecated:    "this has been renamed to the boolean `geo_redundant_backup` and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
							DiffSuppressFunc: suppress.CaseDifference,
						},
					},
				},
			},

			"create_mode": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          string(postgresql.CreateModeDefault),
				DiffSuppressFunc: suppress.CaseDifference,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.CreateModeDefault),
					string(postgresql.CreateModeGeoRestore),
					string(postgresql.CreateModePointInTimeRestore),
					string(postgresql.CreateModeReplica),
					string(postgresql.CreateModeServerPropertiesForCreate),
				}, false),
			},

			"infrastructure_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(postgresql.TLSEnforcementDisabled),
					string(postgresql.TLS10),
					string(postgresql.TLS11),
					string(postgresql.TLS12),
				}, false),
			},

			"ssl_enforcement_enabled": {
				Type:         schema.TypeBool,
				Optional:     true, // required in 3.0
				Computed:     true, // remove computed in 3.0
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
			},

			"ssl_enforcement": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "this has been renamed to the boolean `ssl_enforcement_enabled` and will be removed in version 3.0 of the provider.",
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
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

	infraEncrypt := postgresql.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = postgresql.InfrastructureEncryptionDisabled
	}

	publicAccess := postgresql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = postgresql.PublicNetworkAccessEnumDisabled
	}

	ssl := postgresql.SslEnforcementEnumEnabled
	if v, ok := d.GetOk("ssl_enforcement"); ok && strings.EqualFold(v.(string), string(postgresql.SslEnforcementEnumDisabled)) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}
	if v, ok := d.GetOkExists("ssl_enforcement_enabled"); ok && !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}

	storage := expandAzureRmPostgreSQLStorageProfile(d)

	props := postgresql.ServerForCreate{
		Location: &location,
		Properties: &postgresql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         utils.String(d.Get("administrator_login").(string)),
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			CreateMode:                 postgresql.CreateMode(d.Get("create_mode").(string)),
			InfrastructureEncryption:   infraEncrypt,
			PublicNetworkAccess:        publicAccess,
			SslEnforcement:             ssl,
			StorageProfile:             storage,
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, props)
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

	id, err := parse.PostgresServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("error expanding `sku_name` for PostgreSQL Server %s (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	publicAccess := postgresql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = postgresql.PublicNetworkAccessEnumDisabled
	}

	ssl := postgresql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement"); strings.EqualFold(v.(string), string(postgresql.SslEnforcementEnumDisabled)) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = postgresql.SslEnforcementEnumDisabled
	}

	properties := postgresql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &postgresql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			PublicNetworkAccess:        publicAccess,
			SslEnforcement:             ssl,
			StorageProfile:             expandAzureRmPostgreSQLStorageProfile(d),
			Version:                    postgresql.ServerVersion(d.Get("version").(string)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, properties)
	if err != nil {
		return fmt.Errorf("Error updating PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for update of PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("Error retrieving PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}
	if read.ID == nil {
		return fmt.Errorf("Cannot read PostgreSQL Server %s (resource group %s) ID", id.Name, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmPostgreSQLServerRead(d, meta)
}

func resourceArmPostgreSQLServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PostgresServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] PostgreSQL Server %q was not found (resource group %q)", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Azure PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if p := resp.ServerProperties; p != nil {
		if location := resp.Location; location != nil {
			d.Set("location", azure.NormalizeLocation(*location))
		}

		d.Set("administrator_login", p.AdministratorLogin)
		d.Set("ssl_enforcement", string(p.SslEnforcement))
		d.Set("version", string(p.Version))

		if p.InfrastructureEncryption == postgresql.InfrastructureEncryptionEnabled {
			d.Set("infrastructure_encryption_enabled", true)
		} else if p.InfrastructureEncryption == postgresql.InfrastructureEncryptionDisabled {
			d.Set("infrastructure_encryption_enabled", false)
		}

		if p.PublicNetworkAccess == postgresql.PublicNetworkAccessEnumEnabled {
			d.Set("public_network_access_enabled", true)
		} else if p.PublicNetworkAccess == postgresql.PublicNetworkAccessEnumDisabled {
			d.Set("public_network_access_enabled", false)
		}

		if p.SslEnforcement == postgresql.SslEnforcementEnumEnabled {
			d.Set("ssl_enforcement_enabled", true)
		} else if p.SslEnforcement == postgresql.SslEnforcementEnumDisabled {
			d.Set("ssl_enforcement_enabled", false)
		}

		if err := d.Set("storage_profile", flattenPostgreSQLStorageProfile(p.StorageProfile)); err != nil {
			return fmt.Errorf("Error setting `storage_profile`: %+v", err)
		}

		// Computed
		d.Set("fqdn", p.FullyQualifiedDomainName)
	}
	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPostgreSQLServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Postgres.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PostgresServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing Postgres Server ID : %v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error deleting PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}

		return fmt.Errorf("Error waiting for deletion of PostgreSQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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
	storageMB := storageprofile["storage_mb"].(int)

	autoGrow := postgresql.StorageAutogrowEnabled
	if v, ok := storageprofile["auto_grow"].(string); ok && strings.EqualFold(v, string(postgresql.StorageAutogrowEnabled)) {
		autoGrow = postgresql.StorageAutogrowEnabled
	}
	if v, ok := storageprofile["auto_grow_enabled"].(bool); ok && !v {
		autoGrow = postgresql.StorageAutogrowEnabled
	}

	geoBackup := postgresql.Disabled
	if v, ok := storageprofile["geo_redundant_backup"].(string); ok && strings.EqualFold(v, string(postgresql.Enabled)) {
		geoBackup = postgresql.Enabled
	}
	if v, ok := storageprofile["geo_redundant_backup_enabled"].(bool); ok && v {
		geoBackup = postgresql.Enabled
	}

	return &postgresql.StorageProfile{
		BackupRetentionDays: utils.Int32(int32(backupRetentionDays)),
		GeoRedundantBackup:  geoBackup,
		StorageMB:           utils.Int32(int32(storageMB)),
		StorageAutogrow:     autoGrow,
	}
}

func flattenPostgreSQLStorageProfile(resp *postgresql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["auto_grow"] = string(resp.StorageAutogrow)
	if resp.StorageAutogrow == postgresql.StorageAutogrowEnabled {
		values["auto_grow_enabled"] = true
	} else if resp.StorageAutogrow == postgresql.StorageAutogrowDisabled {
		values["auto_grow_enabled"] = false
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)
	if resp.GeoRedundantBackup == postgresql.Enabled {
		values["geo_redundant_backup_enabled"] = true
	} else if resp.GeoRedundantBackup == postgresql.Disabled {
		values["geo_redundant_backup_enabled"] = false
	}

	return []interface{}{values}
}
