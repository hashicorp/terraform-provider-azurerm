package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2019-12-01/apimanagement"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/satori/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementSubscriptionCreateUpdate,
		Read:   resourceApiManagementSubscriptionRead,
		Update: resourceApiManagementSubscriptionCreateUpdate,
		Delete: resourceApiManagementSubscriptionDelete,
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
			"subscription_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
			},

			// 3.0 this seems to have been renamed to owner id?
			"user_id": azure.SchemaApiManagementChildID(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": azure.SchemaApiManagementName(),

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			// TODO this now sets the scope property - either a scope block needs adding or additional properties `api_id` and maybe `all_apis`
			"product_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
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

			"allow_tracing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceApiManagementSubscriptionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	azureSubscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subscriptionId = uuid.NewV4().String()
	}

	id := parse.NewSubscriptionID(azureSubscriptionId, resourceGroup, serviceName, subscriptionId)

	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for present of existing Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
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
	allowTracing := d.Get("allow_tracing").(bool)

	params := apimanagement.SubscriptionCreateParameters{
		SubscriptionCreateParameterProperties: &apimanagement.SubscriptionCreateParameterProperties{
			DisplayName:  utils.String(displayName),
			Scope:        utils.String(productId),
			State:        apimanagement.SubscriptionState(state),
			OwnerID:      utils.String(userId),
			AllowTracing: utils.Bool(allowTracing),
		},
	}

	if v, ok := d.GetOk("primary_key"); ok {
		params.SubscriptionCreateParameterProperties.PrimaryKey = utils.String(v.(string))
	}

	if v, ok := d.GetOk("secondary_key"); ok {
		params.SubscriptionCreateParameterProperties.SecondaryKey = utils.String(v.(string))
	}

	sendEmail := utils.Bool(false)
	_, err := client.CreateOrUpdate(ctx, resourceGroup, serviceName, subscriptionId, params, sendEmail, "", apimanagement.DeveloperPortal)
	if err != nil {
		return fmt.Errorf("creating/updating Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	_, err = client.Get(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		return fmt.Errorf("retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	d.SetId(id.ID())

	return resourceApiManagementSubscriptionRead(d, meta)
}

func resourceApiManagementSubscriptionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	subscriptionId := id.SubscriptionId

	resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subscription %q was not found in API Management Service %q / Resource Group %q - removing from state!", subscriptionId, serviceName, resourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	d.Set("subscription_id", subscriptionId)
	d.Set("resource_group_name", resourceGroup)
	d.Set("api_management_name", serviceName)

	if props := resp.SubscriptionContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("state", string(props.State))
		d.Set("product_id", props.Scope)
		d.Set("user_id", props.OwnerID)
		d.Set("allow_tracing", props.AllowTracing)
	}

	// Primary and secondary keys must be got from this additional api
	keyResp, err := client.ListSecrets(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		return fmt.Errorf("listing Subscription %q Primary and Secondary Keys (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}
	d.Set("primary_key", keyResp.PrimaryKey)
	d.Set("secondary_key", keyResp.SecondaryKey)

	return nil
}

func resourceApiManagementSubscriptionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).ApiManagement.SubscriptionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubscriptionID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	serviceName := id.ServiceName
	subscriptionId := id.SubscriptionId

	if resp, err := client.Delete(ctx, resourceGroup, serviceName, subscriptionId, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
		}
	}

	return nil
}
