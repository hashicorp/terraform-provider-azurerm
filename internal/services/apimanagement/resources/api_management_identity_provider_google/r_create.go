// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_identity_provider_google

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/identityprovider"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementIdentityProviderGoogleCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.IdentityProviderClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := identityprovider.NewIdentityProviderID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), identityprovider.IdentityProviderTypeGoogle)
	clientID := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_identity_provider_google", id.ID())
			}
		}
	}

	parameters := identityprovider.IdentityProviderCreateContract{
		Properties: &identityprovider.IdentityProviderCreateContractProperties{
			ClientId:     clientID,
			ClientSecret: clientSecret,
			Type:         pointer.To(identityprovider.IdentityProviderTypeGoogle),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters, identityprovider.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementIdentityProviderGoogleRead(d, meta)
}
