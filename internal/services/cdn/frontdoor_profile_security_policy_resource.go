package cdn

import (
	"fmt"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/profiles"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/cdn/sdk/2021-06-01/securitypolicies"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceFrontdoorProfileSecurityPolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceFrontdoorProfileSecurityPolicyCreate,
		Read:   resourceFrontdoorProfileSecurityPolicyRead,
		Update: resourceFrontdoorProfileSecurityPolicyUpdate,
		Delete: resourceFrontdoorProfileSecurityPolicyDelete,

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := securitypolicies.ParseSecurityPoliciesID(id)
			return err
		}),

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"cdn_profile_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: profiles.ValidateProfileID,
			},

			"deployment_status": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"parameters": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,

				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{

						"type": {
							Type:     pluginsdk.TypeString,
							Required: true,
						},
					},
				},
			},

			"profile_name": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"provisioning_state": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceFrontdoorProfileSecurityPolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileSecurityPoliciesClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	profileId, err := profiles.ParseProfileID(d.Get("cdn_profile_id").(string))
	if err != nil {
		return err
	}

	id := securitypolicies.NewSecurityPoliciesID(profileId.SubscriptionId, profileId.ResourceGroupName, profileId.ProfileName, d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id)
		if err != nil {
			if !response.WasNotFound(existing.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(existing.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_cdn_security_policy", id.ID())
		}
	}

	props := securitypolicies.SecurityPolicy{
		Properties: &securitypolicies.SecurityPolicyProperties{
			Parameters: expandSecurityPoliciesSecurityPolicyParameters(d.Get("parameters").([]interface{})),
		},
	}
	if err := client.CreateThenPoll(ctx, id, props); err != nil {

		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceFrontdoorProfileSecurityPolicyRead(d, meta)
}

func resourceFrontdoorProfileSecurityPolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileSecurityPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.SecurityPolicyName)

	d.Set("cdn_profile_id", profiles.NewProfileID(id.SubscriptionId, id.ResourceGroupName, id.ProfileName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("deployment_status", props.DeploymentStatus)

			if err := d.Set("parameters", flattenSecurityPoliciesSecurityPolicyParameters(props.Parameters)); err != nil {
				return fmt.Errorf("setting `parameters`: %+v", err)
			}
			d.Set("profile_name", props.ProfileName)
			d.Set("provisioning_state", props.ProvisioningState)
		}
	}
	return nil
}

func resourceFrontdoorProfileSecurityPolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileSecurityPoliciesClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	props := securitypolicies.SecurityPolicyProperties{
		Parameters: expandSecurityPoliciesSecurityPolicyParameters(d.Get("parameters").([]interface{})),
	}
	if err := client.PatchThenPoll(ctx, *id, props); err != nil {

		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceFrontdoorProfileSecurityPolicyRead(d, meta)
}

func resourceFrontdoorProfileSecurityPolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Cdn.FrontdoorProfileSecurityPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := securitypolicies.ParseSecurityPoliciesID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {

		return fmt.Errorf("deleting %s: %+v", id, err)
	}
	return nil
}

func expandSecurityPoliciesSecurityPolicyParameters(input []interface{}) *securitypolicies.SecurityPolicyPropertiesParameters {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})

	typeValue := securitypolicies.SecurityPolicyType(v["type"].(string))
	return &securitypolicies.SecurityPolicyPropertiesParameters{
		Type: typeValue,
	}
}

func flattenSecurityPoliciesSecurityPolicyParameters(input *securitypolicies.SecurityPolicyPropertiesParameters) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	result := make(map[string]interface{})
	result["type"] = input.Type
	return append(results, result)
}
