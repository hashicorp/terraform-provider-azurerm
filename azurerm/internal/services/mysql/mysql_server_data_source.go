package mysql

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceMySqlServer() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceMySqlServerRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ServerName,
			},

			"administrator_login": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"auto_grow_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"backup_retention_days": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"geo_redundant_backup_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"infrastructure_encryption_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"restore_point_in_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Computed: true,
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

			"ssl_enforcement_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"storage_mb": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},

			"threat_detection_policy": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"disabled_alerts": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},

						"email_account_admins": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"email_addresses": {
							Type:     pluginsdk.TypeSet,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							Set: pluginsdk.HashString,
						},

						"retention_days": {
							Type:     pluginsdk.TypeInt,
							Computed: true,
						},

						"storage_account_access_key": {
							Type:      pluginsdk.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"storage_endpoint": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMySqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.ServersClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	securityClient := meta.(*clients.Client).MySQL.ServerSecurityAlertPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))

	tier := mysql.Basic
	if sku := resp.Sku; sku != nil {
		d.Set("sku_name", sku.Name)
		tier = sku.Tier
	}

	if err := d.Set("identity", flattenServerIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.ServerProperties; props != nil {
		d.Set("administrator_login", props.AdministratorLogin)
		d.Set("infrastructure_encryption_enabled", props.InfrastructureEncryption == mysql.InfrastructureEncryptionEnabled)
		d.Set("public_network_access_enabled", props.PublicNetworkAccess == mysql.PublicNetworkAccessEnumEnabled)
		d.Set("ssl_enforcement_enabled", props.SslEnforcement == mysql.SslEnforcementEnumEnabled)
		d.Set("ssl_minimal_tls_version_enforced", props.MinimalTLSVersion)
		d.Set("version", string(props.Version))

		if storage := props.StorageProfile; storage != nil {
			d.Set("auto_grow_enabled", storage.StorageAutogrow == mysql.StorageAutogrowEnabled)
			d.Set("backup_retention_days", storage.BackupRetentionDays)
			d.Set("geo_redundant_backup_enabled", storage.GeoRedundantBackup == mysql.Enabled)
			d.Set("storage_mb", storage.StorageMB)
		}

		d.Set("fqdn", props.FullyQualifiedDomainName)
	}

	if tier == mysql.GeneralPurpose || tier == mysql.MemoryOptimized {
		secResp, err := securityClient.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil && !utils.ResponseWasNotFound(secResp.Response) {
			return fmt.Errorf("retrieving Security Alert Policy for %s: %+v", id, err)
		}

		accountKey := ""
		if secResp.SecurityAlertPolicyProperties != nil && secResp.SecurityAlertPolicyProperties.StorageAccountAccessKey != nil {
			accountKey = *secResp.SecurityAlertPolicyProperties.StorageAccountAccessKey
		}

		if !utils.ResponseWasNotFound(secResp.Response) {
			block := flattenSecurityAlertPolicy(secResp.SecurityAlertPolicyProperties, accountKey)
			if err := d.Set("threat_detection_policy", block); err != nil {
				return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
