package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementProduct() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementProductCreateUpdate,
		Read:   resourceApiManagementProductRead,
		Update: resourceApiManagementProductCreateUpdate,
		Delete: resourceApiManagementProductDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"product_id": schemaz.SchemaApiManagementChildName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"subscription_required": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"published": {
				Type:     schema.TypeBool,
				Required: true,
			},

			"approval_required": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"terms": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"subscriptions_limit": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
}

func resourceApiManagementProductCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	log.Printf("[INFO] preparing arguments for API Management Product creation.")

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	productId := d.Get("product_id").(string)

	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	terms := d.Get("terms").(string)
	subscriptionRequired := d.Get("subscription_required").(bool)
	approvalRequired := d.Get("approval_required").(bool)
	subscriptionsLimit := d.Get("subscriptions_limit").(int)
	published := d.Get("published").(bool)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, productId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Product %q (API Management Service %q / Resource Group %q): %s", productId, serviceName, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_api_management_product", *existing.ID)
		}
	}
	publishedVal := apimanagement.NotPublished
	if published {
		publishedVal = apimanagement.Published
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

	if _, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, productId, properties, ""); err != nil {
		return fmt.Errorf("creating/updating Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		return fmt.Errorf("retrieving Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	if resp.ID == nil {
		return fmt.Errorf("Cannot read ID for Product %q (API Management Service %q / Resource Group %q)", productId, serviceName, resourceGroup)
	}

	d.SetId(*resp.ID)

	return resourceApiManagementProductRead(d, meta)
}

func resourceApiManagementProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	productId := id.Name

	resp, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Product %q was not found in API Management Service %q / Resource Group %q - removing from state!", productId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	d.Set("product_id", productId)
	d.Set("api_management_name", serviceName)
	d.Set("resource_group_name", resourceGroup)

	if props := resp.ProductContractProperties; props != nil {
		d.Set("approval_required", props.ApprovalRequired)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("published", props.State == apimanagement.Published)
		d.Set("subscriptions_limit", props.SubscriptionsLimit)
		d.Set("subscription_required", props.SubscriptionRequired)
		d.Set("terms", props.Terms)
	}

	return nil
}

func resourceApiManagementProductDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.ProductID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	productId := id.Name

	log.Printf("[DEBUG] Deleting Product %q (API Management Service %q / Resource Grouo %q)", productId, serviceName, resourceGroup)
	deleteSubscriptions := true
	resp, err := client.Delete(ctx, resourceGroup, serviceName, productId, "", utils.Bool(deleteSubscriptions))
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("deleting Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
