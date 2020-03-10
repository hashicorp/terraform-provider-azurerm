package web

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func dataSourceArmAppServiceVirtualNetworkConnectionGateway() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAppServiceVirtualNetworkConnectionGatewayRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"app_service_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateAppServiceName,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"virtual_network_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"certificate_blob": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"certificate_thumbprint": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"dns_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"resync_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"virtual_network_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceAppServiceVirtualNetworkConnectionGatewayRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Web.AppServicesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resGroup := d.Get("resource_group_name").(string)
	appServiceName := d.Get("app_service_name").(string)
	vnetName := d.Get("virtual_network_name").(string)

	resp, err := client.GetVnetConnection(ctx, resGroup, appServiceName, vnetName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] App Service Virtual Network Connection for app %q vnet %q was not found in Resource Group %q - removnig from state!", appServiceName, vnetName, resGroup)
			return nil
		}

		return fmt.Errorf("Error making Read request on App Service Virtual Network Connection for app %q vnet %q (Resource Group %q): %+v", appServiceName, vnetName, resGroup, err)
	}

	if resp.ID != nil || *resp.ID != "" {
		d.SetId(*resp.ID)
	}
	d.Set("name", resp.Name)
	d.Set("resource_group_name", resGroup)
	d.Set("app_service_name", appServiceName)
	if props := resp.VnetInfoProperties; props != nil {
		d.Set("virtual_network_id", props.VnetResourceID)
		d.Set("certificate_thumbprint", props.CertThumbprint)
		d.Set("certificate_blob", props.CertBlob)
		d.Set("resync_required", props.ResyncRequired)
		if props.DNSServers != nil {
			d.Set("dns_servers", strings.Split(*props.DNSServers, ","))
		} else {
			d.Set("dns_servers", []string{})
		}
	}

	return nil
}
