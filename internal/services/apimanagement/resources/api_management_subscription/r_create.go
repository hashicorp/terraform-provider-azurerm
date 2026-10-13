// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_subscription

import (
	"fmt"
	"net/http"

	"github.com/Azure/go-autorest/autorest"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	subName := d.Get("subscription_id").(string)
	if subName == "" {
		subId, err := uuid.NewV4()
		if err != nil {
			return err
		}

		subName = subId.String()
	}
	id := subscription.NewSubscriptions2ID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), subName)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			resp, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(resp.HttpResponse) {
					return fmt.Errorf("checking for present of existing %s: %+v", id, err)
				}
			}

			if !response.WasNotFound(resp.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_subscription", id.ID())
			}
		}
	}

	displayName := d.Get("display_name").(string)
	productId, productSet := d.GetOk("product_id")
	apiId, apiSet := d.GetOk("api_id")
	state := d.Get("state").(string)
	allowTracing := d.Get("allow_tracing").(bool)

	var scope string
	switch {
	case productSet:
		scope = productId.(string)
	case apiSet:
		scope = apiId.(string)
	default:
		scope = "/apis"
	}

	params := subscription.SubscriptionCreateParameters{
		Properties: &subscription.SubscriptionCreateParameterProperties{
			DisplayName:  displayName,
			Scope:        scope,
			State:        pointer.ToEnum[subscription.SubscriptionState](state),
			AllowTracing: pointer.To(allowTracing),
		},
	}
	if v, ok := d.GetOk("user_id"); ok {
		params.Properties.OwnerId = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("primary_key"); ok {
		params.Properties.PrimaryKey = pointer.To(v.(string))
	}

	if v, ok := d.GetOk("secondary_key"); ok {
		params.Properties.SecondaryKey = pointer.To(v.(string))
	}

	if err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, params, subscription.CreateOrUpdateOperationOptions{AppType: pointer.To(subscription.AppTypeDeveloperPortal), Notify: pointer.To(false)}); err != nil {
			// APIM admins set limit on number of subscriptions to a product.  In order to be able to correctly enforce that limit service cannot let simultaneous creations
			// to go through and first one wins/subsequent one gets 412 and that client/user can retry. This ensures that we have proper limits enforces as desired by APIM admin.
			if v, ok := err.(autorest.DetailedError); ok && v.StatusCode == http.StatusPreconditionFailed {
				return pluginsdk.RetryableError(err)
			}
			return pluginsdk.NonRetryableError(err)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementSubscriptionRead(d, meta)
}
