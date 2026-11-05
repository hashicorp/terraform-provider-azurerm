// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"fmt"
	"slices"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
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

	props := cognitiveservicesaccounts.Account{
		Sku:  &sku,
		Kind: pointer.To(d.Get("kind").(string)),
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
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if d.HasChanges("customer_managed_key") {
		old, new := d.GetChange("customer_managed_key")
		// Remove `customer_managed_key` (switch using a customer managed key to microsoft managed), and explicitly specify KeySource as `Microsoft.CognitiveServices`.
		if len(old.([]any)) > 0 && len(new.([]any)) == 0 {
			props.Properties.Encryption = &cognitiveservicesaccounts.Encryption{
				KeySource: pointer.To(cognitiveservicesaccounts.KeySourceMicrosoftPointCognitiveServices),
			}
		}
	}

	if d.HasChange("identity") {
		identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]any))
		if err != nil {
			return fmt.Errorf("expanding `identity`: %+v", err)
		}
		props.Identity = identity
	}

	if err = client.AccountsUpdateThenPoll(ctx, *id, props); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceCognitiveAccountRead(d, meta)
}
