package policy

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2019-09-01/policy"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyAssignment() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmPolicyAssignmentCreate,
		Update: resourceArmPolicyAssignmentUpdate,
		Read:   resourceArmPolicyAssignmentRead,
		Delete: resourceArmPolicyAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PolicyAssignmentID(id)
			return err
		}),

		DeprecationMessage: func() string {
			msg := `The 'azurerm_policy_assignment' resource is deprecated in favour of the:

- 'azurerm_management_group_policy_assignment'
- 'azurerm_resource_policy_assignment'
- 'azurerm_resource_group_policy_assignment'
- 'azurerm_subscription_policy_assignment'

resources and will be removed in version 3.0 of the Azure Provider.
`
			return strings.ReplaceAll(msg, "'", "`")
		}(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"scope": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_definition_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					validate.PolicyDefinitionID,
					validate.PolicySetDefinitionID,
				),
			},

			"description": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"display_name": {
				Type:     pluginsdk.TypeString,
				Optional: true,
			},

			"location": azure.SchemaLocationOptional(),

			//lintignore:XS003
			"identity": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"type": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(policy.None),
								string(policy.SystemAssigned),
							}, false),
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

			"parameters": {
				Type:             pluginsdk.TypeString,
				Optional:         true,
				ForceNew:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
			},

			"enforcement_mode": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},

			"not_scopes": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"metadata": metadataSchema(),
		},
	}
}

func resourceArmPolicyAssignmentCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPolicyAssignmentId(d.Get("scope").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id.Scope, id.Name)
	if err != nil {
		if !utils.ResponseWasNotFound(existing.Response) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !utils.ResponseWasNotFound(existing.Response) {
		return tf.ImportAsExistsError("azurerm_policy_assignment", id.ID())
	}

	assignment := policy.Assignment{
		AssignmentProperties: &policy.AssignmentProperties{
			PolicyDefinitionID: utils.String(d.Get("policy_definition_id").(string)),
			DisplayName:        utils.String(d.Get("display_name").(string)),
			Scope:              utils.String(id.Scope),
			EnforcementMode:    convertEnforcementMode(d.Get("enforcement_mode").(bool)),
		},
	}

	if v := d.Get("description").(string); v != "" {
		assignment.AssignmentProperties.Description = utils.String(v)
	}

	if v := d.Get("location").(string); v != "" {
		assignment.Location = utils.String(azure.NormalizeLocation(v))
	}

	if v, ok := d.GetOk("identity"); ok {
		if assignment.Location == nil {
			return fmt.Errorf("`location` must be set when `identity` is assigned")
		}
		assignment.Identity = expandAzureRmPolicyIdentity(v.([]interface{}))
	}

	if v := d.Get("parameters").(string); v != "" {
		expandedParams, err := expandParameterValuesValueFromString(v)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
		}

		assignment.AssignmentProperties.Parameters = expandedParams
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to parse metadata: %s", err)
		}
		assignment.AssignmentProperties.Metadata = &metaData
	}

	if v, ok := d.GetOk("not_scopes"); ok {
		assignment.AssignmentProperties.NotScopes = expandAzureRmPolicyNotScopes(v.([]interface{}))
	}

	if _, err := client.Create(ctx, id.Scope, id.Name, assignment); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Policy Assignments are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for %s to become available..", id)
	if err := waitForPolicyAssignmentToStabilize(ctx, client, id, true); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	d.SetId(id.ID())
	return resourceArmPolicyAssignmentRead(d, meta)
}

func resourceArmPolicyAssignmentUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyAssignmentID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, id.Scope, id.Name)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}
	if existing.AssignmentProperties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", *id)
	}

	update := policy.Assignment{
		Location:             existing.Location,
		AssignmentProperties: existing.AssignmentProperties,
	}
	if existing.Identity != nil {
		update.Identity = &policy.Identity{
			Type: existing.Identity.Type,
		}
	}

	if d.HasChange("description") {
		update.AssignmentProperties.Description = utils.String(d.Get("description").(string))
	}
	if d.HasChange("display_name") {
		update.AssignmentProperties.DisplayName = utils.String(d.Get("display_name").(string))
	}
	if d.HasChange("enforcement_mode") {
		update.AssignmentProperties.EnforcementMode = convertEnforcementMode(d.Get("enforcement_mode").(bool))
	}
	if d.HasChange("location") {
		update.Location = utils.String(d.Get("location").(string))
	}
	if d.HasChange("policy_definition_id") {
		update.AssignmentProperties.PolicyDefinitionID = utils.String(d.Get("policy_definition_id").(string))
	}

	if d.HasChange("identity") {
		if update.Location == nil {
			return fmt.Errorf("`location` must be set when `identity` is assigned")
		}
		identityRaw := d.Get("identity").([]interface{})
		update.Identity = expandAzureRmPolicyIdentity(identityRaw)
	}

	if d.HasChange("metadata") {
		v := d.Get("metadata").(string)
		update.AssignmentProperties.Metadata = map[string]interface{}{}
		if v != "" {
			metaData, err := pluginsdk.ExpandJsonFromString(v)
			if err != nil {
				return fmt.Errorf("parsing metadata: %+v", err)
			}
			update.AssignmentProperties.Metadata = &metaData
		}
	}

	if d.HasChange("not_scopes") {
		update.AssignmentProperties.NotScopes = expandAzureRmPolicyNotScopes(d.Get("not_scopes").([]interface{}))
	}

	if d.HasChange("parameters") {
		update.AssignmentProperties.Parameters = map[string]*policy.ParameterValuesValue{}

		if v := d.Get("parameters").(string); v != "" {
			expandedParams, err := expandParameterValuesValueFromString(v)
			if err != nil {
				return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
			}
			update.AssignmentProperties.Parameters = expandedParams
		}
	}

	// NOTE: there isn't an Update endpoint
	if _, err := client.Create(ctx, id.Scope, id.Name, update); err != nil {
		return fmt.Errorf("updating %s: %+v", *id, err)
	}

	// Policy Assignments are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for %s to become available..", id)
	if err := waitForPolicyAssignmentToStabilize(ctx, client, *id, true); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	return resourceArmPolicyAssignmentRead(d, meta)
}

func resourceArmPolicyAssignmentRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyAssignmentID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.Scope, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] %s was not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("scope", id.Scope)
	d.Set("location", location.NormalizeNilable(resp.Location))

	if err := d.Set("identity", flattenAzureRmPolicyIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
	}

	if props := resp.AssignmentProperties; props != nil {
		d.Set("policy_definition_id", props.PolicyDefinitionID)
		d.Set("description", props.Description)
		d.Set("display_name", props.DisplayName)
		d.Set("enforcement_mode", props.EnforcementMode == policy.Default)
		d.Set("metadata", flattenJSON(props.Metadata))

		json, err := flattenParameterValuesValueToString(props.Parameters)
		if err != nil {
			return fmt.Errorf("serializing JSON from `parameters`: %+v", err)
		}

		d.Set("parameters", json)
		d.Set("not_scopes", props.NotScopes)
	}

	return nil
}

func resourceArmPolicyAssignmentDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyAssignmentID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.Scope, id.Name); err != nil {
		return fmt.Errorf("deleting Policy Assignment %q: %+v", id, err)
	}

	// Policy Assignments are eventually consistent; wait for it to be gone
	log.Printf("[DEBUG] Waiting for %s to disappear..", id)
	if err := waitForPolicyAssignmentToStabilize(ctx, client, *id, false); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %s", id, err)
	}

	return nil
}

func expandAzureRmPolicyIdentity(input []interface{}) *policy.Identity {
	if len(input) == 0 || input[0] == nil {
		return nil
	}
	identity := input[0].(map[string]interface{})

	return &policy.Identity{
		Type: policy.ResourceIdentityType(identity["type"].(string)),
	}
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

func expandAzureRmPolicyNotScopes(input []interface{}) *[]string {
	notScopesRes := make([]string, 0)

	for _, notScope := range input {
		notScopesRes = append(notScopesRes, notScope.(string))
	}

	return &notScopesRes
}
