package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/identity"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
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
			_, err := parse.AutomationAccountID(id)
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
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.SkuNameEnumBasic),
					string(automation.SkuNameEnumFree),
				}, false),
			},

			"identity": commonschema.SystemAssignedUserAssignedIdentityOptional(),

			"tags": tags.Schema(),

			"dsc_server_endpoint": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"dsc_primary_access_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
			"dsc_secondary_access_key": {
				Type:     pluginsdk.TypeString,
				Computed: true,
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

	id := parse.NewAutomationAccountID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_automation_account", id.ID())
	}

	identity, err := expandAutomationAccountIdentity(d.Get("identity").([]interface{}), true)
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters := automation.AccountCreateOrUpdateParameters{
		AccountCreateOrUpdateProperties: &automation.AccountCreateOrUpdateProperties{
			Sku: &automation.Sku{
				Name: automation.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationAccountID(d.Id())
	if err != nil {
		return err
	}
	identity, err := expandAutomationAccountIdentity(d.Get("identity").([]interface{}), false)
	if err != nil {
		return fmt.Errorf("expanding `identity`: %+v", err)
	}
	parameters := automation.AccountUpdateParameters{
		AccountUpdateProperties: &automation.AccountUpdateProperties{
			Sku: &automation.Sku{
				Name: automation.SkuNameEnum(d.Get("sku_name").(string)),
			},
			PublicNetworkAccess: utils.Bool(d.Get("public_network_access_enabled").(bool)),
		},
		Location: utils.String(location.Normalize(d.Get("location").(string))),
		Identity: identity,
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if _, err := client.Update(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	return resourceAutomationAccountRead(d, meta)
}

func resourceAutomationAccountRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	registrationClient := meta.(*clients.Client).Automation.AgentRegistrationInfoClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	keysResp, err := registrationClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Agent Registration Info for %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Registration Info for %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", location.NormalizeNilable(resp.Location))
	publicNetworkAccessEnabled := true
	if resp.PublicNetworkAccess != nil {
		publicNetworkAccessEnabled = *resp.PublicNetworkAccess
	}
	d.Set("public_network_access_enabled", publicNetworkAccessEnabled)
	skuName := ""
	if sku := resp.Sku; sku != nil {
		skuName = string(resp.Sku.Name)
	}
	d.Set("sku_name", skuName)

	d.Set("dsc_server_endpoint", keysResp.Endpoint)
	if keys := keysResp.Keys; keys != nil {
		d.Set("dsc_primary_access_key", keys.Primary)
		d.Set("dsc_secondary_access_key", keys.Secondary)
	}

	identity, err := flattenAutomationAccountIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("flattening `identity`: %+v", err)
	}
	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceAutomationAccountDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AutomationAccountID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAutomationAccountIdentity(input []interface{}, newResource bool) (*automation.Identity, error) {
	expanded, err := identity.ExpandSystemAndUserAssignedMap(input)
	if err != nil {
		return nil, err
	}

	if newResource && expanded.Type == identity.TypeNone {
		return nil, nil
	}

	out := automation.Identity{
		Type: automation.ResourceIdentityType(string(expanded.Type)),
	}

	if len(expanded.IdentityIds) > 0 {
		ids := make(map[string]*automation.IdentityUserAssignedIdentitiesValue)

		for k := range expanded.IdentityIds {
			ids[k] = &automation.IdentityUserAssignedIdentitiesValue{
				// intentionally empty
			}
		}

		out.UserAssignedIdentities = ids
	}

	return &out, nil
}

func flattenAutomationAccountIdentity(input *automation.Identity) (*[]interface{}, error) {
	var transformed *identity.SystemAndUserAssignedMap
	if input != nil {
		transformed = &identity.SystemAndUserAssignedMap{
			Type:        identity.Type(string(input.Type)),
			IdentityIds: make(map[string]identity.UserAssignedIdentityDetails),
		}
		if input.PrincipalID != nil {
			transformed.PrincipalId = *input.PrincipalID
		}
		if input.TenantID != nil {
			transformed.TenantId = *input.TenantID
		}
		if input.UserAssignedIdentities != nil {
			for k, v := range input.UserAssignedIdentities {
				transformed.IdentityIds[k] = identity.UserAssignedIdentityDetails{
					ClientId:    v.ClientID,
					PrincipalId: v.PrincipalID,
				}
			}
		}
	}

	return identity.FlattenSystemAndUserAssignedMap(transformed)
}
