package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceApiManagementGatewayHostnameConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Read: dataSourceApiManagementGatewayHostnameConfigurationRead,

		Timeouts: &pluginsdk.ResourceTimeout{
			Read: pluginsdk.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"api_management_gateway_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.GatewayID,
			},
			"hostname": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"certificate_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementGatewayHostnameConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostnameConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	gwID, err := parse.GatewayID(d.Get("api_management_gateway_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_gateway_id`: %v", err)
	}

	id := parse.NewGatewayHostnameConfigurationID(gwID.SubscriptionId, gwID.ResourceGroup, gwID.ServiceName, gwID.Name, d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	_, err = parse.GatewayHostnameConfigurationID(*resp.ID)
	if err != nil {
		return fmt.Errorf("parsing GatewayHostnameConfiguration ID %q", *resp.ID)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("api_management_gateway_id", gwID.ID())

	if properties := resp.GatewayHostnameConfigurationContractProperties; properties != nil {
		d.Set("hostname", properties.Hostname)
		d.Set("certificate_id", properties.CertificateID)
	}

	return nil
}
