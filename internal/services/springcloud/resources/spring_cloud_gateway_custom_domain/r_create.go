// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_gateway_custom_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudGatewayCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	client := meta.(*clients.Client).AppPlatform.GatewayCustomDomainClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	gatewayId, err := appplatform_rm.ParseGatewayIDInsensitively(d.Get("spring_cloud_gateway_id").(string))
	if err != nil {
		return err
	}
	id := appplatform_rm.NewGatewayDomainID(subscriptionId, gatewayId.ResourceGroupName, gatewayId.SpringName, gatewayId.GatewayName, d.Get("name").(string))

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.DomainName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_gateway_custom_domain", id.ID())
			}
		}
	}

	gatewayCustomDomainResource := appplatform.GatewayCustomDomainResource{
		Properties: &appplatform.GatewayCustomDomainProperties{
			Thumbprint: pointer.To(d.Get("thumbprint").(string)),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.GatewayName, id.DomainName, gatewayCustomDomainResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	return resourceSpringCloudGatewayCustomDomainRead(d, meta)
}
