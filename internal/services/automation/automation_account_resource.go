package automation

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/automation/mgmt/2020-01-13-preview/automation"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/automation/validate"
	msiparse "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/parse"
	msivalidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/msi/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAutomationAccount() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAutomationAccountCreateUpdate,
		Read:   resourceAutomationAccountRead,
		Update: resourceAutomationAccountCreateUpdate,
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

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"sku_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(automation.SkuNameEnumBasic),
					string(automation.SkuNameEnumFree),
				}, false),
			},

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

			"identity": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(automation.ResourceIdentityTypeSystemAssigned),
								string(automation.ResourceIdentityTypeUserAssigned),
							}, false),
						},
						"user_assigned_identity_id": {
							Type:         pluginsdk.TypeString,
							ValidateFunc: msivalidate.UserAssignedIdentityID,
							Optional:     true,
						},
						"principal_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func expandAutomationAccountIdentity(input []interface{}) *automation.Identity {
	if len(input) == 0 || input[0] == nil {
		return &automation.Identity{
			Type: automation.ResourceIdentityTypeNone,
		}
	}

	values := input[0].(map[string]interface{})

	if automation.ResourceIdentityType(values["type"].(string)) == automation.ResourceIdentityTypeUserAssigned {
		userAssignedIdentities := map[string]*automation.IdentityUserAssignedIdentitiesValue{
			values["user_assigned_identity_id"].(string): {},
		}

		return &automation.Identity{
			Type:                   automation.ResourceIdentityType(values["type"].(string)),
			UserAssignedIdentities: userAssignedIdentities,
		}
	}

	return &automation.Identity{
		Type: automation.ResourceIdentityType(values["type"].(string)),
	}
}

func flattenAutomationAccountIdentity(input *automation.Identity) ([]interface{}, error) {
	// if it's none, omit the block
	if input == nil || input.Type == automation.ResourceIdentityTypeNone {
		return []interface{}{}, nil
	}

	identity := make(map[string]interface{})

	identity["principal_id"] = ""
	if input.PrincipalID != nil {
		identity["principal_id"] = *input.PrincipalID
	}

	identity["tenant_id"] = ""
	if input.TenantID != nil {
		identity["tenant_id"] = *input.TenantID
	}

	identity["user_assigned_identity_id"] = ""
	if input.UserAssignedIdentities != nil {
		keys := []string{}
		for key := range input.UserAssignedIdentities {
			keys = append(keys, key)
		}
		if len(keys) > 0 {
			parsedId, err := msiparse.UserAssignedIdentityIDInsensitively(keys[0])
			if err != nil {
				return nil, err
			}
			identity["user_assigned_identity_id"] = parsedId.ID()
		}
	}

	identity["type"] = string(input.Type)

	return []interface{}{identity}, nil
}

func resourceAutomationAccountCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Automation.AccountClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	sku := automation.Sku{
		Name: automation.SkuNameEnum(d.Get("sku_name").(string)),
	}

	log.Printf("[INFO] preparing arguments for Automation Account create/update.")

	id := parse.NewAutomationAccountID(client.SubscriptionID, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}

		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_automation_account", id.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	t := d.Get("tags").(map[string]interface{})

	parameters := automation.AccountCreateOrUpdateParameters{
		AccountCreateOrUpdateProperties: &automation.AccountCreateOrUpdateProperties{
			Sku: &sku,
		},
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}

	automationAccountIdentityRaw := d.Get("identity").([]interface{})
	if len(automationAccountIdentityRaw) > 0 {
		parameters.Identity = expandAutomationAccountIdentity(automationAccountIdentityRaw)
	} else {
		parameters.Identity = nil
	}

	if _, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

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
			log.Printf("[DEBUG] Automation Account %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request on Automation Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	keysResp, err := registrationClient.Get(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Agent Registration Info for Automation Account %q was not found in Resource Group %q - removing from state!", id.Name, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("making Read request for Agent Registration Info for Automation Account %q (Resource Group %q): %+v", id.Name, id.ResourceGroup, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if sku := resp.Sku; sku != nil {
		if err := d.Set("sku_name", string(sku.Name)); err != nil {
			return fmt.Errorf("setting 'sku_name': %+v", err)
		}
	} else {
		return fmt.Errorf("making Read request on Automation Account %q (Resource Group %q): Unable to retrieve 'sku' value", id.Name, id.ResourceGroup)
	}

	d.Set("dsc_server_endpoint", keysResp.Endpoint)
	if keys := keysResp.Keys; keys != nil {
		d.Set("dsc_primary_access_key", keys.Primary)
		d.Set("dsc_secondary_access_key", keys.Secondary)
	}

	identity, err := flattenAutomationAccountIdentity(resp.Identity)
	if err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if err := d.Set("identity", identity); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if t := resp.Tags; t != nil {
		return tags.FlattenAndSet(d, t)
	}

	return nil
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

		return fmt.Errorf("issuing AzureRM delete request for Automation Account '%s': %+v", id.Name, err)
	}

	return nil
}
