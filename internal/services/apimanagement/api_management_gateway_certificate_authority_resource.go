// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementGatewayCertificateAuthority() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementGatewayCertificateAuthorityCreateUpdate,
		Read:   resourceApiManagementGatewayCertificateAuthorityRead,
		Update: resourceApiManagementGatewayCertificateAuthorityCreateUpdate,
		Delete: resourceApiManagementGatewayCertificateAuthorityDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.GatewayCertificateAuthorityID(id)
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
				ValidateFunc: validate.ApiManagementID,
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

	id := parse.NewGatewayCertificateAuthorityID(apimId.SubscriptionId, apimId.ResourceGroup, apimId.ServiceName, d.Get("gateway_name").(string), d.Get("certificate_name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, d.Get("certificate_name").(string))
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_gateway_certificate_authority", id.ID())
		}
	}

	parameters := apimanagement.GatewayCertificateAuthorityContract{
		GatewayCertificateAuthorityContractProperties: &apimanagement.GatewayCertificateAuthorityContractProperties{
			IsTrusted: utils.Bool(d.Get("is_trusted").(bool)),
		},
	}

	_, err = client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, d.Get("gateway_name").(string), d.Get("certificate_name").(string), parameters, "")
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

	id, err := parse.GatewayCertificateAuthorityID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.CertificateAuthorityName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making read request for %s: %+v", *id, err)
	}

	apimId := parse.NewApiManagementID(id.SubscriptionId, id.ResourceGroup, id.ServiceName)

	d.Set("certificate_name", resp.Name)
	d.Set("api_management_id", apimId.ID())
	d.Set("gateway_name", id.GatewayName)

	if properties := resp.GatewayCertificateAuthorityContractProperties; properties != nil {
		d.Set("is_trusted", properties.IsTrusted)
	}

	return nil
}

func resourceApiManagementGatewayCertificateAuthorityDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.GatewayCertificateAuthorityClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.GatewayCertificateAuthorityID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.GatewayName, id.CertificateAuthorityName, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
