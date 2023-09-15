// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

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
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const (
	mySQLServerResourceName = "azurerm_mysql_server"
)

func resourceMySqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceMySqlServerCreate,
		Read:   resourceMySqlServerRead,
		Update: resourceMySqlServerUpdate,
		Delete: resourceMySqlServerDelete,

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
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
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

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				ForceNew: true,
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

			"identity": commonschema.SystemAssignedIdentityOptional(),

			"ssl_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(servers.MinimalTlsVersionEnumTLSOneTwo),
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.MinimalTlsVersionEnumTLSEnforcementDisabled),
					string(servers.MinimalTlsVersionEnumTLSOneZero),
					string(servers.MinimalTlsVersionEnumTLSOneOne),
					string(servers.MinimalTlsVersionEnumTLSOneTwo),
				}, false),
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.All(
					validation.IntBetween(5120, 16777216),
					validation.IntDivisibleBy(1024),
				),
			},

			"threat_detection_policy": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"disabled_alerts": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Set:      pluginsdk.HashString,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.StringInSlice([]string{
									"Sql_Injection",
									"Sql_Injection_Vulnerability",
									"Access_Anomaly",
									"Data_Exfiltration",
									"Unsafe_Action",
								}, false),
							},
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_account_admins": {
							Type:     pluginsdk.TypeBool,
							Optional: true,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"email_addresses": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								// todo email validation in code
							},
							Set: pluginsdk.HashString,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"retention_days": {
							Type:         pluginsdk.TypeInt,
							Optional:     true,
							ValidateFunc: validation.IntAtLeast(0),
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_account_access_key": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							Sensitive:    true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},

						"storage_endpoint": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringIsNotEmpty,
							AtLeastOneOf: []string{
								"threat_detection_policy.0.enabled", "threat_detection_policy.0.disabled_alerts", "threat_detection_policy.0.email_account_admins",
								"threat_detection_policy.0.email_addresses", "threat_detection_policy.0.retention_days", "threat_detection_policy.0.storage_account_access_key",
								"threat_detection_policy.0.storage_endpoint",
							},
						},
					},
				},
			},

			"tags": commonschema.Tags(),

			"version": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(servers.ServerVersionFivePointSeven),
					string(servers.ServerVersionEightPointZero),
				}, false),
				ForceNew: true,
			},
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, diff *pluginsdk.ResourceDiff, v interface{}) error {
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
		}),
	}
}

func resourceMySqlServerCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.Servers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	securityClient := meta.(*clients.Client).MySQL.MySqlClient.ServerSecurityAlertPolicies
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Server creation.")

	location := azure.NormalizeLocation(d.Get("location").(string))

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_mysql_server", id.ID())
		}
	}

	mode := servers.CreateMode(d.Get("create_mode").(string))
	tlsMin := servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))
	source := d.Get("creation_source_server_id").(string)
	version := servers.ServerVersion(d.Get("version").(string))

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for %s: %v", id, err)
	}

	infraEncrypt := servers.InfrastructureEncryptionEnabled
	if v := d.Get("infrastructure_encryption_enabled"); !v.(bool) {
		infraEncrypt = servers.InfrastructureEncryptionDisabled
	}

	if pointer.From(sku.Tier) == servers.SkuTierBasic && infraEncrypt == servers.InfrastructureEncryptionEnabled {
		return fmt.Errorf("`infrastructure_encryption_enabled` is not supported for sku Tier `Basic` for %s", id)
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled"); !v.(bool) {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled"); !v.(bool) {
		ssl = servers.SslEnforcementEnumDisabled
	}

	storage := expandMySQLStorageProfile(d)

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

		// check admin
		props = &servers.ServerPropertiesForDefaultCreate{
			AdministratorLogin:         admin,
			AdministratorLoginPassword: pass,
			InfrastructureEncryption:   pointer.To(infraEncrypt),
			PublicNetworkAccess:        pointer.To(publicAccess),
			MinimalTlsVersion:          pointer.To(tlsMin),
			SslEnforcement:             pointer.To(ssl),
			StorageProfile:             storage,
			Version:                    pointer.To(version),
		}
	case servers.CreateModePointInTimeRestore:
		v, ok := d.GetOk("restore_point_in_time")
		if !ok || v.(string) == "" {
			return fmt.Errorf("restore_point_in_time must be set when create_mode is PointInTimeRestore")
		}
		t, _ := time.Parse(time.RFC3339, v.(string)) // should be validated by the schema

		props = &servers.ServerPropertiesForRestore{
			SourceServerId:           source,
			RestorePointInTime:       t.Format(time.RFC3339),
			InfrastructureEncryption: pointer.To(infraEncrypt),
			PublicNetworkAccess:      pointer.To(publicAccess),
			MinimalTlsVersion:        pointer.To(tlsMin),
			SslEnforcement:           pointer.To(ssl),
			StorageProfile:           storage,
			Version:                  pointer.To(version),
		}
	case servers.CreateModeGeoRestore:
		props = &servers.ServerPropertiesForGeoRestore{
			SourceServerId:           source,
			InfrastructureEncryption: pointer.To(infraEncrypt),
			PublicNetworkAccess:      pointer.To(publicAccess),
			MinimalTlsVersion:        pointer.To(tlsMin),
			SslEnforcement:           pointer.To(ssl),
			StorageProfile:           storage,
			Version:                  pointer.To(version),
		}
	case servers.CreateModeReplica:
		props = &servers.ServerPropertiesForReplica{
			SourceServerId:           source,
			InfrastructureEncryption: pointer.To(infraEncrypt),
			PublicNetworkAccess:      pointer.To(publicAccess),
			MinimalTlsVersion:        pointer.To(tlsMin),
			SslEnforcement:           pointer.To(ssl),
			Version:                  pointer.To(version),
		}
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	server := servers.ServerForCreate{
		Identity:   expandedIdentity,
		Location:   location,
		Properties: props,
		Sku:        sku,
		Tags:       tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err = client.CreateThenPoll(ctx, id, server); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		securityAlertId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			if err = securityClient.CreateOrUpdateThenPoll(ctx, securityAlertId, *alert); err != nil {
				return fmt.Errorf("updating of Security Alert Policy for %s: %+v", id, err)
			}

		}
	}

	// Issue tracking the REST API update failure: https://github.com/Azure/azure-rest-api-specs/issues/14117
	if mode == servers.CreateModeReplica {
		log.Printf("[INFO] changing `public_network_access_enabled` for %s", id)
		properties := servers.ServerUpdateParameters{
			Properties: &servers.ServerUpdateParametersProperties{
				PublicNetworkAccess: pointer.To(publicAccess),
			},
		}

		if err = client.UpdateThenPoll(ctx, id, properties); err != nil {
			return fmt.Errorf("updating %s: %+v", id, err)
		}
	}

	return resourceMySqlServerRead(d, meta)
}

func resourceMySqlServerUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.Servers
	securityClient := meta.(*clients.Client).MySQL.MySqlClient.ServerSecurityAlertPolicies
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// TODO: support for Delta updates

	log.Printf("[INFO] preparing arguments for AzureRM MySQL Server update.")

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing MySQL Server ID : %v", err)
	}

	sku, err := expandServerSkuName(d.Get("sku_name").(string))
	if err != nil {
		return fmt.Errorf("expanding sku_name for %s: %v", id, err)
	}

	publicAccess := servers.PublicNetworkAccessEnumEnabled
	if v := d.Get("public_network_access_enabled").(bool); !v {
		publicAccess = servers.PublicNetworkAccessEnumDisabled
	}

	ssl := servers.SslEnforcementEnumEnabled
	if v := d.Get("ssl_enforcement_enabled").(bool); !v {
		ssl = servers.SslEnforcementEnumDisabled
	}

	storageProfile := expandMySQLStorageProfile(d)

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	properties := servers.ServerUpdateParameters{
		Identity: expandedIdentity,
		Properties: &servers.ServerUpdateParametersProperties{
			AdministratorLoginPassword: utils.String(d.Get("administrator_login_password").(string)),
			PublicNetworkAccess:        pointer.To(publicAccess),
			SslEnforcement:             pointer.To(ssl),
			MinimalTlsVersion:          pointer.To(servers.MinimalTlsVersionEnum(d.Get("ssl_minimal_tls_version_enforced").(string))),
			StorageProfile:             storageProfile,
			Version:                    pointer.To(servers.ServerVersion(d.Get("version").(string))),
		},
		Sku:  sku,
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.UpdateThenPoll(ctx, *id, properties); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	if v, ok := d.GetOk("threat_detection_policy"); ok {
		securityAlertId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
		alert := expandSecurityAlertPolicy(v)
		if alert != nil {
			if err = securityClient.CreateOrUpdateThenPoll(ctx, securityAlertId, *alert); err != nil {
				return fmt.Errorf("updataing mysql server security alert policy: %v", err)
			}
		}
	}

	return resourceMySqlServerRead(d, meta)
}

func resourceMySqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.Servers
	securityClient := meta.(*clients.Client).MySQL.MySqlClient.ServerSecurityAlertPolicies
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
		d.Set("location", location.Normalize(model.Location))

		tier := servers.SkuTierBasic
		if sku := model.Sku; sku != nil {
			d.Set("sku_name", sku.Name)
			tier = pointer.From(sku.Tier)
		}

		if err := d.Set("identity", identity.FlattenSystemAssigned(model.Identity)); err != nil {
			return fmt.Errorf("setting `identity`: %+v", err)
		}

		if props := model.Properties; props != nil {
			d.Set("administrator_login", props.AdministratorLogin)
			d.Set("infrastructure_encryption_enabled", pointer.From(props.InfrastructureEncryption) == servers.InfrastructureEncryptionEnabled)
			d.Set("public_network_access_enabled", pointer.From(props.PublicNetworkAccess) == servers.PublicNetworkAccessEnumEnabled)
			d.Set("ssl_enforcement_enabled", pointer.From(props.SslEnforcement) == servers.SslEnforcementEnumEnabled)
			minimalTlsVersion := ""
			if props.MinimalTlsVersion != nil {
				minimalTlsVersion = string(*props.MinimalTlsVersion)
			}
			d.Set("ssl_minimal_tls_version_enforced", minimalTlsVersion)
			version := ""
			if props.Version != nil {
				version = string(*props.Version)
			}
			d.Set("version", version)

			if storage := props.StorageProfile; storage != nil {
				d.Set("auto_grow_enabled", pointer.From(storage.StorageAutogrow) == servers.StorageAutogrowEnabled)
				d.Set("backup_retention_days", storage.BackupRetentionDays)
				d.Set("geo_redundant_backup_enabled", pointer.From(storage.GeoRedundantBackup) == servers.GeoRedundantBackupEnabled)
				d.Set("storage_mb", storage.StorageMB)
			}

			// Computed
			d.Set("fqdn", props.FullyQualifiedDomainName)
		}

		// the basic does not support threat detection policies
		if tier == servers.SkuTierGeneralPurpose || tier == servers.SkuTierMemoryOptimized {
			securityAlertId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
			secResp, err := securityClient.Get(ctx, securityAlertId)
			if err != nil && !response.WasNotFound(secResp.HttpResponse) {
				return fmt.Errorf("retrieving Security Alert Policy for %s: %+v", *id, err)
			}

			if !response.WasNotFound(secResp.HttpResponse) && secResp.Model != nil {
				block := flattenSecurityAlertPolicy(secResp.Model.Properties, d.Get("threat_detection_policy.0.storage_account_access_key").(string))
				if err := d.Set("threat_detection_policy", block); err != nil {
					return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
				}
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}

func resourceMySqlServerDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.Servers
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := servers.ParseServerID(d.Id())
	if err != nil {
		return fmt.Errorf("parsing MySQL Server ID : %v", err)
	}

	if err = client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandServerSkuName(skuName string) (*servers.Sku, error) {
	parts := strings.Split(skuName, "_")
	if len(parts) != 3 {
		return nil, fmt.Errorf("sku_name (%s) has the worng numberof parts (%d) after splitting on _", skuName, len(parts))
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
		return nil, fmt.Errorf("cannot convert skuname %s capcity %s to int", skuName, parts[2])
	}

	return &servers.Sku{
		Name:     skuName,
		Tier:     pointer.To(tier),
		Capacity: utils.Int64(int64(capacity)),
		Family:   utils.String(parts[1]),
	}, nil
}

func expandMySQLStorageProfile(d *pluginsdk.ResourceData) *servers.StorageProfile {
	storage := servers.StorageProfile{}

	// now override whatever we may have from the block with the top level properties
	if v, ok := d.GetOk("auto_grow_enabled"); ok {
		storage.StorageAutogrow = pointer.To(servers.StorageAutogrowDisabled)
		if v.(bool) {
			storage.StorageAutogrow = pointer.To(servers.StorageAutogrowEnabled)
		}
	}

	if v, ok := d.GetOk("backup_retention_days"); ok {
		storage.BackupRetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := d.GetOk("geo_redundant_backup_enabled"); ok {
		storage.GeoRedundantBackup = pointer.To(servers.GeoRedundantBackupDisabled)
		if v.(bool) {
			storage.GeoRedundantBackup = pointer.To(servers.GeoRedundantBackupEnabled)
		}
	}

	if v, ok := d.GetOk("storage_mb"); ok {
		storage.StorageMB = utils.Int64(int64(v.(int)))
	}

	return &storage
}

func expandSecurityAlertPolicy(i interface{}) *serversecurityalertpolicies.ServerSecurityAlertPolicy {
	slice := i.([]interface{})
	if len(slice) == 0 {
		return nil
	}

	block := slice[0].(map[string]interface{})

	state := serversecurityalertpolicies.ServerSecurityAlertPolicyStateEnabled
	if !block["enabled"].(bool) {
		state = serversecurityalertpolicies.ServerSecurityAlertPolicyStateDisabled
	}

	props := &serversecurityalertpolicies.SecurityAlertPolicyProperties{
		State: state,
	}

	if v, ok := block["disabled_alerts"]; ok {
		props.DisabledAlerts = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := block["email_addresses"]; ok {
		props.EmailAddresses = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
	}

	if v, ok := block["email_account_admins"]; ok {
		props.EmailAccountAdmins = utils.Bool(v.(bool))
	}

	if v, ok := block["retention_days"]; ok {
		props.RetentionDays = utils.Int64(int64(v.(int)))
	}

	if v, ok := block["storage_account_access_key"]; ok && v.(string) != "" {
		props.StorageAccountAccessKey = utils.String(v.(string))
	}

	if v, ok := block["storage_endpoint"]; ok && v.(string) != "" {
		props.StorageEndpoint = utils.String(v.(string))
	}

	return &serversecurityalertpolicies.ServerSecurityAlertPolicy{
		Properties: props,
	}
}

func flattenSecurityAlertPolicy(props *serversecurityalertpolicies.SecurityAlertPolicyProperties, accessKey string) interface{} {
	if props == nil {
		return nil
	}

	// check if its an empty block as in its never been set before
	if props.DisabledAlerts != nil && len(*props.DisabledAlerts) == 1 && (*props.DisabledAlerts)[0] == "" &&
		props.EmailAddresses != nil && len(*props.EmailAddresses) == 1 && (*props.EmailAddresses)[0] == "" &&
		props.StorageAccountAccessKey != nil && *props.StorageAccountAccessKey == "" &&
		props.StorageEndpoint != nil && *props.StorageEndpoint == "" &&
		props.RetentionDays != nil && *props.RetentionDays == 0 &&
		props.EmailAccountAdmins != nil && !*props.EmailAccountAdmins &&
		props.State == serversecurityalertpolicies.ServerSecurityAlertPolicyStateDisabled {
		return nil
	}

	block := map[string]interface{}{}

	block["enabled"] = props.State == serversecurityalertpolicies.ServerSecurityAlertPolicyStateEnabled

	// the service will return "disabledAlerts":[""] for empty
	if props.DisabledAlerts == nil || len(*props.DisabledAlerts) == 0 || (*props.DisabledAlerts)[0] == "" {
		block["disabled_alerts"] = []interface{}{}
	} else {
		block["disabled_alerts"] = utils.FlattenStringSlice(props.DisabledAlerts)
	}

	// the service will return "emailAddresses":[""] for empty
	if props.EmailAddresses == nil || len(*props.EmailAddresses) == 0 || (*props.EmailAddresses)[0] == "" {
		block["email_addresses"] = []interface{}{}
	} else {
		block["email_addresses"] = utils.FlattenStringSlice(props.EmailAddresses)
	}

	if v := props.EmailAccountAdmins; v != nil {
		block["email_account_admins"] = *v
	}
	if v := props.RetentionDays; v != nil {
		block["retention_days"] = *v
	}
	if v := props.StorageEndpoint; v != nil {
		block["storage_endpoint"] = *v
	}

	block["storage_account_access_key"] = accessKey

	return []interface{}{block}
}
