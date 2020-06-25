package policy

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyAssignmentCreateUpdate,
		Update: resourceArmPolicyAssignmentCreateUpdate,
		Read:   resourceArmPolicyAssignmentRead,
		Delete: resourceArmPolicyAssignmentDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PolicyAssignmentID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scope": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_definition_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validate.PolicyDefinitionID,
					validate.PolicySetDefinitionID,
				),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"location": azure.SchemaLocationOptional(),

			"identity": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(policy.None),
								string(policy.SystemAssigned),
							}, false),
						},
						"principal_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tenant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"enforcement_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},

			"not_scopes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceArmPolicyAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)
	enforcementMode := convertEnforcementMode(d.Get("enforcement_mode").(bool))
	policyDefinitionId := d.Get("policy_definition_id").(string)
	displayName := d.Get("display_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Policy Assignment %q: %s", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_assignment", *existing.ID)
		}
	}

	assignment := policy.Assignment{
		AssignmentProperties: &policy.AssignmentProperties{
			PolicyDefinitionID: utils.String(policyDefinitionId),
			DisplayName:        utils.String(displayName),
			Scope:              utils.String(scope),
			EnforcementMode:    enforcementMode,
		},
	}

	if v := d.Get("description").(string); v != "" {
		assignment.AssignmentProperties.Description = utils.String(v)
	}

	if _, ok := d.GetOk("identity"); ok {
		if v := d.Get("location").(string); v == "" {
			return fmt.Errorf("`location` must be set when `identity` is assigned")
		}
		policyIdentity := expandAzureRmPolicyIdentity(d)
		assignment.Identity = policyIdentity
	}

	if v := d.Get("location").(string); v != "" {
		assignment.Location = utils.String(azure.NormalizeLocation(v))
	}

	if v := d.Get("parameters").(string); v != "" {
		expandedParams, err := expandParameterValuesValueFromString(v)
		if err != nil {
			return fmt.Errorf("Error expanding JSON from Parameters %q: %+v", v, err)
		}

		assignment.AssignmentProperties.Parameters = expandedParams
	}

	if _, ok := d.GetOk("not_scopes"); ok {
		notScopes := expandAzureRmPolicyNotScopes(d)
		assignment.AssignmentProperties.NotScopes = notScopes
	}

	if _, err := client.Create(ctx, scope, name, assignment); err != nil {
		return err
	}

	// Policy Assignments are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Assignment %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policyAssignmentRefreshFunc(ctx, client, scope, name),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Policy Assignment %q to become available: %s", name, err)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmPolicyAssignmentRead(d, meta)
}

func resourceArmPolicyAssignmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()

	resp, err := client.GetByID(ctx, id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Assignment %q - removing from state", id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Policy Assignment %q: %+v", id, err)
	}

	d.Set("name", resp.Name)

	if err := d.Set("identity", flattenAzureRmPolicyIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("Error setting `identity`: %+v", err)
	}

	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}

	if props := resp.AssignmentProperties; props != nil {
		d.Set("scope", props.Scope)
		d.Set("policy_definition_id", props.PolicyDefinitionID)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("enforcement_mode", props.EnforcementMode == policy.Default)

		if params := props.Parameters; params != nil {
			json, err := flattenParameterValuesValueToString(params)
			if err != nil {
				return fmt.Errorf("Error serializing JSON from Parameters: %+v", err)
			}

			d.Set("parameters", json)
		}

		d.Set("not_scopes", props.NotScopes)
	}

	return nil
}

func resourceArmPolicyAssignmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := d.Id()

	resp, err := client.DeleteByID(ctx, id)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			return nil
		}

		return fmt.Errorf("Error deleting Policy Assignment %q: %+v", id, err)
	}

	return nil
}

func policyAssignmentRefreshFunc(ctx context.Context, client *policy.AssignmentsClient, scope string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q (Scope: %q): %s", name, scope, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func expandAzureRmPolicyIdentity(d *schema.ResourceData) *policy.Identity {
	v := d.Get("identity")
	identities := v.([]interface{})
	identity := identities[0].(map[string]interface{})

	identityType := policy.ResourceIdentityType(identity["type"].(string))

	policyIdentity := policy.Identity{
		Type: identityType,
	}

	return &policyIdentity
}

func flattenAzureRmPolicyIdentity(identity *policy.Identity) []interface{} {
	if identity == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})
	result["type"] = string(identity.Type)
	if identity.PrincipalID != nil {
		result["principal_id"] = *identity.PrincipalID
	}

	if identity.TenantID != nil {
		result["tenant_id"] = *identity.TenantID
	}

	return []interface{}{result}
}

func expandAzureRmPolicyNotScopes(d *schema.ResourceData) *[]string {
	notScopes := d.Get("not_scopes").([]interface{})
	notScopesRes := make([]string, 0)

	for _, notScope := range notScopes {
		notScopesRes = append(notScopesRes, notScope.(string))
	}

	return &notScopesRes
}

func convertEnforcementMode(mode bool) policy.EnforcementMode {
	if mode {
		return policy.Default
	} else {
		return policy.DoNotEnforce
	}
}
