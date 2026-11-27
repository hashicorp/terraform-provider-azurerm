// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package spring_cloud_api_portal_custom_domain

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	springcloudClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/springcloud/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/appplatform/2023-05-01-preview/appplatform"
)

func resourceSpringCloudAPIPortalCustomDomainCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	subscriptionId := account.From(meta).SubscriptionId
	client := springcloudClient.From(meta).APIPortalCustomDomainClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	// the ID is parsed insensitively to preserve the behaviour of the legacy parser this replaced -
	// we are ok with this remaining insensitive as Azure Spring Apps is deprecated and will be removed
	portalId, err := appplatform_rm.ParseApiPortalIDInsensitively(d.Get("spring_cloud_api_portal_id").(string))
	if err != nil {
		return err
	}
	id := appplatform_rm.NewApiPortalDomainID(subscriptionId, portalId.ResourceGroupName, portalId.SpringName, portalId.ApiPortalName, d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id.ResourceGroupName, id.SpringName, id.ApiPortalName, id.DomainName)
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_spring_cloud_api_portal_custom_domain", id.ID())
			}
		}
	}

	apiPortalCustomDomainResource := appplatform.APIPortalCustomDomainResource{
		Properties: &appplatform.APIPortalCustomDomainProperties{
			Thumbprint: pointer.To(d.Get("thumbprint").(string)),
		},
	}
	future, err := client.CreateOrUpdate(ctx, id.ResourceGroupName, id.SpringName, id.ApiPortalName, id.DomainName, apiPortalCustomDomainResource)
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation/update of %s: %+v", id, err)
	}

	return resourceSpringCloudAPIPortalCustomDomainRead(d, meta)
}
