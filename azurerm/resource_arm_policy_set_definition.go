package azurerm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicySetDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicySetDefinitionCreateUpdate,
		Update: resourceArmPolicySetDefinitionCreateUpdate,
		Read:   resourceArmPolicySetDefinitionRead,
		Delete: resourceArmPolicySetDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.TypeBuiltIn),
					string(policy.TypeCustom),
				}, false),
			},

			"management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"display_name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"policy_definitions": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: policyDefinitionsDiffSuppressFunc,
			},
		},
	}
}

func policyDefinitionsDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	var oldPolicyDefinitions []policy.DefinitionReference
	errOld := json.Unmarshal([]byte(old), &oldPolicyDefinitions)
	if errOld != nil {
		return false
	}

	var newPolicyDefinitions []policy.DefinitionReference
	errNew := json.Unmarshal([]byte(new), &newPolicyDefinitions)
	if errNew != nil {
		return false
	}

	return reflect.DeepEqual(oldPolicyDefinitions, newPolicyDefinitions)
}

func resourceArmPolicySetDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	managementGroupID := d.Get("management_group_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := getPolicySetDefinition(ctx, client, name, managementGroupID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Policy Set Definition %q: %s", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_set_definition", *existing.ID)
		}
	}

	properties := policy.SetDefinitionProperties{
		PolicyType:  policy.Type(policyType),
		DisplayName: utils.String(displayName),
		Description: utils.String(description),
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to expand metadata json: %s", err)
		}
		properties.Metadata = &metaData
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := structure.ExpandJsonFromString(parametersString)
		if err != nil {
			return fmt.Errorf("unable to expand parameters json: %s", err)
		}
		properties.Parameters = &parameters
	}

	if policyDefinitionsString := d.Get("policy_definitions").(string); policyDefinitionsString != "" {
		var policyDefinitions []policy.DefinitionReference
		err := json.Unmarshal([]byte(policyDefinitionsString), &policyDefinitions)
		if err != nil {
			return fmt.Errorf("unable to expand parameters json: %s", err)
		}
		properties.PolicyDefinitions = &policyDefinitions
	}

	definition := policy.SetDefinition{
		Name:                    utils.String(name),
		SetDefinitionProperties: &properties,
	}

	var err error
	if managementGroupID == "" {
		_, err = client.CreateOrUpdate(ctx, name, definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, name, definition, managementGroupID)
	}

	if err != nil {
		return fmt.Errorf("Error creating/updating Policy Set Definition %q: %s", name, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Set Definition %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, name, managementGroupID),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if features.SupportsCustomTimeouts() {
		if d.IsNewResource() {
			stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
		} else {
			stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
		}
	} else {
		stateConf.Timeout = 5 * time.Minute
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Policy Set Definition %q to become available: %s", name, err)
	}

	var resp policy.SetDefinition
	resp, err = getPolicySetDefinition(ctx, client, name, managementGroupID)
	if err != nil {
		return fmt.Errorf("Error retrieving Policy Set Definition %q: %s", name, err)
	}

	d.SetId(*resp.ID)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := parsePolicySetDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	managementGroupID := parseManagementGroupIdFromPolicySetId(d.Id())

	resp, err := getPolicySetDefinition(ctx, client, name, managementGroupID)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Set Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Policy Set Definition %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("management_group_id", managementGroupID)

	if props := resp.SetDefinitionProperties; props != nil {
		d.Set("policy_type", string(props.PolicyType))
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)

		if metadata := props.Metadata; metadata != nil {
			metadataVal := metadata.(map[string]interface{})
			metadataStr, err := structure.FlattenJsonToString(metadataVal)
			if err != nil {
				return fmt.Errorf("unable to flatten JSON for `metadata`: %s", err)
			}

			d.Set("metadata", metadataStr)
		}

		if parameters := props.Parameters; parameters != nil {
			paramsVal := parameters.(map[string]interface{})
			parametersStr, err := structure.FlattenJsonToString(paramsVal)
			if err != nil {
				return fmt.Errorf("unable to flatten JSON for `parameters`: %s", err)
			}

			d.Set("parameters", parametersStr)
		}

		if policyDefinitions := props.PolicyDefinitions; policyDefinitions != nil {
			policyDefinitionsRes, err := json.Marshal(policyDefinitions)
			if err != nil {
				return fmt.Errorf("unable to flatten JSON for `policy_defintions`: %s", err)
			}

			policyDefinitionsStr := string(policyDefinitionsRes)
			d.Set("policy_definitions", policyDefinitionsStr)
		}
	}

	return nil
}

func resourceArmPolicySetDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.SetDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name, err := parsePolicySetDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	managementGroupID := parseManagementGroupIdFromPolicySetId(d.Id())

	var resp autorest.Response
	if managementGroupID == "" {
		resp, err = client.Delete(ctx, name)
	} else {
		resp, err = client.DeleteAtManagementGroup(ctx, name, managementGroupID)
	}

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Policy Set Definition %q: %+v", name, err)
	}

	return nil
}

func parsePolicySetDefinitionNameFromId(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Azure Policy Set Definition Id is empty or not formatted correctly: %s", id)
	}

	return components[len(components)-1], nil
}

func parseManagementGroupIdFromPolicySetId(id string) string {
	r, _ := regexp.Compile("managementgroups/(.+)/providers/.*$")

	if r.MatchString(id) {
		matches := r.FindAllStringSubmatch(id, -1)[0]
		return matches[1]
	}

	return ""
}

func policySetDefinitionRefreshFunc(ctx context.Context, client *policy.SetDefinitionsClient, name string, managementGroupId string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getPolicySetDefinition(ctx, client, name, managementGroupId)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error issuing read request in policySetDefinitionRefreshFunc for Policy Set Definition %q: %s", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func getPolicySetDefinition(ctx context.Context, client *policy.SetDefinitionsClient, name string, managementGroupID string) (res policy.SetDefinition, err error) {
	if managementGroupID == "" {
		res, err = client.Get(ctx, name)
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupID)
	}

	return res, err
}
