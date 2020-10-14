package aadmgmt

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/domainservices/mgmt/2017-06-01/aad"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmActiveDirectoryDomainService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceArmActiveDirectoryDomainServiceRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotWhiteSpace,
			},

			"location": azure.SchemaLocationForDataSource(),

			"resource_group_name": azure.SchemaResourceGroupNameForDataSource(),

			"domain_controller_ip_addresses": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"security": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ntlm_v1": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"sync_ntlm_passwords": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"tls_v1": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"filtered_sync": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"ldaps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"external_access": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ldaps": {
							Type:     schema.TypeBool,
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

			"notifications": {
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
							Type:     schema.TypeBool,
							Computed: true,
						},
						"notify_global_admins": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceArmActiveDirectoryDomainServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).AadMgmt.DomainServicesClient
	ctx := meta.(*clients.Client).StopContext

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
		if err := d.Set("domain_controller_ip_addresses", domainServiceProperties.DomainControllerIPAddress); err != nil {
			return fmt.Errorf("Error setting `domain_controller_ip_addresses`: %+v", err)
		}
		if err := d.Set("security", flattenDomainServiceSecurity(domainServiceProperties.DomainSecuritySettings)); err != nil {
			return fmt.Errorf("Error setting `security`: %+v", err)
		}
		if err := d.Set("ldaps", flattenDomainServiceLdapsSettings(domainServiceProperties.LdapsSettings)); err != nil {
			return fmt.Errorf("Error setting `ldaps_settings`: %+v", err)
		}
		if err := d.Set("notifications", flattenDomainServiceNotification(domainServiceProperties.NotificationSettings)); err != nil {
			return fmt.Errorf("Error setting `notification_settings`: %+v", err)
		}
		d.Set("filtered_sync", false)
		if domainServiceProperties.FilteredSync == aad.FilteredSyncEnabled {
			d.Set("filtered_sync", true)
		}
		d.Set("subnet_id", domainServiceProperties.SubnetID)
	}

	return nil
}