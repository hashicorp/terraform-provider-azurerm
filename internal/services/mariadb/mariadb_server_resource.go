// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mariadb

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mariadb/2018-06-01/servers"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mariadb/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceMariaDbServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMariaDbServerCreate,
		Read:   resourceMariaDbServerRead,
		Update: resourceMariaDbServerUpdate,
		Delete: resourceMariaDbServerDelete,

		Importer: pluginsdk.ImporterValidatingResourceIdThen(func(id string) error {
			_, err := servers.ParseServerID(id)
			return err
		}, func(ctx context.Context, d *pluginsdk.ResourceData, meta interface{}) ([]*pluginsdk.ResourceData, error) {
			d.Set("create_mode", "Default")
			if v, ok := d.GetOk("create_mode"); ok && v.(string) != "" {
				d.Set("create_mode", v)
			}

			return []*pluginsdk.ResourceData{d}, nil
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(90 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ServerName,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"administrator_login_password": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"auto_grow_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"backup_retention_days": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validation.IntBetween(7, 35),
			},

			"create_mode": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(servers.CreateModeDefault),
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.CreateModeDefault),
					string(servers.CreateModeGeoRestore),
					string(servers.CreateModePointInTimeRestore),
					string(servers.CreateModeReplica),
				}, false),
			},

			"creation_source_server_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: servers.ValidateServerID,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Computed: true,
			},

			"location": commonschema.Location(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"restore_point_in_time": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
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

			"ssl_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.IntBetween(5120, 4194304),
					validation.IntDivisibleBy(1024),
				),
			},
			"ssl_minimal_tls_version_enforced": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Default:      string(servers.MinimalTlsVersionEnumTLSOneTwo),
				ValidateFunc: validation.StringInSlice(servers.PossibleValuesForMinimalTlsVersionEnum(), false),
			},
			"tags": commonschema.Tags(),

			"version": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"10.2",
					"10.3",
				}, false),
			},
		},
	}
}

func resourceMariaDbServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mariadb_server", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	mode := servers.CreateMode(d.Get("create_mode").(string))
	source := d.Get("creation_source_server_id").(string)
	version := servers.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name`: %+v", err)
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled").(bool); !v {
		ssl = servers.SslEnforcementEnumDisabled
	}

	tlsMin := servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	if ssl == servers.SslEnforcementEnumDisabled && tlsMin != servers.MinimalTlsVersionEnumTLSEnforcementDisabled {
		return fmt.Errorf("`ssl_minimal_tls_version_enforced` must be set to `TLSEnforcementDisabled` if `ssl_enforcement_enabled` is set to `false`")
	}

	storage := expandMariaDbStorageProfile(d)

	var props servers.ServerPropertiesForCreate
	switch mode {
	case servers.CreateModeDefault:
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

		props = servers.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         admin,
			AdministratorLoginPassword: pass,
			PublicNetworkAccess:        &publicAccess,
			MinimalTlsVersion:          &tlsMin,
			SslEnforcement:             &ssl,
			StorageProfile:             storage,
			Version:                    &version,
		}

	case servers.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}

		props = &servers.ServerPropertiesForRestore{
			SourceServerId:      source,
			RestorePointInTime:  v.(string),
			PublicNetworkAccess: &publicAccess,
			SslEnforcement:      &ssl,
			StorageProfile:      storage,
			Version:             &version,
		}
	case servers.CreateModeGeoRestore:
		props = &servers.ServerPropertiesForGeoRestore{
			SourceServerId:      source,
			PublicNetworkAccess: &publicAccess,
			SslEnforcement:      &ssl,
			StorageProfile:      storage,
			Version:             &version,
		}
	case servers.CreateModeReplica:
		props = &servers.ServerPropertiesForReplica{
			SourceServerId:      source,
			PublicNetworkAccess: &publicAccess,
			SslEnforcement:      &ssl,
			Version:             &version,
		}
	}

	server := servers.ServerForCreate{
		Location:   location,
		Properties: props,
		Sku:        sku,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateThenPoll(ctx, id, server); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceMariaDbServerRead(d, meta)
}

func resourceMariaDbServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MariaDB Server update.")

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding `sku_name`: %+v", err)
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled").(bool); !v {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled").(bool); !v {
		ssl = servers.SslEnforcementEnumDisabled
	}

	tlsMin := servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))

	if ssl == servers.SslEnforcementEnumDisabled && tlsMin != servers.MinimalTlsVersionEnumTLSEnforcementDisabled {
		return fmt.Errorf("`ssl_minimal_tls_version_enforced` must be set to `TLSEnforcementDisabled` if `ssl_enforcement_enabled` is set to `false`")
	}

	storageProfile := expandMariaDbStorageProfile(d)
	serverVersion := servers.ServerVersion(d.Get("version").(string))
	properties := servers.ServerUpdateParameters{
		Properties: &servers.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			PublicNetworkAccess:        &publicAccess,
			MinimalTlsVersion:          &tlsMin,
			SslEnforcement:             &ssl,
			StorageProfile:             storageProfile,
			Version:                    &serverVersion,
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceMariaDbServerRead(d, meta)
}

func resourceMariaDbServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.ServerName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", azure.NormalizeLocation(model.Location))

		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
		}

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("ssl_minimal_tls_version_enforced", string(pointer.From(props.MinimalTlsVersion)))

			publicNetworkAccess := false
			if props.PublicNetworkAccess != nil {
				publicNetworkAccess = *props.PublicNetworkAccess == servers.PublicNetworkAccessEnumEnabled
			}
			d.Set("public_network_access_enabled", publicNetworkAccess)

			sslEnforcement := false
			if props.SslEnforcement != nil {
				sslEnforcement = *props.SslEnforcement == servers.SslEnforcementEnumEnabled
			}
			d.Set("ssl_enforcement_enabled", sslEnforcement)

			version := ""
			if props.Version != nil {
				version = string(*props.Version)
			}
			d.Set("version", version)

			if storage := props.StorageProfile; storage != nil {
				autoGrow := false
				if storage.StorageAutogrow != nil {
					autoGrow = *storage.StorageAutogrow == servers.StorageAutogrowEnabled
				}
				d.Set("auto_grow_enabled", autoGrow)

				geoRedundant := false
				if storage.GeoRedundantBackup != nil {
					geoRedundant = *storage.GeoRedundantBackup == servers.GeoRedundantBackupEnabled
				}
				d.Set("geo_redundant_backup_enabled", geoRedundant)
				d.Set("backup_retention_days", storage.BackupRetentionDays)
				d.Set("storage_mb", storage.StorageMB)
			}

			// Computed
			d.Set("fqdn", props.FullyQualifiedDomainName)
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceMariaDbServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MariaDB.ServersClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandServerSkuName(skuName string) (*servers.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the wrong number of parts (%d) after splitting on _", skuName, len(parts))
	}

	var tier servers.SkuTier
	switch parts[0] {
	case "B":
		tier = servers.SkuTierBasic
	case "GP":
		tier = servers.SkuTierGeneralPurpose
	case "MO":
		tier = servers.SkuTierMemoryOptimized
	default:
		return nil, fmt.Errorf("sku_name %s has unknown sku tier %s", skuName, parts[0])
	}

	capacity, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, fmt.Errorf("cannot convert `sku_name` %q capacity %s to int", skuName, parts[2])
	}

	return &servers.Sku{
		Name:     skuName,
		Tier:     &tier,
		Capacity: utils.Int64(int64(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandMariaDbStorageProfile(d *pluginsdk.ResourceData) *servers.StorageProfile {
	storage := servers.StorageProfile{}
	// now override whatever we may have from the block with the top level properties
	if v, ok := d.GetOk("auto_grow_enabled"); ok {
		autogrowEnabled := servers.StorageAutogrowDisabled
		if v.(bool) {
			autogrowEnabled = servers.StorageAutogrowEnabled
		}
		storage.StorageAutogrow = &autogrowEnabled
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		storage.BackupRetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("geo_redundant_backup_enabled"); ok {
		geoRedundantBackup := servers.GeoRedundantBackupDisabled
		if v.(bool) {
			geoRedundantBackup = servers.GeoRedundantBackupEnabled
		}
		storage.GeoRedundantBackup = &geoRedundantBackup
	}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageMB = utils.Int64(int64(v.(int)))
	}

	return &storage
}
