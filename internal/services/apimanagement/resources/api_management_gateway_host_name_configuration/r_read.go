// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_gateway_host_name_configuration

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2024-05-01/apimanagementservice"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayHostNameConfigurationRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayhostnameconfiguration.ParseHostnameConfigurationID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", *id, err)
	}

	apimId := apimanagementservice.NewServiceID(id.SubscriptionId, id.ResourceGroupName, id.ServiceName)

	d.Set("api_management_id", apimId.ID())
	d.Set("gateway_name", id.GatewayId)

	if model := resp.Model; model != nil {
		d.Set("name", pointer.From(model.Name))
		if properties := model.Properties; properties != nil {
			d.Set("host_name", pointer.From(properties.Hostname))
			d.Set("certificate_id", pointer.From(properties.CertificateId))
			d.Set("request_client_certificate_enabled", pointer.From(properties.NegotiateClientCertificate))
			d.Set("tls10_enabled", pointer.From(properties.Tls10Enabled))
			d.Set("tls11_enabled", pointer.From(properties.Tls11Enabled))
			d.Set("http2_enabled", pointer.From(properties.HTTP2Enabled))
		}
	}

	return nil
}
