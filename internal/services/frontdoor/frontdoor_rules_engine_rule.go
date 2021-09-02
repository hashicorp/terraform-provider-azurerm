package frontdoor

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/frontdoor/mgmt/2020-05-01/frontdoor"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/location"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/parse"
	azValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/frontdoor/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceFrontDoorRulesEngine() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontDoorRulesEngineCreateUpdate,
		Read:   resourceFrontDoorRulesEngineRead,
		Update: resourceFrontDoorRulesEngineCreateUpdate,
		Delete: resourceFrontDoorRulesEngineDelete,

		SchemaVersion: 1,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.WebApplicationFirewallPolicyIDInsensitively(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(6 * time.Hour),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(6 * time.Hour),
			Delete: pluginsdk.DefaultTimeout(6 * time.Hour),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},
			"frontdoor_name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azValidate.FrontDoorName,
			},
			"location": location.SchemaComputed(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceFrontDoorRulesEngineCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	frontdoorRulesEngineRuleHeaderAction := frontdoor.HeaderAction{
		HeaderActionType: "Append", // HeaderActionType - 'Append', 'Delete', 'Overwrite'
		HeaderName:       utils.String("X-TEST-HEADER"),
		Value:            utils.String("this is a test"),
	}

	frontdoorRulesEngineRuleHeaderActions := make([]frontdoor.HeaderAction, 0)
	frontdoorRulesEngineRuleHeaderActions = append(frontdoorRulesEngineRuleHeaderActions, frontdoorRulesEngineRuleHeaderAction)

	frontdoorRulesEngineRuleAction := frontdoor.RulesEngineAction{
		RequestHeaderActions: &frontdoorRulesEngineRuleHeaderActions,
		//ResponseHeaderActions
		//RouteConfigurationOverride
	}

	frontdoorRulesEngineRule := frontdoor.RulesEngineRule{
		Name:     utils.String("testrule"),
		Priority: utils.Int32(1),
		Action:   &frontdoorRulesEngineRuleAction,
		//MatchConditions:
		//MatchProcessingBehavior:
	}

	frontdoorRulesEngineRules := make([]frontdoor.RulesEngineRule, 0)
	frontdoorRulesEngineRules = append(frontdoorRulesEngineRules, frontdoorRulesEngineRule)

	frontdoorRulesEngineProperties := frontdoor.RulesEngineProperties{
		Rules: &frontdoorRulesEngineRules,
	}

	frontdoorRulesEngine := frontdoor.RulesEngine{
		Name:                  utils.String(rulesEngineName),
		RulesEngineProperties: &frontdoorRulesEngineProperties,
		//Type
		//ID
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, frontDoorName, rulesEngineName, frontdoorRulesEngine)
	if err != nil {
		return fmt.Errorf("creating Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}
	return resourceFrontDoorRulesEngineRead(d, meta)
}

func resourceFrontDoorRulesEngineRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	resp, err := client.Get(ctx, resourceGroup, frontDoorName, rulesEngineName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Front Door Rules Engine %q does not exist - removing from state", rulesEngineName)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}

	return nil
}

func resourceFrontDoorRulesEngineDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Frontdoor.FrontDoorsRulesEnginesClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	frontDoorName := d.Get("frontdoor_name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	rulesEngineName := d.Get("name").(string)

	future, err := client.Delete(ctx, resourceGroup, frontDoorName, rulesEngineName)
	if err != nil {

		return fmt.Errorf("deleting Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("waiting for deleting Front Door Rules Engine %q (Resource Group %q): %+v", rulesEngineName, resourceGroup, err)
		}
	}

	return nil
}
