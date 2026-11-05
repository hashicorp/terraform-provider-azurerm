// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package cognitive_account

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/cognitive/2026-03-01/cognitiveservicesaccounts"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/subnets"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceCognitiveAccountRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Cognitive.AccountsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := cognitiveservicesaccounts.ParseAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.AccountsGet(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	return resourceCognitiveAccountFlatten(ctx, client, d, id, resp.Model, true)
}

func serviceAssociationLinkStateRefreshFunc(ctx context.Context, client *subnets.SubnetsClient, subnetId commonids.SubnetId) pluginsdk.StateRefreshFunc {
	return func() (any, string, error) {
		resp, err := client.Get(ctx, subnetId, subnets.DefaultGetOperationOptions())
		if err != nil {
			return nil, "Error", fmt.Errorf("retrieving subnet %s: %+v", subnetId, err)
		}

		if resp.Model != nil && resp.Model.Properties != nil && resp.Model.Properties.ServiceAssociationLinks != nil {
			serviceAssociationLinks := *resp.Model.Properties.ServiceAssociationLinks

			if len(serviceAssociationLinks) > 0 {
				for _, sal := range serviceAssociationLinks {
					// The SAL name associated with the Cognitive Services account is "legionservicelink"
					if pointer.From(sal.Name) == "legionservicelink" && sal.Properties != nil && sal.Properties.ProvisioningState != nil {
						log.Printf("[DEBUG] Found Service Association Link %s with provisioning state %s",
							pointer.From(sal.Name), *sal.Properties.ProvisioningState)
						return resp, "SALExists", nil
					}
				}
			}
		}

		log.Printf("[DEBUG] No Service Association Links found on subnet %s", subnetId.ID())
		return resp, "SALDeleted", nil
	}
}
