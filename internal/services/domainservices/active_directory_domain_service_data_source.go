package domainservices

import (
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2020-01-01/aad"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceActiveDirectoryDomainService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceActiveDirectoryDomainServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"deployment_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"domain_configuration_type": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"domain_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"filtered_sync_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"location": azure.SchemaLocationForDataSource(),

			"notifications": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"additional_recipients": {
							Type:     pluginsdk.TypeList,
							Computed: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"notify_dc_admins": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
						"notify_global_admins": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"replica_sets": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: dataSourceActiveDirectoryDomainServiceReplicaSetSchema(),
				},
			},

			"resource_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"secure_ldap": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"external_access_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"certificate_expiry": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"certificate_thumbprint": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},

						"public_certificate": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"security": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"ntlm_v1_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"sync_kerberos_passwords": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"sync_ntlm_passwords": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"sync_on_prem_passwords": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"tls_v1_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"sku": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"sync_owner": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.SchemaDataSource(),

			"tenant_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"version": {
				Type:     pluginsdk.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceActiveDirectoryDomainServiceReplicaSetSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		// TODO: add health-related attributes

		"domain_controller_ip_addresses": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"external_access_ip_address": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": azure.SchemaLocationForDataSource(),

		"service_status": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"subnet_id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},
	}
}

func dataSourceActiveDirectoryDomainServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	resp, err := client.Get(ctx, resourceGroup, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}
		return err
	}

	if resp.ID == nil {
		return fmt.Errorf("reading Domain Service: ID was returned nil")
	}
	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if resp.Location == nil {
		return fmt.Errorf("reading Domain Service %q: location was returned nil", d.Id())
	}
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.DomainServiceProperties; props != nil {
		d.Set("deployment_id", props.DeploymentID)

		domainConfigType := ""
		if v := props.DomainConfigurationType; v != nil {
			domainConfigType = *v
		}
		d.Set("domain_configuration_type", domainConfigType)

		d.Set("domain_name", props.DomainName)

		d.Set("filtered_sync_enabled", false)
		if props.FilteredSync == aad.FilteredSyncEnabled {
			d.Set("filtered_sync_enabled", true)
		}

		d.Set("resource_id", resp.ID)
		d.Set("sku", props.Sku)
		d.Set("sync_owner", props.SyncOwner)
		d.Set("tenant_id", props.TenantID)
		d.Set("version", props.Version)

		if err := d.Set("notifications", flattenDomainServiceNotifications(props.NotificationSettings)); err != nil {
			return fmt.Errorf("setting `notifications`: %+v", err)
		}

		if err := d.Set("secure_ldap", flattenDomainServiceLdaps(d, props.LdapsSettings, true)); err != nil {
			return fmt.Errorf("setting `secure_ldap`: %+v", err)
		}

		if err := d.Set("security", flattenDomainServiceSecurity(props.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("setting `security`: %+v", err)
		}

		replicaSets := flattenDomainServiceReplicaSets(props.ReplicaSets)
		if err := d.Set("replica_sets", replicaSets); err != nil {
			return fmt.Errorf("setting `replica_sets`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
