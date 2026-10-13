// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package api_management_product

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProductCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := product.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	terms := d.Get("terms").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)
	approvalRequired := d.Get("approval_required").(bool)
	subscriptionsLimit := d.Get("subscriptions_limit").(int)
	published := d.Get("published").(bool)

	if d.IsNewResource() {
		if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.Get(ctx, id)
			if err != nil {
				if !response.WasNotFound(existing.HttpResponse) {
					return fmt.Errorf("checking for presence of existing %s: %s", id, err)
				}
			}

			if !response.WasNotFound(existing.HttpResponse) {
				return tf.ImportAsExistsError("azurerm_api_management_product", id.ID())
			}
		}
	}
	publishedVal := product.ProductStateNotPublished
	if published {
		publishedVal = product.ProductStatePublished
	}

	properties := product.ProductContract{
		Properties: &product.ProductContractProperties{
			Description:          pointer.To(description),
			DisplayName:          displayName,
			State:                pointer.To(publishedVal),
			SubscriptionRequired: pointer.To(subscriptionRequired),
			Terms:                pointer.To(terms),
		},
	}

	// Swagger says: Can be present only if subscriptionRequired property is present and has a value of false.
	// API/Portal says: Cannot provide values for approvalRequired and subscriptionsLimit when subscriptionRequired is set to false in the request payload
	if subscriptionRequired {
		if approvalRequired && subscriptionsLimit <= 0 {
			return fmt.Errorf("`subscriptions_limit` must be greater than 0 to use `approval_required`")
		}
		// Set `subscriptions_limit` to null or omit to allow unlimited per user subscriptions
		// When `subscriptions_limit` is specified as `0` it means the maximum number of subscriptions is 0, rather than allowing unlimited per user subscriptions
		if !pluginsdk.IsExplicitlyNullInConfig(d, "subscriptions_limit") && subscriptionsLimit >= 0 {
			properties.Properties.ApprovalRequired = pointer.To(approvalRequired)
			properties.Properties.SubscriptionsLimit = pointer.To(int64(subscriptionsLimit))
		}
	} else if approvalRequired {
		return fmt.Errorf("`subscription_required` must be true to use `approval_required`")
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties, product.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementProductRead(d, meta)
}
