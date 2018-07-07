package azurerm

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceApiManagementService() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApiManagementRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"resource_group_name": resourceGroupNameForDataSourceSchema(),

			"location": locationForDataSourceSchema(),

			"publisher_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"publisher_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"sku": apiManagementDataSourceSkuSchema(),

			"notification_sender_email": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"gateway_regional_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"portal_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"management_api_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"scm_url": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vnet_subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"vnet_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"additional_location": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"location": locationForDataSourceSchema(),

						"sku": apiManagementDataSourceSkuSchema(),

						"vnet_subnet_id": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"gateway_regional_url": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"static_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"store_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"certificate_info": apiManagementDataSourceCertificateInfoSchema(),
					},
				},
			},

			"custom_properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},

			"hostname_configuration": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},

						"certificate_info": apiManagementDataSourceCertificateInfoSchema(),

						"default_ssl_binding": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"negotiate_client_certificate": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},

			"tags": tagsForDataSourceSchema(),
		},
	}
}

func dataSourceApiManagementRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagementServiceClient

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	ctx := meta.(*ArmClient).StopContext
	resp, err := client.Get(ctx, resourceGroup, name)

	if err != nil {
		return fmt.Errorf("Error making Read request on API Management Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if utils.ResponseWasNotFound(resp.Response) {
		return fmt.Errorf("Error: API Management Service %q (Resource Group %q) was not found", name, resourceGroup)
	}

	d.SetId(*resp.ID)

	d.Set("name", name)
	d.Set("resource_group_name", resourceGroup)

	if location := resp.Location; location != nil {
		d.Set("location", azureRMNormalizeLocation(*location))
	}

	if props := resp.ServiceProperties; props != nil {
		d.Set("publisher_email", props.PublisherEmail)
		d.Set("publisher_name", props.PublisherName)

		d.Set("notification_sender_email", props.NotificationSenderEmail)
		d.Set("created", props.CreatedAtUtc.Local().Format(time.RFC3339))
		d.Set("gateway_url", props.GatewayURL)
		d.Set("gateway_regional_url", props.GatewayRegionalURL)
		d.Set("portal_url", props.PortalURL)
		d.Set("management_api_url", props.ManagementAPIURL)
		d.Set("scm_url", props.ScmURL)
		d.Set("static_ips", props.StaticIps)
		d.Set("vnet_type", string(props.VirtualNetworkType))
		d.Set("custom_properties", &props.CustomProperties)

		err, hostnameConfigurations := flattenApiManagementHostnameConfigurations(d, props.HostnameConfigurations)
		if err != nil {
			return err
		}

		if err := d.Set("hostname_configuration", hostnameConfigurations); err != nil {
			return fmt.Errorf("Error setting `hostname_configuration`: %+v", err)
		}

		additionalLocations := flattenApiManagementAdditionalLocations(props.AdditionalLocations)
		if err := d.Set("additional_location", additionalLocations); err != nil {
			return fmt.Errorf("Error setting `additional_location`: %+v", err)
		}

		err, certificates := flattenApiManagementCertificates(d, props.Certificates)

		if err != nil {
			return err
		}

		if err := d.Set("certificate", certificates); err != nil {
			return fmt.Errorf("Error setting `certificate`: %+v", err)
		}

		if vnetConfig := props.VirtualNetworkConfiguration; vnetConfig != nil {
			if subnetId := vnetConfig.SubnetResourceID; subnetId != nil {
				d.Set("vnet_subnet_id", subnetId)
			}
		}
	}

	if sku := resp.Sku; sku != nil {
		d.Set("sku", flattenApiManagementServiceSku(sku))
	}

	flattenAndSetTags(d, resp.Tags)

	return nil
}

func apiManagementDataSourceSkuSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"capacity": {
					Type:     schema.TypeInt,
					Computed: true,
				},
			},
		},
	}
}

func apiManagementDataSourceCertificateInfoSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"expiry": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"thumbprint": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"subject": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}
