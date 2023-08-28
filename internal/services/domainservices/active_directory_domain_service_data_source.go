// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package domainservices

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

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

			"location": commonschema.LocationComputed(),

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
						"kerberos_armoring_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

						"kerberos_rc4_encryption_enabled": {
							Type:     pluginsdk.TypeBool,
							Computed: true,
						},

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

			"tags": commonschema.Tags(),

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

		"location": commonschema.ResourceGroupNameForDataSource(),

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
	subscrptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	idsdk := domainservices.NewDomainServiceID(subscrptionId, resourceGroup, name)

	resp, err := client.Get(ctx, idsdk)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}
		return err
	}

	model := resp.Model
	if model == nil {
		return fmt.Errorf("reading Domain Service: model was returned nil")
	}

	if model.Id == nil {
		return fmt.Errorf("reading Domain Service: ID was returned nil")
	}

	d.SetId(idsdk.ID())

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	d.Set("location", location.NormalizeNilable(model.Location))
	if err := tags.FlattenAndSet(d, model.Tags); err != nil {
		return err
	}

	if props := model.Properties; props != nil {
		d.Set("deployment_id", props.DeploymentId)

		domainConfigType := ""
		if v := props.DomainConfigurationType; v != nil {
			domainConfigType = *v
		}
		d.Set("domain_configuration_type", domainConfigType)

		d.Set("domain_name", props.DomainName)

		d.Set("filtered_sync_enabled", false)
		if props.FilteredSync != nil && *props.FilteredSync == domainservices.FilteredSyncEnabled {
			d.Set("filtered_sync_enabled", true)
		}

		d.Set("resource_id", model.Id)
		d.Set("sku", props.Sku)
		d.Set("sync_owner", props.SyncOwner)
		d.Set("tenant_id", props.TenantId)
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

	return nil
}
