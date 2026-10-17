// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_subscription

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementSubscriptionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscription.ParseSubscriptions2ID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("subscription_id", id.SubscriptionName)
	d.Set("resource_group_name", id.ResourceGroupName)
	d.Set("api_management_name", id.ServiceName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("display_name", pointer.From(props.DisplayName))
			d.Set("state", string(props.State))
			productId := ""
			apiId := ""
			// check if the subscription is for all apis or a specific product/ api
			if props.Scope != "" && !strings.HasSuffix(props.Scope, "/apis") {
				// the scope is either a product or api id
				parseId, err := product.ParseProductIDInsensitively(props.Scope)
				if err == nil {
					productId = parseId.ID()
				} else {
					parsedApiId, err := api.ParseApiIDInsensitively(props.Scope)
					if err != nil {
						return fmt.Errorf("parsing scope into product/ api id %q: %+v", props.Scope, err)
					}
					apiId = parsedApiId.ID()
				}
			}
			d.Set("product_id", productId)
			d.Set("api_id", apiId)
			d.Set("user_id", pointer.From(props.OwnerId))
			d.Set("allow_tracing", pointer.From(props.AllowTracing))
		}
	}

	// Primary and secondary keys must be got from this additional api
	keyResp, err := client.ListSecrets(ctx, *id)
	if err != nil {
		return fmt.Errorf("listing Subscription %q Primary and Secondary Keys (API Management Service %q / Resource Group %q): %+v", id.SubscriptionId, id.ServiceName, id.ResourceGroupName, err)
	}
	if model := keyResp.Model; model != nil {
		d.Set("primary_key", pointer.From(model.PrimaryKey))
		d.Set("secondary_key", pointer.From(model.SecondaryKey))
	}

	return nil
}
