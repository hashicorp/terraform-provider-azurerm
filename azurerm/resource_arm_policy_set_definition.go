package azurerm

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"strings"

	"time"

	"strconv"

	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
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

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.TypeBuiltIn),
					string(policy.TypeCustom),
					string(policy.TypeNotSpecified),
				}, true)},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
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
	errNew := json.Unmarshal([]byte(old), &newPolicyDefinitions)
	if errNew != nil {
		return false
	}

	return reflect.DeepEqual(oldPolicyDefinitions, newPolicyDefinitions)
}

func resourceArmPolicySetDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policySetDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)

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

	if _, err := client.CreateOrUpdate(ctx, name, definition); err != nil {
		return err
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Set Definition %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, name),
		Timeout:                   5 * time.Minute,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Policy Set Definition %q to become available: %s", name, err)
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policySetDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name, err := parsePolicySetDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Policy Definition %+v", err)
	}

	d.Set("name", resp.Name)

	if props := resp.SetDefinitionProperties; props != nil {
		d.Set("policy_type", props.PolicyType)
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
			paramsVal := props.Parameters.(map[string]interface{})
			parametersStr, err := structure.FlattenJsonToString(paramsVal)
			if err != nil {
				return fmt.Errorf("unable to flatten JSON for `parameters`: %s", err)
			}

			d.Set("parameters", parametersStr)
		}

		if policyDefinitions := props.PolicyDefinitions; policyDefinitions != nil {
			policyDefinitionsRes, err := json.Marshal(props.PolicyDefinitions)
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
	client := meta.(*ArmClient).policySetDefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name, err := parsePolicySetDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Delete(ctx, name)
	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("Error deleting Policy Definition %q: %+v", name, err)
	}

	return nil
}

func parsePolicySetDefinitionNameFromId(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Azure Policy Set Definition Id is empty or not formatted correctly: %s", id)
	}

	if len(components) != 7 {
		return "", fmt.Errorf("Azure Policy Set Definition Id should have 6 segments, got %d: '%s'", len(components)-1, id)
	}

	return components[6], nil
}

func policySetDefinitionRefreshFunc(ctx context.Context, client policy.SetDefinitionsClient, name string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, name)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q: %s", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}
