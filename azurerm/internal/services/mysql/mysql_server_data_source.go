package mysql

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/parse"

	"github.com/Azure/azure-sdk-for-go/services/mysql/mgmt/2020-01-01/mysql"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/mysql/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceMySqlServer() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMySqlServerRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.ServerName,
			},

			"administrator_login": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"auto_grow_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"backup_retention_days": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"fqdn": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"geo_redundant_backup_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"infrastructure_encryption_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"public_network_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"restore_point_in_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"identity": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"ssl_enforcement_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"ssl_minimal_tls_version_enforced": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"storage_mb": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"threat_detection_policy": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"disabled_alerts": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"email_account_admins": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"email_addresses": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Set: schema.HashString,
						},

						"retention_days": {
							Type:     schema.TypeInt,
							Computed: true,
						},

						"storage_account_access_key": {
							Type:      schema.TypeString,
							Computed:  true,
							Sensitive: true,
						},

						"storage_endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"tags": tags.SchemaDataSource(),

			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMySqlServerRead(d *schema.ResourceData, meta interface{}) error {
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
