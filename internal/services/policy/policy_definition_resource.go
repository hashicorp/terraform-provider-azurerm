// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceArmPolicyDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmPolicyDefinitionCreateUpdate,
		Update: resourceArmPolicyDefinitionCreateUpdate,
		Read:   resourceArmPolicyDefinitionRead,
		Delete: resourceArmPolicyDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PolicyDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourceArmPolicyDefinitionSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			// `parameters` cannot have values removed so we'll ForceNew if there are less parameters between Terraform runs
			if d.HasChange("parameters") {
				oldParametersRaw, newParametersRaw := d.GetChange("parameters")
				if oldParametersString := oldParametersRaw.(string); oldParametersString != "" {
					newParametersString := newParametersRaw.(string)
					if newParametersString == "" {
						return d.ForceNew("parameters")
					}

					oldParameters, err := expandParameterDefinitionsValue(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := expandParameterDefinitionsValue(newParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					if len(*newParameters) < len(*oldParameters) {
						return d.ForceNew("parameters")
					}
				}
			}

			return nil
		}),
	}
}

func resourceArmPolicyDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	var id any
	managementGroupName := ""
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupId, err := mgmtGrpParse.ManagementGroupID(v.(string))
		if err != nil {
			return err
		}
		managementGroupName = managementGroupId.Name
		id = policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, name)
	} else {
		id = policydefinitions.NewProviderPolicyDefinitionID(subscriptionId, name)
	}

	if d.IsNewResource() {
		resp, _, err := getPolicyDefinitionByID(ctx, client, id)
		if err != nil {
			if !response.WasNotFound(resp) {
				return fmt.Errorf("checking for presence of existing Policy Definition %q: %+v", name, err)
			}
		}

		if !response.WasNotFound(resp) {
			var idString string
			switch typedId := id.(type) {
			case policydefinitions.ProviderPolicyDefinitionId:
				idString = typedId.ID()
			case policydefinitions.Providers2PolicyDefinitionId:
				idString = typedId.ID()
			}
			return tf.ImportAsExistsError("azurerm_policy_definition", idString)
		}
	}

	parameters := policydefinitions.PolicyDefinition{
		Properties: &policydefinitions.PolicyDefinitionProperties{
			DisplayName: pointer.To(d.Get("display_name").(string)),
			Description: pointer.To(d.Get("description").(string)),
			PolicyType:  pointer.To(policydefinitions.PolicyType(d.Get("policy_type").(string))),
			Mode:        pointer.To(d.Get("mode").(string)),
		},
	}
	props := parameters.Properties

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := pluginsdk.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_rule`: %+v", err)
		}
		var iPolicyRule interface{} = policyRule
		props.PolicyRule = &iPolicyRule
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		var iMetadata interface{} = metaData
		props.Metadata = &iMetadata
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		params, err := expandParameterDefinitionsValueForPolicyDefinition(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		props.Parameters = params
	}

	switch typedId := id.(type) {
	case policydefinitions.ProviderPolicyDefinitionId:
		if _, err := client.CreateOrUpdate(ctx, typedId, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", typedId, err)
		}
	case policydefinitions.Providers2PolicyDefinitionId:
		if _, err := client.CreateOrUpdateAtManagementGroup(ctx, typedId, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", typedId, err)
		}
	default:
		return fmt.Errorf("unsupported ID type: %T", id)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policyDefinitionRefreshFunc(ctx, client, id),
		MinTimeout:                10 * time.Second,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		ContinuousTargetOccurence: 10,
	}

	if !d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for Policy Definition %q to become available: %+v", name, err)
	}

	var idString string
	switch typedId := id.(type) {
	case policydefinitions.ProviderPolicyDefinitionId:
		idString = typedId.ID()
	case policydefinitions.Providers2PolicyDefinitionId:
		idString = typedId.ID()
	}

	resourceId, err := parse.PolicyDefinitionID(idString)
	if err != nil {
		return err
	}

	d.SetId(resourceId.Id)

	return resourceArmPolicyDefinitionRead(d, meta)
}

func resourceArmPolicyDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	oldId, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	var id any
	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	switch scopeId := oldId.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
		id = policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, oldId.Name)
	default:
		id = policydefinitions.NewProviderPolicyDefinitionID(meta.(*clients.Client).Account.SubscriptionId, oldId.Name)
	}

	resp, model, err := getPolicyDefinitionByID(ctx, client, id)
	if err != nil {
		if response.WasNotFound(resp) {
			log.Printf("[INFO] policy definition was not found - removing from state")
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading policy definition: %+v", err)
	}

	d.Set("name", oldId.Name)

	d.Set("management_group_id", "")
	if managementGroupName != "" {
		d.Set("management_group_id", managementGroupId.ID())
	}

	if model != nil {
		if props := model.Properties; props != nil {
			d.Set("policy_type", string(pointer.From(props.PolicyType)))
			d.Set("mode", props.Mode)
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)

			if policyRuleStr := flattenJSON(props.PolicyRule); policyRuleStr != "" {
				d.Set("policy_rule", policyRuleStr)
				roleIDs, _ := getPolicyRoleDefinitionIDs(policyRuleStr)
				d.Set("role_definition_ids", roleIDs)
			}

			if metadataStr := flattenJSON(props.Metadata); metadataStr != "" {
				d.Set("metadata", metadataStr)
			}

			if parametersStr, err := flattenParameterDefinitionsValueToStringForPolicyDefinition(props.Parameters); err == nil {
				d.Set("parameters", parametersStr)
			} else {
				return fmt.Errorf("flattening policy definition parameters: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmPolicyDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	oldId, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	switch scopeId := oldId.PolicyScopeId.(type) { // nolint gocritic
	case parse.ScopeAtManagementGroup:
		id := policydefinitions.NewProviders2PolicyDefinitionID(scopeId.ManagementGroupName, oldId.Name)
		if _, err := client.DeleteAtManagementGroup(ctx, id); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	default:
		id := policydefinitions.NewProviderPolicyDefinitionID(meta.(*clients.Client).Account.SubscriptionId, oldId.Name)
		if _, err := client.Delete(ctx, id); err != nil {
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	}

	return nil
}

func policyDefinitionRefreshFunc(ctx context.Context, client *policydefinitions.PolicyDefinitionsClient, id any) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, result, err := getPolicyDefinitionByID(ctx, client, id)
		if err != nil {
			return nil, strconv.Itoa(resp.StatusCode), fmt.Errorf("issuing read request in policyDefinitionRefreshFunc for %s: %+v", id, err)
		}

		return result, strconv.Itoa(resp.StatusCode), nil
	}
}

func getPolicyDefinition(ctx context.Context, client *policydefinitions.PolicyDefinitionsClient, id policydefinitions.ProviderPolicyDefinitionId) (*http.Response, *policydefinitions.PolicyDefinition, error) {
	builtinId := policydefinitions.NewPolicyDefinitionID(id.PolicyDefinitionName)
	builtInResp, err := client.GetBuiltIn(ctx, builtinId)
	if err != nil && !response.WasNotFound(builtInResp.HttpResponse) {
		return builtInResp.HttpResponse, nil, err
	}
	if !response.WasNotFound(builtInResp.HttpResponse) {
		return builtInResp.HttpResponse, builtInResp.Model, nil
	}

	resp, err := client.Get(ctx, id)
	return resp.HttpResponse, resp.Model, err
}

func flattenJSON(stringMap interface{}) string {
	if stringMap != nil {
		if v, ok := stringMap.(*interface{}); ok {
			stringMap = *v
		}
		value := stringMap.(map[string]interface{})
		jsonString, err := pluginsdk.FlattenJsonToString(value)
		if err == nil {
			return jsonString
		}
	}

	return ""
}

func resourceArmPolicyDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(policydefinitions.PolicyTypeBuiltIn),
				string(policydefinitions.PolicyTypeCustom),
				string(policydefinitions.PolicyTypeNotSpecified),
				string(policydefinitions.PolicyTypeStatic),
			}, false),
		},

		"mode": {
			Type:     pluginsdk.TypeString,
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
			Type:     pluginsdk.TypeString,
			Optional: true,
			ForceNew: true,
		},

		"display_name": {
			Type:     pluginsdk.TypeString,
			Required: true,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"policy_rule": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},

		"metadata": metadataSchema(),
	}
}
