package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2018-01-01/apimanagement"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/satori/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmApiManagementSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmApiManagementSubscriptionCreateUpdate,
		Read:   resourceArmApiManagementSubscriptionRead,
		Update: resourceArmApiManagementSubscriptionCreateUpdate,
		Delete: resourceArmApiManagementSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"subscription_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validate.UUIDOrEmpty,
			},

			"user_id": azure.SchemaApiManagementChildID(),

			"product_id": azure.SchemaApiManagementChildID(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(apimanagement.Submitted),
				ValidateFunc: validation.StringInSlice([]string{
					string(apimanagement.Active),
					string(apimanagement.Cancelled),
					string(apimanagement.Expired),
					string(apimanagement.Rejected),
					string(apimanagement.Submitted),
					string(apimanagement.Suspended),
				}, false),
			},

			"primary_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},

			"secondary_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceArmApiManagementSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.SubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = uuid.NewV4().String()
	}

	if features.ShouldResourcesBeImported() {
		resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_api_management_subscription", *resp.ID)
		}
	}

	displayName := d.Get("display_name").(string)
	productId := d.Get("product_id").(string)
	state := d.Get("state").(string)
	userId := d.Get("user_id").(string)

	params := apimanagement.SubscriptionCreateParameters{
		SubscriptionCreateParameterProperties: &apimanagement.SubscriptionCreateParameterProperties{
			DisplayName: utils.String(displayName),
			ProductID:   utils.String(productId),
			State:       apimanagement.SubscriptionState(state),
			UserID:      utils.String(userId),
		},
	}

	if v, ok := d.GetOk("primary_key"); ok {
		params.SubscriptionCreateParameterProperties.PrimaryKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("secondary_key"); ok {
		params.SubscriptionCreateParameterProperties.SecondaryKey = utils.String(v.(string))
	}

	sendEmail := utils.Bool(false)
	_, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, subscriptionId, params, sendEmail, "")
	if err != nil {
		return fmt.Errorf("Error creating/updating Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		return fmt.Errorf("Error retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourceArmApiManagementSubscriptionRead(d, meta)
}

func resourceArmApiManagementSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.SubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	subscriptionId := id.Path["subscriptions"]

	resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subscription %q was not found in API Management Service %q / Resource Group %q - removing from state!", subscriptionId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	d.Set("subscription_id", subscriptionId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.SubscriptionContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("primary_key", props.PrimaryKey)
		d.Set("secondary_key", props.SecondaryKey)
		d.Set("state", string(props.State))
		d.Set("product_id", props.ProductID)
		d.Set("user_id", props.UserID)
	}

	return nil
}

func resourceArmApiManagementSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).apiManagement.SubscriptionsClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.Path["service"]
	subscriptionId := id.Path["subscriptions"]

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, subscriptionId, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("Error removing Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
