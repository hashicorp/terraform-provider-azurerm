package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/gofrs/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementSubscription() *schema.Resource {
	return &schema.Resource{
		Create: resourceApiManagementSubscriptionCreateUpdate,
		Read:   resourceApiManagementSubscriptionRead,
		Update: resourceApiManagementSubscriptionCreateUpdate,
		Delete: resourceApiManagementSubscriptionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
			"user_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

			"api_management_name": schemaz.SchemaApiManagementName(),

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
				ValidateFunc: validate.ProductID,
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
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceGroup := d.Get("resource_group_name").(string)
	serviceName := d.Get("api_management_name").(string)
	subscriptionId := d.Get("subscription_id").(string)
	if subscriptionId == "" {
		subId, err := uuid.NewV4()
		if err != nil {
			return err
		}

		subscriptionId = subId.String()
	}

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
	allowTracing := d.Get("allow_tracing").(bool)

	params := apimanagement.SubscriptionCreateParameters{
		SubscriptionCreateParameterProperties: &apimanagement.SubscriptionCreateParameterProperties{
			DisplayName:  utils.String(displayName),
			Scope:        utils.String(productId),
			State:        apimanagement.SubscriptionState(state),
			AllowTracing: utils.Bool(allowTracing),
		},
	}
	if v, ok := d.GetOk("user_id"); ok {
		params.SubscriptionCreateParameterProperties.OwnerID = utils.String(v.(string))
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

	resp, err := client.Get(ctx, resourceGroup, serviceName, subscriptionId)
	if err != nil {
		return fmt.Errorf("retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", subscriptionId, serviceName, resourceGroup, err)
	}

	d.SetId(*resp.ID)

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

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subscription %q was not found in API Management Service %q / Resource Group %q - removing from state!", id.Name, id.ServiceName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Subscription %q (API Management Service %q / Resource Group %q): %+v", id.Name, id.ServiceName, id.ResourceGroup, err)
	}

	d.Set("subscription_id", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("api_management_name", id.ServiceName)

	if props := resp.SubscriptionContractProperties; props != nil {
		d.Set("display_name", props.DisplayName)
		d.Set("state", string(props.State))
		productId := ""
		if *props.Scope != "" {
			parseId, err := parse.ProductID(*props.Scope)
			if err != nil {
				return fmt.Errorf("parsing product id %q: %+v", *props.Scope, err)
			}
			productId = parseId.ID()
		}
		d.Set("product_id", productId)
		d.Set("user_id", props.OwnerID)
		d.Set("allow_tracing", props.AllowTracing)
	}

	// Primary and secondary keys must be got from this additional api
	keyResp, err := client.ListSecrets(ctx, id.ResourceGroup, id.ServiceName, id.Name)
	if err != nil {
		return fmt.Errorf("listing Subscription %q Primary and Secondary Keys (API Management Service %q / Resource Group %q): %+v", id.Name, id.ServiceName, id.ResourceGroup, err)
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

	if resp, err := client.Delete(ctx, id.ResourceGroup, id.ServiceName, id.Name, ""); err != nil {
		if !utils.ResponseWasNotFound(resp) {
			return fmt.Errorf("removing Subscription %q (API Management Service %q / Resource Group %q): %+v", id.Name, id.ServiceName, id.ResourceGroup, err)
		}
	}

	return nil
}
