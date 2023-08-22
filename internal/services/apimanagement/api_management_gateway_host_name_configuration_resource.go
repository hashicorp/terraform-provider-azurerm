// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/apimanagementservice"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gatewayhostnameconfiguration"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayHostNameConfiguration() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Read:   resourceApiManagementGatewayHostNameConfigurationRead,
		Update: resourceApiManagementGatewayHostNameConfigurationCreateUpdate,
		Delete: resourceApiManagementGatewayHostNameConfigurationDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gatewayhostnameconfiguration.ParseHostnameConfigurationID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": schemaz.SchemaApiManagementChildName(),

			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"gateway_name": schemaz.SchemaApiManagementChildName(),

			"certificate_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validate.CertificateID,
			},

			"host_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"request_client_certificate_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"http2_enabled": {
				Type:     pluginsdk.TypeBool,
				Default:  true,
				Optional: true,
			},

			"tls10_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"tls11_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementGatewayHostNameConfigurationCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := apimanagementservice.ParseServiceID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gatewayhostnameconfiguration.NewHostnameConfigurationID(apimId.SubscriptionId, apimId.ResourceGroupName, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_host_name_configuration", id.ID())
		}
	}

	parameters := gatewayhostnameconfiguration.GatewayHostnameConfigurationContract{
		Properties: &gatewayhostnameconfiguration.GatewayHostnameConfigurationContractProperties{
			Hostname:                   pointer.To(d.Get("host_name").(string)),
			CertificateId:              pointer.To(d.Get("certificate_id").(string)),
			NegotiateClientCertificate: pointer.To(d.Get("request_client_certificate_enabled").(bool)),
			Tls10Enabled:               pointer.To(d.Get("tls10_enabled").(bool)),
			Tls11Enabled:               pointer.To(d.Get("tls11_enabled").(bool)),
			HTTP2Enabled:               pointer.To(d.Get("http2_enabled").(bool)),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id, parameters, gatewayhostnameconfiguration.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayHostNameConfigurationRead(d, meta)
}

func resourceApiManagementGatewayHostNameConfigurationRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementGatewayHostNameConfigurationDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayHostNameConfigurationClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewayhostnameconfiguration.ParseHostnameConfigurationID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if resp, err := client.Delete(ctx, *id, gatewayhostnameconfiguration.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
