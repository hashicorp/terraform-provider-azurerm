// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
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
				ValidateFunc: apimanagementservice.ValidateServiceID,
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

	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gatewayhostnameconfiguration.NewHostnameConfigurationID(apimId.SubscriptionId, apimId.ResourceGroupName, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("%s was not found", id)
		}

		return fmt.Errorf("making read request %s: %+v", id, err)
	}

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		_, err = gatewayhostnameconfiguration.ParseHostnameConfigurationID(*model.Id)
		if err != nil {
			return fmt.Errorf("parsing GatewayHostnameConfiguration ID %q", *model.Id)
		}
		if props := model.Properties; props != nil {
			d.Set("host_name", props.Hostname)
			d.Set("certificate_id", props.CertificateId)
			d.Set("request_client_certificate_enabled", props.NegotiateClientCertificate)
			d.Set("tls10_enabled", props.Tls10Enabled)
			d.Set("tls11_enabled", props.Tls11Enabled)
			d.Set("http2_enabled", props.HTTP2Enabled)
		}
	}

	d.SetId(id.ID())
	d.Set("api_management_id", apimId.ID())
	d.Set("gateway_name", id.GatewayId)

	return nil
}
