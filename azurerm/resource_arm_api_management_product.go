package azurerm

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementProduct() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementProductCreateUpdate,
		Read:   resourceArmApiManagementProductRead,
		Update: resourceArmApiManagementProductCreateUpdate,
		Delete: resourceArmApiManagementProductDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"product_id": azure.SchemaApiManagementChildName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
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

func resourceArmApiManagementProductCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, serviceName, productId)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Product %q (API Management Service %q / Resource Group %q): %s", productId, serviceName, resourceGroup, err)
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
		return fmt.Errorf("Error creating/updating Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	read, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		return fmt.Errorf("Error retrieving Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
	}

	if read.ID == nil {
		return fmt.Errorf("Cannot read ID for Product %q (API Management Service %q / Resource Group %q)", productId, serviceName, resourceGroup)
	}

	d.SetId(*read.ID)

	return resourceArmApiManagementProductRead(d, meta)
}

func resourceArmApiManagementProductRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}

	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	productId := id.Path["products"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, productId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Product %q was not found in API Management Service %q / Resource Group %q - removing from state!", productId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error making Read request on Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
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

func resourceArmApiManagementProductDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.ProductsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	productId := id.Path["products"]

	log.Printf("[DEBUG] Deleting Product %q (API Management Service %q / Resource Grouo %q)", productId, serviceName, resourceGroup)
	deleteSubscriptions := true
	resp, err := client.Delete(ctx, resourceGroup, serviceName, productId, "", utils.Bool(deleteSubscriptions))
	if err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error deleting Product %q (API Management Service %q / Resource Group %q): %+v", productId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
