package azurerm

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"

	"time"

	"strconv"

	"github.com/Azure/azure-sdk-for-go/services/resources/mgmt/2018-05-01/policy"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/structure"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyDefinitionCreateUpdate,
		Update: resourceArmPolicyDefinitionCreateUpdate,
		Read:   resourceArmPolicyDefinitionRead,
		Delete: resourceArmPolicyDefinitionDelete,
		Importer: &schema.ResourceImporter{
			State: resourceArmPolicyDefinitionImport,
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

			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.All),
					string(policy.Indexed),
					string(policy.NotSpecified),
				}, true),
			},

			"management_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"policy_rule": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.ValidateJsonString,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},
		},
	}
}

func resourceArmPolicyDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policy.DefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	mode := d.Get("mode").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	managementGroupID := d.Get("management_group_id").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := getPolicyDefinition(ctx, client, name, managementGroupID)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Policy Definition %q: %s", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_definition", *existing.ID)
		}
	}

	properties := policy.DefinitionProperties{
		PolicyType:  policy.Type(policyType),
		Mode:        policy.Mode(mode),
		DisplayName: utils.String(displayName),
		Description: utils.String(description),
	}

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := structure.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("unable to parse policy_rule: %s", err)
		}
		properties.PolicyRule = &policyRule
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("unable to parse metadata: %s", err)
		}
		properties.Metadata = &metaData
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := structure.ExpandJsonFromString(parametersString)
		if err != nil {
			return fmt.Errorf("unable to parse parameters: %s", err)
		}
		properties.Parameters = &parameters
	}

	definition := policy.Definition{
		Name:                 utils.String(name),
		DefinitionProperties: &properties,
	}

	var err error

	if managementGroupID == "" {
		_, err = client.CreateOrUpdate(ctx, name, definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, name, definition, managementGroupID)
	}

	if err != nil {
		return err
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policyDefinitionRefreshFunc(ctx, client, name, managementGroupID),
		Timeout:                   5 * time.Minute,
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Error waiting for Policy Definition %q to become available: %s", name, err)
	}

	resp, err := getPolicyDefinition(ctx, client, name, managementGroupID)
	if err != nil {
		return err
	}

	d.SetId(*resp.ID)

	return resourceArmPolicyDefinitionRead(d, meta)
}

func resourceArmPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policy.DefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name, err := parsePolicyDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	managementGroupID := parseManagementGroupIdFromPolicyId(d.Id())

	resp, err := getPolicyDefinition(ctx, client, name, managementGroupID)

	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("Error reading Policy Definition %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("management_group_id", managementGroupID)

	if props := resp.DefinitionProperties; props != nil {
		d.Set("policy_type", props.PolicyType)
		d.Set("mode", props.Mode)
		d.Set("display_name", props.DisplayName)
		d.Set("description", props.Description)

		if policyRuleStr := flattenJSON(props.PolicyRule); policyRuleStr != "" {
			d.Set("policy_rule", policyRuleStr)
		}

		if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
			d.Set("metadata", metadataStr)
		}

		if parametersStr := flattenJSON(props.Parameters); parametersStr != "" {
			d.Set("parameters", parametersStr)
		}
	}

	return nil
}

func resourceArmPolicyDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).policy.DefinitionsClient
	ctx := meta.(*ArmClient).StopContext

	name, err := parsePolicyDefinitionNameFromId(d.Id())
	if err != nil {
		return err
	}

	managementGroupID := parseManagementGroupIdFromPolicyId(d.Id())

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

		return fmt.Errorf("Error deleting Policy Definition %q: %+v", name, err)
	}

	return nil
}

func parsePolicyDefinitionNameFromId(id string) (string, error) {
	components := strings.Split(id, "/")

	if len(components) == 0 {
		return "", fmt.Errorf("Azure Policy Definition Id is empty or not formatted correctly: %s", id)
	}

	return components[len(components)-1], nil
}

func parseManagementGroupIdFromPolicyId(id string) string {
	r, _ := regexp.Compile("managementgroups/(.+)/providers/.*$")

	if r.MatchString(id) {
		parms := r.FindAllStringSubmatch(id, -1)[0]
		return parms[1]
	}

	return ""
}

func policyDefinitionRefreshFunc(ctx context.Context, client *policy.DefinitionsClient, name string, managementGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getPolicyDefinition(ctx, client, name, managementGroupID)

		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("Error issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q: %s", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
}

func resourceArmPolicyDefinitionImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if managementGroupID := parseManagementGroupIdFromPolicyId(d.Id()); managementGroupID != "" {
		d.Set("management_group_id", managementGroupID)

		if name, _ := parsePolicyDefinitionNameFromId(d.Id()); name != "" {
			d.Set("name", name)
		}

		return []*schema.ResourceData{d}, nil
	}

	//import a subscription based policy as before
	return schema.ImportStatePassthrough(d, meta)
}

func getPolicyDefinition(ctx context.Context, client *policy.DefinitionsClient, name string, managementGroupID string) (res policy.Definition, err error) {
	if managementGroupID == "" {
		res, err = client.Get(ctx, name)
	} else {
		res, err = client.GetAtManagementGroup(ctx, name, managementGroupID)
	}

	return res, err
}

func flattenJSON(stringMap interface{}) string {
	if stringMap != nil {
		value := stringMap.(map[string]interface{})
		jsonString, err := structure.FlattenJsonToString(value)
		if err == nil {
			return jsonString
		}
	}

	return ""
}
