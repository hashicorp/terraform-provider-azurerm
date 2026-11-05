// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := cognitiveservicesaccounts.NewAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		existing, err := client.AccountsGet(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError(azureCognitiveAccountResourceName, id.ID())
		}
	}

	sku := cognitiveservicesaccounts.Sku{
		Name: d.Get("sku_name").(string),
	}

	networkAcls, subnetIds := expandCognitiveAccountNetworkAcls(d)

	// also lock on the Virtual Network ID's since modifications in the networking stack are exclusive
	virtualNetworkNames := make([]string, 0)
	for _, v := range subnetIds {
		id, err := commonids.ParseSubnetIDInsensitively(v)
		if err != nil {
			return err
		}
		if !slices.Contains(virtualNetworkNames, id.VirtualNetworkName) {
			virtualNetworkNames = append(virtualNetworkNames, id.VirtualNetworkName)
		}
	}

	locks.MultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(&virtualNetworkNames, network.VirtualNetworkResourceName)

	publicNetworkAccess := cognitiveservicesaccounts.PublicNetworkAccessEnabled
	if !d.Get("public_network_access_enabled").(bool) {
		publicNetworkAccess = cognitiveservicesaccounts.PublicNetworkAccessDisabled
	}

	apiProps, err := expandCognitiveAccountAPIProperties(d)
	if err != nil {
		return err
	}

	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}

	props := cognitiveservicesaccounts.Account{
		Kind:     pointer.To(d.Get("kind").(string)),
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Sku:      &sku,
		Properties: &cognitiveservicesaccounts.AccountProperties{
			ApiProperties:                 apiProps,
			NetworkAcls:                   networkAcls,
			CustomSubDomainName:           pointer.To(d.Get("custom_subdomain_name").(string)),
			AllowedFqdnList:               pluginsdk.ExpandStringSlice(d.Get("fqdns").([]any)),
			PublicNetworkAccess:           &publicNetworkAccess,
			UserOwnedStorage:              expandCognitiveAccountStorage(d.Get("storage").([]any)),
			RestrictOutboundNetworkAccess: pointer.To(d.Get("outbound_network_access_restricted").(bool)),
			DisableLocalAuth:              pointer.To(!d.Get("local_auth_enabled").(bool)),
			DynamicThrottlingEnabled:      pointer.To(d.Get("dynamic_throttling_enabled").(bool)),
			AllowProjectManagement:        pointer.To(d.Get("project_management_enabled").(bool)),
			NetworkInjections:             expandCognitiveAccountNetworkInjection(d.Get("network_injection").([]any)),
			Encryption:                    expandCognitiveAccountCustomerManagedKey(d.Get("customer_managed_key").([]any)),
		},
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.AccountsCreateCallbackThenPoll(ctx, id, props, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceCognitiveAccountRead(d, meta)
}
