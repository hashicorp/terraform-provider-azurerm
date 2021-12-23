package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceAfdRuleSet() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceAfdRuleSetCreate,
		Read:   resourceAfdRuleSetRead,
		Delete: resourceAfdRuleSetDelete,

		SchemaVersion: 1,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.AfdRuleSetsID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.ProfileID,
			},
		},
	}
}

func resourceAfdRuleSetCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	// rule set name
	ruleSetName := d.Get("name").(string)

	// parse profile_id
	profileId := d.Get("profile_id").(string)
	profile, err := parse.ProfileID(profileId)
	if err != nil {
		return err
	}

	id := parse.NewAfdRuleSetsID(profile.SubscriptionId, profile.ResourceGroup, profile.Name, ruleSetName)

	future, err := client.Create(ctx, profile.ResourceGroup, profile.Name, ruleSetName)
	if err != nil {
		return fmt.Errorf("creating rule set %s: %+v", ruleSetName, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the creation of rule set %s: %+v", ruleSetName, err)
	}

	d.SetId(id.ID())

	return resourceAfdRuleSetRead(d, meta)
}

func resourceAfdRuleSetRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdRuleSetsID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("making read request on Azure CDN Front Door Rule Set %q (Resource Group %q): %+v", id.RuleSetName, id.ResourceGroup, err)
	}

	d.Set("name", resp.Name)

	// profile id
	profileId := parse.NewProfileID(id.SubscriptionId, id.ResourceGroup, id.ProfileName).ID()
	d.Set("profile_id", string(profileId))

	return nil
}

func resourceAfdRuleSetDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.AFDRuleSetsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.AfdRuleSetsID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.ProfileName, id.RuleSetName)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
	}

	return err
}
