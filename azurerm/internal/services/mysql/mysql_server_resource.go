package mysql

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2017-12-01/mysql"
	"github.com/Azure/go-autorest/autorest/date"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
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
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				if _, err := parse.MysqlServerServerID(d.Id()); err != nil {
					return []*schema.ResourceData{d}, err
				}

				d.Set("create_mode", "Default")
				if v, ok := d.GetOk("create_mode"); ok && v.(string) != "" {
					d.Set("create_mode", v)
				}

				return []*schema.ResourceData{d}, nil
			},
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
				ValidateFunc: validate.MysqlServerServerName,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"auto_grow_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true, // TODO: remove in 3.0 and default to true
				ConflictsWith: []string{"storage_profile.0.auto_grow"},
			},

			"backup_retention_days": {
				Type:          schema.TypeInt,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_profile.0.backup_retention_days"},
				ValidateFunc:  validation.IntBetween(7, 35),
			},

			"create_mode": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(mysql.CreateModeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.CreateModeDefault),
					string(mysql.CreateModeGeoRestore),
					string(mysql.CreateModePointInTimeRestore),
					string(mysql.CreateModeReplica),
				}, false),
			},

			"creation_source_server_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validate.MysqlServerServerID,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"geo_redundant_backup_enabled": {
				Type:          schema.TypeBool,
				Optional:      true,
				Computed:      true,
				ConflictsWith: []string{"storage_profile.0.geo_redundant_backup"},
			},

			"infrastructure_encryption_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},

			"location": azure.SchemaLocation(),

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"restore_point_in_time": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

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

			"ssl_enforcement": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				Deprecated:   "this has been moved to the boolean attribute `ssl_enforcement_enabled` and will be removed in version 3.0 of the provider.",
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.SslEnforcementEnumDisabled),
					string(mysql.SslEnforcementEnumEnabled),
				}, true),
				DiffSuppressFunc: suppress.CaseDifference,
			},

			"ssl_enforcement_enabled": {
				Type:         schema.TypeBool,
				Optional:     true, // required in 3.0
				ExactlyOneOf: []string{"ssl_enforcement", "ssl_enforcement_enabled"},
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(mysql.TLSEnforcementDisabled),
				ValidateFunc: validation.StringInSlice([]string{
					string(mysql.TLSEnforcementDisabled),
					string(mysql.TLS10),
					string(mysql.TLS11),
					string(mysql.TLS12),
				}, false),
			},

			"storage_mb": {
				Type:         schema.TypeInt,
				Optional:     true,
				Computed:     true,
				ExactlyOneOf: []string{"storage_profile.0.storage_mb"},
				ValidateFunc: validation.All(
					validation.IntBetween(5120, 4194304),
					validation.IntDivisibleBy(1024),
				),
			},

			"storage_profile": {
				Type:       schema.TypeList,
				Optional:   true,
				Computed:   true,
				MaxItems:   1,
				Deprecated: "all storage_profile properties have been moved to the top level. This block will be removed in version 3.0 of the provider.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"auto_grow": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"auto_grow_enabled"},
							Deprecated:       "this has been moved to the top level boolean attribute `auto_grow_enabled` and will be removed in version 3.0 of the provider.",
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								string(mysql.StorageAutogrowEnabled),
								string(mysql.StorageAutogrowDisabled),
							}, false),
						},
						"backup_retention_days": {
							Type:          schema.TypeInt,
							Optional:      true,
							Computed:      true,
							ConflictsWith: []string{"backup_retention_days"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc:  validation.IntBetween(7, 35),
						},
						"geo_redundant_backup": {
							Type:             schema.TypeString,
							Optional:         true,
							Computed:         true,
							ConflictsWith:    []string{"geo_redundant_backup_enabled"},
							Deprecated:       "this has been moved to the top level boolean attribute `geo_redundant_backup_enabled` and will be removed in version 3.0 of the provider.",
							DiffSuppressFunc: suppress.CaseDifference,
							ValidateFunc: validation.StringInSlice([]string{
								"Enabled",
								"Disabled",
							}, true),
						},
						"storage_mb": {
							Type:          schema.TypeInt,
							Optional:      true,
							ConflictsWith: []string{"storage_mb"},
							Deprecated:    "this has been moved to the top level and will be removed in version 3.0 of the provider.",
							ValidateFunc: validation.All(
								validation.IntBetween(5120, 4194304),
								validation.IntDivisibleBy(1024),
							),
						},
					},
				},
			},

			"tags": tags.Schema(),

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
		},

		CustomizeDiff: func(diff *schema.ResourceDiff, v interface{}) error {
			tier, _ := diff.GetOk("sku_name")

			var storageMB int
			if v, ok := diff.GetOk("storage_mb"); ok {
				storageMB = v.(int)
			} else if v, ok := diff.GetOk("storage_profile.0.storage_mb"); ok {
				storageMB = v.(int)
			}

			if strings.HasPrefix(tier.(string), "B_") && storageMB > 1048576 {
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
				return fmt.Errorf("checking for presence of existing MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_mysql_server", *existing.ID)
		}
	}

	mode := mysql.CreateMode(d.Get("create_mode").(string))
	tlsMin := mysql.MinimalTLSVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	source := d.Get("creation_source_server_id").(string)
	version := mysql.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for MySQL Server %q (Resource Group %q): %v", name, resourceGroup, err)
	}

	infraEncrypt := mysql.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = mysql.InfrastructureEncryptionDisabled
	}

	publicAccess := mysql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = mysql.PublicNetworkAccessEnumDisabled
	}

	ssl := mysql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = mysql.SslEnforcementEnumDisabled
	}

	storage := expandMySQLStorageProfile(d)

	var props mysql.BasicServerPropertiesForCreate
	switch mode {
	case mysql.CreateModeDefault:
		admin := d.Get("administrator_login").(string)
		pass := d.Get("administrator_login_password").(string)

		if admin == "" {
			return fmt.Errorf("`administrator_login` must not be empty when `create_mode` is `default`")
		}
		if pass == "" {
			return fmt.Errorf("`administrator_login_password` must not be empty when `create_mode` is `default`")
		}

		if _, ok := d.GetOk("restore_point_in_time"); ok {
			return fmt.Errorf("`restore_point_in_time` cannot be set when `create_mode` is `default`")
		}

		// check admin
		props = &mysql.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         &admin,
			AdministratorLoginPassword: &pass,
			CreateMode:                 mode,
			InfrastructureEncryption:   infraEncrypt,
			PublicNetworkAccess:        publicAccess,
			MinimalTLSVersion:          tlsMin,
			SslEnforcement:             ssl,
			StorageProfile:             storage,
			Version:                    version,
		}
	case mysql.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}
		time, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema

		props = &mysql.ServerPropertiesForRestore{
			CreateMode:     mode,
			SourceServerID: &source,
			RestorePointInTime: &date.Time{
				Time: time,
			},
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case mysql.CreateModeGeoRestore:
		props = &mysql.ServerPropertiesForGeoRestore{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			StorageProfile:           storage,
			Version:                  version,
		}
	case mysql.CreateModeReplica:
		props = &mysql.ServerPropertiesForReplica{
			CreateMode:               mode,
			SourceServerID:           &source,
			InfrastructureEncryption: infraEncrypt,
			PublicNetworkAccess:      publicAccess,
			MinimalTLSVersion:        tlsMin,
			SslEnforcement:           ssl,
			Version:                  version,
		}
	}

	server := mysql.ServerForCreate{
		Location:   &location,
		Properties: props,
		Sku:        sku,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Create(ctx, resourceGroup, name, server)
	if err != nil {
		return fmt.Errorf("creating MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		return fmt.Errorf("retrieving MySQL Server %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read MySQL Server %q (Resource Group %q) ID", name, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySqlServerRead(d, meta)
}

func resourceArmMySqlServerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: support for Delta updates

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Server update.")

	id, err := parse.MysqlServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing MySQL Server ID : %v", err)
	}

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for MySQL Server %q (Resource Group %q): %v", id.Name, id.ResourceGroup, err)
	}

	publicAccess := mysql.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled").(bool); !v {
		publicAccess = mysql.PublicNetworkAccessEnumDisabled
	}

	ssl := mysql.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled").(bool); !v {
		ssl = mysql.SslEnforcementEnumDisabled
	}

	storageProfile := expandMySQLStorageProfile(d)

	properties := mysql.ServerUpdateParameters{
		ServerUpdateParametersProperties: &mysql.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			PublicNetworkAccess:        publicAccess,
			SslEnforcement:             ssl,
			StorageProfile:             storageProfile,
			Version:                    mysql.ServerVersion(d.Get("version").(string)),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.Update(ctx, id.ResourceGroup, id.Name, properties)
	if err != nil {
		return fmt.Errorf("updating MySQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for MySQL Server %q (Resource Group %q) to finish updating: %+v", id.Name, id.ResourceGroup, err)
	}

	read, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving MySQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("cannot read MySQL Server %q (Resource Group %q) ID", id.Name, id.ResourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmMySqlServerRead(d, meta)
}

func resourceArmMySqlServerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MysqlServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing MySQL Server ID : %v", err)
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[WARN] MySQL Server %q was not found (Resource Group %q)", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Azure MySQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", id.ResourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("infrastructure_encryption_enabled", props.InfrastructureEncryption == mysql.InfrastructureEncryptionEnabled)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == mysql.PublicNetworkAccessEnumEnabled)
		d.Set("ssl_enforcement", string(props.SslEnforcement))
		d.Set("ssl_enforcement_enabled", props.SslEnforcement == mysql.SslEnforcementEnumEnabled)
		d.Set("ssl_minimal_tls_version_enforced", props.MinimalTLSVersion)
		d.Set("version", string(props.Version))

		if err := d.Set("storage_profile", flattenMySQLStorageProfile(resp.StorageProfile)); err != nil {
			return fmt.Errorf("setting `storage_profile`: %+v", err)
		}

		if storage := props.StorageProfile; storage != nil {
			d.Set("auto_grow_enabled", storage.StorageAutogrow == mysql.StorageAutogrowEnabled)
			d.Set("backup_retention_days", storage.BackupRetentionDays)
			d.Set("geo_redundant_backup_enabled", storage.GeoRedundantBackup == mysql.Enabled)
			d.Set("storage_mb", storage.StorageMB)
		}

		// Computed
		d.Set("fqdn", props.FullyQualifiedDomainName)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmMySqlServerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.MysqlServerServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing MySQL Server ID : %v", err)
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("deleting MySQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("waiting for deletion of MySQL Server %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
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

func expandMySQLStorageProfile(d *schema.ResourceData) *mysql.StorageProfile {
	storage := mysql.StorageProfile{}
	if v, ok := d.GetOk("storage_profile"); ok {
		storageprofile := v.([]interface{})[0].(map[string]interface{})

		storage.BackupRetentionDays = utils.Int32(int32(storageprofile["backup_retention_days"].(int)))
		storage.GeoRedundantBackup = mysql.GeoRedundantBackup(storageprofile["geo_redundant_backup"].(string))
		storage.StorageAutogrow = mysql.StorageAutogrow(storageprofile["auto_grow"].(string))
		storage.StorageMB = utils.Int32(int32(storageprofile["storage_mb"].(int)))
	}

	// now override whatever we may have from the block with the top level properties
	if v, ok := d.GetOk("auto_grow_enabled"); ok {
		storage.StorageAutogrow = mysql.StorageAutogrowDisabled
		if v.(bool) {
			storage.StorageAutogrow = mysql.StorageAutogrowEnabled
		}
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		storage.BackupRetentionDays = utils.Int32(int32(v.(int)))
	}

	if v, ok := d.GetOk("geo_redundant_backup_enabled"); ok {
		storage.GeoRedundantBackup = mysql.Disabled
		if v.(bool) {
			storage.GeoRedundantBackup = mysql.Enabled
		}
	}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageMB = utils.Int32(int32(v.(int)))
	}

	return &storage
}

func flattenMySQLStorageProfile(resp *mysql.StorageProfile) []interface{} {
	values := map[string]interface{}{}

	values["auto_grow"] = string(resp.StorageAutogrow)

	values["backup_retention_days"] = nil
	if backupRetentionDays := resp.BackupRetentionDays; backupRetentionDays != nil {
		values["backup_retention_days"] = *backupRetentionDays
	}

	values["geo_redundant_backup"] = string(resp.GeoRedundantBackup)

	values["storage_mb"] = nil
	if storageMB := resp.StorageMB; storageMB != nil {
		values["storage_mb"] = *storageMB
	}

	return []interface{}{values}
}
