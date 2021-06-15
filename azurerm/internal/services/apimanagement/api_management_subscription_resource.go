package apimanagement

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2020-12-01/apimanagement"
	"github.com/gofrs/uuid"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/schemaz"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/apimanagement/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceApiManagementSubscription() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceApiManagementSubscriptionCreateUpdate,
		Read:   resourceApiManagementSubscriptionRead,
		Update: resourceApiManagementSubscriptionCreateUpdate,
		Delete: resourceApiManagementSubscriptionDelete,
		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

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
				ValidateFunc: validation.Any(validation.IsUUID, validation.StringIsEmpty),
			},

			// 3.0 this seems to have been renamed to owner id?
			"user_id": {
				Type:         pluginsdk.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"resource_group_name": azure.SchemaResourceGroupName(),

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
				ValidateFunc:  validate.ProductID,
				ConflictsWith: []string{"api_id"},
			},

			"api_id": {
				Type:          pluginsdk.TypeString,
				Optional:      true,
				ForceNew:      true,
				ValidateFunc:  validate.ApiID,
				ConflictsWith: []string{"product_id"},
			},

			"state": {
				Type:     pluginsdk.TypeString,
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
		scope = "all_apis"
	}

	params := apimanagement.SubscriptionCreateParameters{
		SubscriptionCreateParameterProperties: &apimanagement.SubscriptionCreateParameterProperties{
			DisplayName:  utils.String(displayName),
			Scope:        utils.String(scope),
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

func resourceApiManagementSubscriptionRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
		apiId := ""
		if *props.Scope != "" {
			// the scope is either a product or api id or "all_apis" constant
			parseId, err := parse.ProductID(*props.Scope)
			if err == nil {
				productId = parseId.ID()
			} else {
				parsedApiId, err := parse.ApiID(*props.Scope)
				if err != nil {
					return fmt.Errorf("parsing scope into product/ api id %q: %+v", *props.Scope, err)
				}
				apiId = parsedApiId.ID()
			}
		}
		d.Set("product_id", productId)
		d.Set("api_id", apiId)
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

func resourceApiManagementSubscriptionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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
