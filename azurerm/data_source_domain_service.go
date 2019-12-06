package azurerm

import (
	"fmt"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmDomainService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmDomainServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"domain_controller_ip_address": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"domain_security_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ntlm_v1": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sync_ntlm_passwords": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls_v1": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"filtered_sync": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"ldaps_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_access": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ldaps": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pfx_certificate": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pfx_certificate_password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_access_ip_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"notification_settings": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_recipients": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"notify_dc_admins": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"notify_global_admins": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func dataSourceArmDomainServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).DomainServices.DomainServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*ArmClient).StopContext, d)
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

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if domainServiceProperties := resp.DomainServiceProperties; domainServiceProperties != nil {
		d.Set("domain_controller_ip_address", utils.FlattenStringSlice(domainServiceProperties.DomainControllerIPAddress))
		if err := d.Set("domain_security_settings", flattenArmDomainServiceDomainSecuritySettings(domainServiceProperties.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("Error setting `domain_security_settings`: %+v", err)
		}
		d.Set("filtered_sync", string(domainServiceProperties.FilteredSync))
		if err := d.Set("ldaps_settings", flattenArmDomainServiceLdapsSettings(domainServiceProperties.LdapsSettings)); err != nil {
			return fmt.Errorf("Error setting `ldaps_settings`: %+v", err)
		}
		if err := d.Set("notification_settings", flattenArmDomainServiceNotificationSettings(domainServiceProperties.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notification_settings`: %+v", err)
		}
		d.Set("subnet_id", domainServiceProperties.SubnetID)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}
