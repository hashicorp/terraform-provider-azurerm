// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package mysql

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/servers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/mysql/2017-12-01/serversecurityalertpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/mysql/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"location": commonschema.LocationComputed(),

			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"restore_point_in_time": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"identity": commonschema.SystemAssignedIdentityComputed(),

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

			"tags": commonschema.TagsDataSource(),

			"version": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMySqlServerRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).MySQL.MySqlClient.Servers
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	securityClient := meta.(*clients.Client).MySQL.MySqlClient.ServerSecurityAlertPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := servers.NewServerID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.SetId(id.ID())
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

			d.Set("fqdn", props.FullyQualifiedDomainName)
		}

		if tier == servers.SkuTierGeneralPurpose || tier == servers.SkuTierMemoryOptimized {
			securityServerId := serversecurityalertpolicies.NewServerID(id.SubscriptionId, id.ResourceGroupName, id.ServerName)
			secResp, err := securityClient.Get(ctx, securityServerId)
			if err != nil && !response.WasNotFound(secResp.HttpResponse) {
				return fmt.Errorf("retrieving Security Alert Policy for %s: %+v", id, err)
			}

			accountKey := ""
			if secResp.Model != nil && secResp.Model.Properties != nil && secResp.Model.Properties.StorageAccountAccessKey != nil {
				accountKey = *secResp.Model.Properties.StorageAccountAccessKey
			}

			if !response.WasNotFound(secResp.HttpResponse) {
				block := flattenSecurityAlertPolicy(secResp.Model.Properties, accountKey)
				if err := d.Set("threat_detection_policy", block); err != nil {
					return fmt.Errorf("setting `threat_detection_policy`: %+v", err)
				}
			}
		}

		return tags.FlattenAndSet(d, model.Tags)
	}

	return nil
}
