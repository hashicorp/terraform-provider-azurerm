// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package active_directory_domain_service

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/aad/2021-05-01/domainservices"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/domainservices/helpers"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func RegisterDataSource() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Description: "Gets information about an Active Directory Domain Service.",
		Read:        dataSourceActiveDirectoryDomainServiceRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Description:  "The display name for your managed Active Directory Domain Service resource.",
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"resource_group_name": commonschema.ResourceGroupNameForDataSource(),

			"deployment_id": {
				Description: "A unique ID for the managed domain deployment.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"domain_configuration_type": {
				Description: "The forest type used by the managed domain. One of `ResourceTrusting`, for a _Resource Forest_, or blank, for a _User Forest_.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"domain_name": {
				Description: "The Active Directory domain of the Domain Service. See [official documentation](https://learn.microsoft.com/azure/active-directory-domain-services/tutorial-create-instance#create-a-managed-domain) for constraints and recommendations.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
			},

			"filtered_sync_enabled": {
				Description: "Whether group-based filtered sync (also called scoped synchronisation) is enabled.",
				Type:        pluginsdk.TypeBool,
				Computed:    true,
			},

			"location": commonschema.LocationComputed(),

			"notifications": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"additional_recipients": {
							Description: "A list of additional email addresses to notify when there are alerts in the managed domain.",
							Type:        pluginsdk.TypeList,
							Computed:    true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
						},
						"notify_dc_admins": {
							Description: "Whether members of the _AAD DC Administrators_ group are notified when there are alerts in the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
						"notify_global_admins": {
							Description: "Whether all Global Administrators are notified when there are alerts in the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
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
							Description: "Whether secure LDAP is enabled for the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"external_access_enabled": {
							Description: "Whether external access to LDAPS over the Internet, is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
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
							Description: "Whether the Kerberos Armoring is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"kerberos_rc4_encryption_enabled": {
							Description: "Whether the Kerberos RC4 Encryption is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"ntlm_v1_enabled": {
							Description: "Whether legacy NTLM v1 support is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"sync_kerberos_passwords": {
							Description: "Whether Kerberos password hashes are synchronized to the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"sync_ntlm_passwords": {
							Description: "Whether NTLM password hashes are synchronized to the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"sync_on_prem_passwords": {
							Description: "Whether on-premises password hashes are synchronized to the managed domain.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},

						"tls_v1_enabled": {
							Description: "Whether legacy TLS v1 support is enabled.",
							Type:        pluginsdk.TypeBool,
							Computed:    true,
						},
					},
				},
			},

			"sku": {
				Description: "The SKU of the Domain Service resource.",
				Type:        pluginsdk.TypeString,
				Computed:    true,
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
			Description: "A list of subnet IP addresses for the domain controllers in the replica set, typically two.",
			Type:        pluginsdk.TypeList,
			Computed:    true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"external_access_ip_address": {
			Description: "The publicly routable IP address for the domain controllers in the replica set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"id": {
			Type:     pluginsdk.TypeString,
			Computed: true,
		},

		"location": commonschema.ResourceGroupNameForDataSource(),

		"service_status": {
			Description: "The current service status for the replica set.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},

		"subnet_id": {
			Description: "The ID of the subnet in which the replica set resides.",
			Type:        pluginsdk.TypeString,
			Computed:    true,
		},
	}
}

func dataSourceActiveDirectoryDomainServiceRead(d *pluginsdk.ResourceData, meta any) error {
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

		d.Set("domain_configuration_type", pointer.From(props.DomainConfigurationType))

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

		if err := d.Set("replica_sets", helpers.FlattenDomainServiceReplicaSets(props.ReplicaSets)); err != nil {
			return fmt.Errorf("setting `replica_sets`: %+v", err)
		}
	}

	return nil
}
