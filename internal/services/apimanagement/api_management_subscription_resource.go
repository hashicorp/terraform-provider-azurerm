// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/api"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/product"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2022-08-01/subscription"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementSubscriptionCreateUpdate,
		Read:   resourceApiManagementSubscriptionRead,
		Update: resourceApiManagementSubscriptionCreateUpdate,
		Delete: resourceApiManagementSubscriptionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := subscription.ParseSubscriptions2ID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"subscription_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.Any(validate.ApiManagementChildName, validation.StringIsEmpty),
			},

			// 3.0 this seems to have been renamed to owner id?
			"user_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"product_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  product.ValidateProductID,
				ConflictsWith: []string{"api_id"},
			},

			"api_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  api.ValidateApiID,
				ConflictsWith: []string{"product_id"},
			},

			"state": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(subscription.SubscriptionStateSubmitted),
				ValidateFunc: validation.StringInSlice([]string{
					string(subscription.SubscriptionStateActive),
					string(subscription.SubscriptionStateCancelled),
					string(subscription.SubscriptionStateExpired),
					string(subscription.SubscriptionStateRejected),
					string(subscription.SubscriptionStateSubmitted),
					string(subscription.SubscriptionStateSuspended),
				}, false),
			},

			"primary_key": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      pluginsdk.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},

			"allow_tracing": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceApiManagementSubscriptionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
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
			State:        pointer.To(subscription.SubscriptionState(state)),
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

	err := pluginsdk.Retry(d.Timeout(pluginsdk.TimeoutCreate), func() *pluginsdk.RetryError {
		if _, err := client.CreateOrUpdate(ctx, id, params, subscription.CreateOrUpdateOperationOptions{AppType: pointer.To(subscription.AppTypeDeveloperPortal), Notify: pointer.To(false)}); err != nil {
			// APIM admins set limit on number of subscriptions to a product.  In order to be able to correctly enforce that limit service cannot let simultaneous creations
			// to go through and first one wins/subsequent one gets 412 and that client/user can retry. This ensures that we have proper limits enforces as desired by APIM admin.
			if v, ok := err.(autorest.DetailedError); ok && v.StatusCode == http.StatusPreconditionFailed {
				return pluginsdk.RetryableError(err)
			}
			return pluginsdk.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementSubscriptionRead(d, meta)
}

func resourceApiManagementSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

func resourceApiManagementSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := subscription.ParseSubscriptions2ID(d.Id())
	if err != nil {
		return err
	}

	if resp, err := client.Delete(ctx, *id, subscription.DeleteOperationOptions{}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
