// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package purview_account

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/purview/2021-12-01/account"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePurviewAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Purview.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := account.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_purview_account", id.ID())
		}
	}

	purviewAccount := account.Account{
		Properties: &account.AccountProperties{},
		Location:   pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:       tags.Expand(d.Get("tags").(map[string]any)),
	}

	expandedIdentity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	purviewAccount.Identity = expandedIdentity

	publicNetworkAccessEnabled := account.PublicNetworkAccessDisabled
	if d.Get("public_network_enabled").(bool) {
		publicNetworkAccessEnabled = account.PublicNetworkAccessEnabled
	}
	purviewAccount.Properties.PublicNetworkAccess = &publicNetworkAccessEnabled

	if v, ok := d.GetOk("managed_resource_group_name"); ok {
		purviewAccount.Properties.ManagedResourceGroupName = pointer.To(v.(string))
	}

	if d.Get("managed_event_hub_enabled").(bool) {
		purviewAccount.Properties.ManagedEventHubState = pointer.To(account.ManagedEventHubStateEnabled)
	} else {
		purviewAccount.Properties.ManagedEventHubState = pointer.To(account.ManagedEventHubStateDisabled)
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, purviewAccount, sdk.SetIDCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	d.SetId(id.ID())

	return resourcePurviewAccountRead(d, meta)
}
