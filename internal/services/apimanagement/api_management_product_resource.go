// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/apimanagement/2021-08-01/product"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceApiManagementProduct() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductCreateUpdate,
		Read:   resourceApiManagementProductRead,
		Update: resourceApiManagementProductCreateUpdate,
		Delete: resourceApiManagementProductDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := product.ParseProductID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"product_id": schemaz.SchemaApiManagementChildName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_required": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"published": {
				Type:     pluginsdk.TypeBool,
				Required: true,
			},

			"approval_required": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"terms": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"subscriptions_limit": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementProductCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management Product creation.")

	id := product.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	terms := d.Get("terms").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)
	approvalRequired := d.Get("approval_required").(bool)
	subscriptionsLimit := d.Get("subscriptions_limit").(int)
	published := d.Get("published").(bool)

	if d.IsNewResource() {
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
	if subscriptionRequired && subscriptionsLimit > 0 {
		properties.Properties.ApprovalRequired = pointer.To(approvalRequired)
		properties.Properties.SubscriptionsLimit = pointer.To(int64(subscriptionsLimit))
	} else if approvalRequired {
		return fmt.Errorf("`subscription_required` must be true and `subscriptions_limit` must be greater than 0 to use `approval_required`")
	}

	if _, err := client.CreateOrUpdate(ctx, id, properties, product.CreateOrUpdateOperationOptions{}); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementProductRead(d, meta)
}

func resourceApiManagementProductRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := product.ParseProductID(d.Id())
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

	d.Set("product_id", id.ProductId)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("approval_required", pointer.From(props.ApprovalRequired))
			d.Set("description", pointer.From(props.Description))
			d.Set("display_name", props.DisplayName)
			d.Set("published", pointer.From(props.State) == product.ProductStatePublished)
			d.Set("subscriptions_limit", pointer.From(props.SubscriptionsLimit))
			d.Set("subscription_required", pointer.From(props.SubscriptionRequired))
			d.Set("terms", pointer.From(props.Terms))
		}
	}

	return nil
}

func resourceApiManagementProductDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := product.ParseProductID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	if resp, err := client.Delete(ctx, *id, product.DeleteOperationOptions{DeleteSubscriptions: pointer.To(true)}); err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
