package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func dataSourceApiManagementGatewayHostNameConfiguration() *pluginsdk.Resource {
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

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.ApiManagementID,
			},

			"gateway_name": schemaz.SchemaApiManagementChildDataSourceName(),

			"certificate_id": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"host_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"request_client_certificate_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"http2_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tls10_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},

			"tls11_enabled": {
				Type:     pluginsdk.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceApiManagementGatewayHostnameConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := parse.ApiManagementID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := parse.NewGatewayHostNameConfigurationID(apimId.SubscriptionId, apimId.ResourceGroup, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.HostnameConfigurationName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	_, err = parse.GatewayHostNameConfigurationID(*resp.ID)
	if err != nil {
		return fmt.Errorf("parsing GatewayHostnameConfiguration ID %q", *resp.ID)
	}

	d.SetId(id.ID())

	d.Set("name", resp.Name)
	d.Set("api_management_id", apimId.ID())
	d.Set("gateway_name", id.GatewayName)

	if properties := resp.GatewayHostnameConfigurationContractProperties; properties != nil {
		d.Set("host_name", properties.Hostname)
		d.Set("certificate_id", properties.CertificateID)
		d.Set("request_client_certificate_enabled", properties.NegotiateClientCertificate)
		d.Set("tls10_enabled", properties.TLS10Enabled)
		d.Set("tls11_enabled", properties.TLS11Enabled)
		d.Set("http2_enabled", properties.HTTP2Enabled)
	}

	return nil
}
