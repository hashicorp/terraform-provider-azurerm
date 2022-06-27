package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/automation/2021-06-22/automationaccount"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationAccountCreate,
		Read:   resourceAutomationAccountRead,
		Update: resourceAutomationAccountUpdate,
		Delete: resourceAutomationAccountDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := automationaccount.ParseAutomationAccountID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.AutomationAccount(),
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"sku_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice(automationaccount.PossibleValuesForSkuNameEnum(), false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": tags.Schema(),

			"dsc_server_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"dsc_primary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"dsc_secondary_access_key": {
				Type:      pluginsdk.TypeString,
				Computed:  true,
				Sensitive: true,
			},
			"public_network_access_enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceAutomationAccountCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := automationaccount.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_automation_account", id.ID())
	}

	identityVal, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters := automationaccount.AutomationAccountCreateOrUpdateParameters{
		Properties: &automationaccount.AutomationAccountCreateOrUpdateProperties{
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
	}
	// for create account do not set identity property (even TypeNone is not allowed), or api will response error
	if identityVal.Type != identity.TypeNone {
		parameters.Identity = identityVal
	}
	if tagsVal := expandTags(d.Get("tags").(map[string]interface{})); tagsVal != nil {
		parameters.Tags = &tagsVal
	}

	if _, err := client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}
	identity, err := identity.ExpandSystemAndUserAssignedMap(d.Get("identity").([]interface{}))
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters := automationaccount.AutomationAccountUpdateParameters{
		Properties: &automationaccount.AutomationAccountUpdateProperties{
			Sku: &automationaccount.Sku{
				Name: automationaccount.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: identity,
	}

	if tagsVal := expandTags(d.Get("tags").(map[string]interface{})); tagsVal != nil {
		parameters.Tags = &tagsVal
	}

	if _, err := client.Update(ctx, *id, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	registrationClient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := registrationClient.Get(ctx, id.ResourceGroupName, id.AutomationAccountName)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] Agent Registration Info for %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Info for %s: %+v", *id, err)
	}

	d.Set("name", id.AutomationAccountName)
	d.Set("resource_group_name", id.ResourceGroupName)

	d.Set("location", location.NormalizeNilable(resp.Model.Location))
	publicNetworkAccessEnabled := true
	if resp.Model == nil || resp.Model.Properties == nil {
		return fmt.Errorf("retrieving Automation Account got empty Model")
	}
	prop := resp.Model.Properties
	if prop.PublicNetworkAccess != nil {
		publicNetworkAccessEnabled = *prop.PublicNetworkAccess
	}
	d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
	skuName := ""
	if sku := prop.Sku; sku != nil {
		skuName = string(prop.Sku.Name)
	}
	d.Set("sku_name", skuName)

	d.Set("dsc_server_endpoint", keysResp.Endpoint)
	if keys := keysResp.Keys; keys != nil {
		d.Set("dsc_primary_access_key", keys.Primary)
		d.Set("dsc_secondary_access_key", keys.Secondary)
	}

	identity, err := identity.FlattenSystemAndUserAssignedMap(resp.Model.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if resp.Model.Tags != nil {
		return flattenAndSetTags(d, *resp.Model.Tags)
	}
	return nil
}

func resourceAutomationAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := automationaccount.ParseAutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}
