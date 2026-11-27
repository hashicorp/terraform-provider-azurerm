// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package api_management_openid_connect_provider

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/openidconnectprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	apimanagementClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementOpenIDConnectProviderCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := apimanagementClient.From(meta).OpenIdConnectClient
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := openidconnectprovider.NewOpenidConnectProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_openid_connect_provider", id.ID())
			}
		}
	}

	parameters := openidconnectprovider.OpenidConnectProviderContract{
		Properties: &openidconnectprovider.OpenidConnectProviderContractProperties{
			ClientId:         d.Get("client_id").(string),
			ClientSecret:     pointer.To(d.Get("client_secret").(string)),
			Description:      pointer.To(d.Get("description").(string)),
			DisplayName:      d.Get("display_name").(string),
			MetadataEndpoint: d.Get("metadata_endpoint").(string),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, openidconnectprovider.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementOpenIDConnectProviderRead(d, meta)
}
