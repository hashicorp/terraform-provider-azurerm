// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package data_share_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datashare/2019-11-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	clientsAccount "github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	datashareClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/datashare/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataShareAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := datashareClient.From(meta).AccountClient
	subscriptionId := clientsAccount.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(clientsAccount.StopContext(meta), d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}
		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_data_share_account", id.ID())
		}
	}

	expandedIdentity, err := identity.ExpandSystemAssigned(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	account := account.Account{
		Name:     pointer.To(id.AccountName),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Identity: *expandedIdentity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateCallbackThenPoll(ctx, id, account, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceDataShareAccountRead(d, meta)
}
