package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/apimanagement/schemaz"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceApiManagementProduct() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementProductCreateUpdate,
		Read:   resourceApiManagementProductRead,
		Update: resourceApiManagementProductCreateUpdate,
		Delete: resourceApiManagementProductDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.ProductID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_required": func() *schema.Schema {
				if features.ThreePointOh() {
					return &schema.Schema{
						Type:     pluginsdk.TypeBool,
						Optional: true,
						Default:  true,
					}
				}

				return &schema.Schema{
					Type:     pluginsdk.TypeBool,
					Required: true,
				}
			}(),

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

	id := parse.NewProductID(subscriptionId, d.Get("resource_group_name").(string), d.Get("api_management_name").(string), d.Get("product_id").(string))

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	terms := d.Get("terms").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)
	approvalRequired := d.Get("approval_required").(bool)
	subscriptionsLimit := d.Get("subscriptions_limit").(int)
	published := d.Get("published").(bool)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_product", id.ID())
		}
	}
	publishedVal := apimanagement.ProductStateNotPublished
	if published {
		publishedVal = apimanagement.ProductStatePublished
	}

	properties := apimanagement.ProductContract{
		ProductContractProperties: &apimanagement.ProductContractProperties{
			Description:          utils.String(description),
			DisplayName:          utils.String(displayName),
			State:                publishedVal,
			SubscriptionRequired: utils.Bool(subscriptionRequired),
			Terms:                utils.String(terms),
		},
	}

	// Swagger says: Can be present only if subscriptionRequired property is present and has a value of false.
	// API/Portal says: Cannot provide values for approvalRequired and subscriptionsLimit when subscriptionRequired is set to false in the request payload
	if subscriptionRequired && subscriptionsLimit > 0 {
		properties.ProductContractProperties.ApprovalRequired = utils.Bool(approvalRequired)
		properties.ProductContractProperties.SubscriptionsLimit = utils.Int32(int32(subscriptionsLimit))
	} else if approvalRequired {
		return fmt.Errorf("`subscription_required` must be true and `subscriptions_limit` must be greater than 0 to use `approval_required`")
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.ServiceName, id.Name, properties, ""); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementProductRead(d, meta)
}

func resourceApiManagementProductRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("%s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on %s: %+v", *id, err)
	}

	d.Set("product_id", id.Name)
	d.Set("api_management_name", id.ServiceName)
	d.Set("resource_group_name", id.ResourceGroup)

	if props := resp.ProductContractProperties; props != nil {
		d.Set("approval_required", props.ApprovalRequired)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("published", props.State == apimanagement.ProductStatePublished)
		d.Set("subscriptions_limit", props.SubscriptionsLimit)
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("terms", props.Terms)
	}

	return nil
}

func resourceApiManagementProductDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting %s", *id)
	deleteSubscriptions := true
	resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, "", utils.Bool(deleteSubscriptions))
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting %s: %+v", *id, err)
		}
	}

	return nil
}
