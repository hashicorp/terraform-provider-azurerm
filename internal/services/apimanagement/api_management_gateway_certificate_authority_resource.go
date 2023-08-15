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
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/gatewaycertificateauthority"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementGatewayCertificateAuthority() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayCertificateAuthorityCreateUpdate,
		Read:   resourceApiManagementGatewayCertificateAuthorityRead,
		Update: resourceApiManagementGatewayCertificateAuthorityCreateUpdate,
		Delete: resourceApiManagementGatewayCertificateAuthorityDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := gatewaycertificateauthority.ParseCertificateAuthorityID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"api_management_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: apimanagementservice.ValidateServiceID,
			},

			"certificate_name": schemaz.SchemaApiManagementChildName(),

			"gateway_name": schemaz.SchemaApiManagementChildName(),

			"is_trusted": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementGatewayCertificateAuthorityCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayCertificateAuthorityClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	apimId, err := parse.ApiManagementID(d.Get("api_management_id").(string))
	if err != nil {
		return fmt.Errorf("parsing `api_management_id`: %v", err)
	}

	id := gatewaycertificateauthority.NewCertificateAuthorityID(apimId.SubscriptionId, apimId.ResourceGroup, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("certificate_name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_certificate_authority", id.ID())
		}
	}

	parameters := gatewaycertificateauthority.GatewayCertificateAuthorityContract{
		Properties: &gatewaycertificateauthority.GatewayCertificateAuthorityContractProperties{
			IsTrusted: pointer.To(d.Get("is_trusted").(bool)),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id, parameters, gatewaycertificateauthority.CreateOrUpdateOperationOptions{})
	if err != nil {
		return fmt.Errorf("creating or updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementGatewayCertificateAuthorityRead(d, meta)
}

func resourceApiManagementGatewayCertificateAuthorityRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayCertificateAuthorityClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewaycertificateauthority.ParseCertificateAuthorityID(d.Id())
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
		d.Set("certificate_name", pointer.From(model.Name))
		if props := model.Properties; props != nil {
			d.Set("is_trusted", pointer.From(props.IsTrusted))
		}
	}

	return nil
}

func resourceApiManagementGatewayCertificateAuthorityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayCertificateAuthorityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := gatewaycertificateauthority.ParseCertificateAuthorityID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if resp, err := client.Delete(ctx, *id, gatewaycertificateauthority.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
