package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
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
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyAssignment() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyAssignmentCreateUpdate,
		Update: resourceArmPolicyAssignmentCreateUpdate,
		Read:   resourceArmPolicyAssignmentRead,
		Delete: resourceArmPolicyAssignmentDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
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

			//lintignore:XS003
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

			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: policyAssignmentsMetadataDiffSuppressFunc,
			},
		},
	}
}

func policyAssignmentsMetadataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	var oldPolicyAssignmentsMetadata map[string]interface{}
	errOld := json.Unmarshal([]byte(old), &oldPolicyAssignmentsMetadata)
	if errOld != nil {
		return false
	}

	var newPolicyAssignmentsMetadata map[string]interface{}
	if new != "" {
		errNew := json.Unmarshal([]byte(new), &newPolicyAssignmentsMetadata)
		if errNew != nil {
			return false
		}
	}

	// Ignore the following keys if they're found in the metadata JSON
	ignoreKeys := [5]string{"assignedBy", "createdBy", "createdOn", "updatedBy", "updatedOn"}
	for _, key := range ignoreKeys {
		delete(oldPolicyAssignmentsMetadata, key)
		delete(newPolicyAssignmentsMetadata, key)
	}

	return reflect.DeepEqual(oldPolicyAssignmentsMetadata, newPolicyAssignmentsMetadata)
}

func resourceArmPolicyAssignmentCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.AssignmentsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	scope := d.Get("scope").(string)

	if d.IsNewResource() {
		existing, err := client.Get(ctx, scope, name)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Policy Assignment %q: %s", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_assignment", *existing.ID)
		}
	}

	assignment := policy.Assignment{
		AssignmentProperties: &policy.AssignmentProperties{
			PolicyDefinitionID: utils.String(d.Get("policy_definition_id").(string)),
			DisplayName:        utils.String(d.Get("display_name").(string)),
			Scope:              utils.String(scope),
			EnforcementMode:    convertEnforcementMode(d.Get("enforcement_mode").(bool)),
		},
	}

	if v := d.Get("description").(string); v != "" {
		assignment.AssignmentProperties.Description = utils.String(v)
	}

	if v, ok := d.GetOk("identity"); ok {
		if location := d.Get("location").(string); location == "" {
			return fmt.Errorf("`location` must be set when `identity` is assigned")
		}
		assignment.Identity = expandAzureRmPolicyIdentity(v.([]interface{}))
	}

	if v := d.Get("location").(string); v != "" {
		assignment.Location = utils.String(azure.NormalizeLocation(v))
	}

	if v := d.Get("parameters").(string); v != "" {
		expandedParams, err := expandParameterValuesValueFromString(v)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters` %q: %+v", v, err)
		}

		assignment.AssignmentProperties.Parameters = expandedParams
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to parse metadata: %s", err)
		}
		assignment.AssignmentProperties.Metadata = &metaData
	}

	if v, ok := d.GetOk("not_scopes"); ok {
		assignment.AssignmentProperties.NotScopes = expandAzureRmPolicyNotScopes(v.([]interface{}))
	}

	if _, err := client.Create(ctx, scope, name, assignment); err != nil {
		return fmt.Errorf("creating/updating Policy Assignment %q (Scope %q): %+v", name, scope, err)
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
		return fmt.Errorf("waiting for Policy Assignment %q to become available: %s", name, err)
	}

	resp, err := client.Get(ctx, scope, name)
	if err != nil {
		return fmt.Errorf("retrieving Policy Assignment %q (Scope %q): %+v", name, scope, err)
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Assignment %q (Scope %q)", name, scope)
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

		return fmt.Errorf("reading Policy Assignment %q: %+v", id, err)
	}

	d.Set("name", resp.Name)

	if err := d.Set("identity", flattenAzureRmPolicyIdentity(resp.Identity)); err != nil {
		return fmt.Errorf("setting `identity`: %+v", err)
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

		if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}

		if params := props.Parameters; params != nil {
			json, err := flattenParameterValuesValueToString(params)
			if err != nil {
				return fmt.Errorf("serializing JSON from `parameters`: %+v", err)
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

		return fmt.Errorf("deleting Policy Assignment %q: %+v", id, err)
	}

	return nil
}

func policyAssignmentRefreshFunc(ctx context.Context, client *policy.AssignmentsClient, scope string, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, scope, name)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q (Scope: %q): %s", name, scope, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
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

func convertEnforcementMode(mode bool) policy.EnforcementMode {
	if mode {
		return policy.Default
	} else {
		return policy.DoNotEnforce
	}
}
