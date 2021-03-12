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
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/structure"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/policy/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPolicyDefinition() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPolicyDefinitionCreateUpdate,
		Update: resourceArmPolicyDefinitionCreateUpdate,
		Read:   resourceArmPolicyDefinitionRead,
		Delete: resourceArmPolicyDefinitionDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.PolicyDefinitionID(id)
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

			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(policy.BuiltIn),
					string(policy.Custom),
					string(policy.NotSpecified),
					string(policy.Static),
				}, true),
			},

			"mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice(
					[]string{
						"All",
						"Indexed",
						"Microsoft.ContainerService.Data",
						"Microsoft.CustomerLockbox.Data",
						"Microsoft.DataCatalog.Data",
						"Microsoft.KeyVault.Data",
						"Microsoft.Kubernetes.Data",
						"Microsoft.MachineLearningServices.Data",
						"Microsoft.Network.Data",
						"Microsoft.Synapse.Data",
					}, false,
				),
			},

			"management_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true,
				ConflictsWith: []string{"management_group_name"},
				Deprecated:    "Deprecated in favour of `management_group_name`", // TODO -- remove this in next major version
			},

			"management_group_name": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Computed:      true, // TODO -- remove this when deprecation resolves
				ConflictsWith: []string{"management_group_id"},
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
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"parameters": {
				Type:             schema.TypeString,
				Optional:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: structure.SuppressJsonDiff,
			},

			"metadata": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: policyDefinitionsMetadataDiffSuppressFunc,
			},
		},
	}
}

func policyDefinitionsMetadataDiffSuppressFunc(_, old, new string, _ *schema.ResourceData) bool {
	var oldPolicyDefinitionsMetadata map[string]interface{}
	errOld := json.Unmarshal([]byte(old), &oldPolicyDefinitionsMetadata)
	if errOld != nil {
		return false
	}

	var newPolicyDefinitionsMetadata map[string]interface{}
	errNew := json.Unmarshal([]byte(new), &newPolicyDefinitionsMetadata)
	if errNew != nil {
		return false
	}

	// Ignore the following keys if they're found in the metadata JSON
	ignoreKeys := [4]string{"createdBy", "createdOn", "updatedBy", "updatedOn"}
	for _, key := range ignoreKeys {
		delete(oldPolicyDefinitionsMetadata, key)
		delete(newPolicyDefinitionsMetadata, key)
	}

	return reflect.DeepEqual(oldPolicyDefinitionsMetadata, newPolicyDefinitionsMetadata)
}

func resourceArmPolicyDefinitionCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	policyType := d.Get("policy_type").(string)
	mode := d.Get("mode").(string)
	displayName := d.Get("display_name").(string)
	description := d.Get("description").(string)
	managementGroupName := ""
	if v, ok := d.GetOk("management_group_name"); ok {
		managementGroupName = v.(string)
	}
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupName = v.(string)
	}

	if d.IsNewResource() {
		existing, err := getPolicyDefinitionByName(ctx, client, name, managementGroupName)
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing Policy Definition %q: %+v", name, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_policy_definition", *existing.ID)
		}
	}

	properties := policy.DefinitionProperties{
		PolicyType:  policy.Type(policyType),
		Mode:        utils.String(mode),
		DisplayName: utils.String(displayName),
		Description: utils.String(description),
	}

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := structure.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_rule`: %+v", err)
		}
		properties.PolicyRule = &policyRule
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := structure.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		properties.Metadata = &metaData
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		parameters, err := expandParameterDefinitionsValueFromString(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		properties.Parameters = parameters
	}

	definition := policy.Definition{
		Name:                 utils.String(name),
		DefinitionProperties: &properties,
	}

	var err error

	if managementGroupName == "" {
		_, err = client.CreateOrUpdate(ctx, name, definition)
	} else {
		_, err = client.CreateOrUpdateAtManagementGroup(ctx, name, definition, managementGroupName)
	}

	if err != nil {
		return fmt.Errorf("creating/updating Policy Definition %q: %+v", name, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	stateConf := &resource.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policyDefinitionRefreshFunc(ctx, client, name, managementGroupName),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(schema.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(schema.TimeoutUpdate)
	}

	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for Policy Definition %q to become available: %+v", name, err)
	}

	resp, err := getPolicyDefinitionByName(ctx, client, name, managementGroupName)
	if err != nil {
		return err
	}

	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("empty or nil ID returned for Policy Assignment %q", name)
	}
	d.SetId(*resp.ID)

	return resourceArmPolicyDefinitionRead(d, meta)
}

func resourceArmPolicyDefinitionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupName = scopeId.ManagementGroupName
	}

	resp, err := getPolicyDefinitionByName(ctx, client, id.Name, managementGroupName)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Error reading Policy Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading Policy Definition %+v", err)
	}

	d.Set("name", resp.Name)
	d.Set("management_group_id", managementGroupName)
	d.Set("management_group_name", managementGroupName)

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

		if parametersStr, err := flattenParameterDefinitionsValueToString(props.Parameters); err == nil {
			d.Set("parameters", parametersStr)
		} else {
			return fmt.Errorf("flattening policy definition parameters %+v", err)
		}
	}

	return nil
}

func resourceArmPolicyDefinitionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.DefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	switch scopeId := id.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupName = scopeId.ManagementGroupName
	}

	var resp autorest.Response
	if managementGroupName == "" {
		resp, err = client.Delete(ctx, id.Name)
	} else {
		resp, err = client.DeleteAtManagementGroup(ctx, id.Name, managementGroupName)
	}

	if err != nil {
		if utils.ResponseWasNotFound(resp) {
			return nil
		}

		return fmt.Errorf("deleting Policy Definition %q: %+v", id.Name, err)
	}

	return nil
}

func policyDefinitionRefreshFunc(ctx context.Context, client *policy.DefinitionsClient, name, managementGroupID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := getPolicyDefinitionByName(ctx, client, name, managementGroupID)
		if err != nil {
			return nil, strconv.Itoa(res.StatusCode), fmt.Errorf("issuing read request in policyAssignmentRefreshFunc for Policy Assignment %q: %+v", name, err)
		}

		return res, strconv.Itoa(res.StatusCode), nil
	}
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
