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

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

// BEGIN
// TODO: Remove from here until the `END` comment post 5.0
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
			Type:       pluginsdk.TypeString,
			Optional:   true,
			ForceNew:   true,
			Deprecated: "`management_group_id` has been deprecated in favour of the `azurerm_management_group_policy_definition` resource and will be removed in v5.0 of the AzureRM Provider.",
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

// END TODO: Remove post 5.0

type PolicyDefinitionResource struct{}

type PolicyDefinitionResourceModel struct {
	Name        string `tfschema:"name"`
	PolicyType  string `tfschema:"policy_type"`
	Mode        string `tfschema:"mode"`
	DisplayName string `tfschema:"display_name"`
	Description string `tfschema:"description"`
	Metadata    string `tfschema:"metadata"`
	Parameters  string `tfschema:"parameters"`
	PolicyRule  string `tfschema:"policy_rule"`
}

var (
	_ sdk.ResourceWithUpdate         = PolicyDefinitionResource{}
	_ sdk.ResourceWithStateMigration = PolicyDefinitionResource{}
	_ sdk.ResourceWithCustomizeDiff  = PolicyDefinitionResource{}
)

func (r PolicyDefinitionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.PolicyDefinitionV0ToV1{},
		},
	}
}

func (r PolicyDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"policy_type": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringInSlice(policydefinitions.PossibleValuesForPolicyType(), false),
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

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:         pluginsdk.TypeString,
			Optional:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"metadata": metadataSchema(),

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"policy_rule": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},
	}
}

func (r PolicyDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"role_definition_ids": {
			Type:     pluginsdk.TypeList,
			Computed: true,
			Elem: &pluginsdk.Schema{
				Type: pluginsdk.TypeString,
			},
		},
	}
}

func (r PolicyDefinitionResource) ModelObject() interface{} {
	return &PolicyDefinitionResourceModel{}
}

func (r PolicyDefinitionResource) ResourceType() string {
	return "azurerm_policy_definition"
}

func (r PolicyDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model PolicyDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := policydefinitions.NewProviderPolicyDefinitionID(subscriptionId, model.Name)

			resp, _, err := getPolicyDefinition(ctx, client, id)
			if err != nil && !response.WasNotFound(resp) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(resp) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := policydefinitions.PolicyDefinition{
				Name: pointer.To(model.Name),
				Properties: &policydefinitions.PolicyDefinitionProperties{
					Description: pointer.To(model.Description),
					DisplayName: pointer.To(model.DisplayName),
					PolicyType:  pointer.To(policydefinitions.PolicyType(model.PolicyType)),
					Mode:        pointer.To(model.Mode),
				},
			}
			props := parameters.Properties

			if model.PolicyRule != "" {
				policyRule, err := pluginsdk.ExpandJsonFromString(model.PolicyRule)
				if err != nil {
					return fmt.Errorf("expanding `policy_rule`: %+v", err)
				}

				var iPolicyRule interface{} = policyRule
				props.PolicyRule = &iPolicyRule
			}

			if model.Metadata != "" {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(model.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata interface{} = expandedMetadata
				props.Metadata = &iMetadata
			}

			if model.Parameters != "" {
				expandedParameters, err := expandParameterDefinitionsValueForPolicyDefinition(model.Parameters)
				if err != nil {
					return fmt.Errorf("expanding `parameters`: %+v", err)
				}

				props.Parameters = expandedParameters
			}

			if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PolicyDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviderPolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, model, err := getPolicyDefinition(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PolicyDefinitionResourceModel{
				Name: id.PolicyDefinitionName,
			}

			if model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.PolicyType = string(pointer.From(props.PolicyType))
					state.Mode = pointer.From(props.Mode)

					if v, ok := pointer.From(props.Metadata).(map[string]interface{}); ok {
						flattenedMetadata, err := pluginsdk.FlattenJsonToString(v)
						if err != nil {
							return fmt.Errorf("flattening `metadata`: %+v", err)
						}

						state.Metadata = flattenedMetadata
					}

					if policyRule, ok := pointer.From(props.PolicyRule).(map[string]interface{}); ok {
						flattenedPolicyRule, err := pluginsdk.FlattenJsonToString(policyRule)
						if err != nil {
							return fmt.Errorf("flattening `policy_rule`: %+v", err)
						}
						state.PolicyRule = flattenedPolicyRule

						roleIDs, _ := getPolicyRoleDefinitionIDs(flattenedPolicyRule)
						metadata.ResourceData.Set("role_definition_ids", roleIDs)
					}

					flattenedParameters, err := flattenParameterDefinitionsValueToStringForPolicyDefinition(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `parameters`: %+v", err)
					}
					state.Parameters = flattenedParameters
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PolicyDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviderPolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config PolicyDefinitionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, model, err := getPolicyDefinition(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			if model == nil {
				return fmt.Errorf("retrieving %s: `model` was nil", *id)
			}

			if model.Properties == nil {
				return fmt.Errorf("retrieving %s: `properties` was nil", *id)
			}
			props := model.Properties

			if metadata.ResourceData.HasChange("display_name") {
				props.DisplayName = pointer.To(config.DisplayName)
			}

			if metadata.ResourceData.HasChange("description") {
				props.Description = pointer.To(config.Description)
			}

			if metadata.ResourceData.HasChange("mode") {
				props.Mode = pointer.To(config.Mode)
			}

			if metadata.ResourceData.HasChange("policy_rule") {
				props.PolicyRule = nil
				if config.PolicyRule != "" {
					policyRule, err := pluginsdk.ExpandJsonFromString(config.PolicyRule)
					if err != nil {
						return fmt.Errorf("expanding `policy_rule`: %+v", err)
					}

					var iPolicyRule interface{} = policyRule
					props.PolicyRule = &iPolicyRule
				}
			}

			if metadata.ResourceData.HasChange("metadata") {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(config.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata interface{} = expandedMetadata
				props.Metadata = &iMetadata
			}

			if metadata.ResourceData.HasChange("parameters") {
				props.Parameters = nil
				if config.Parameters != "" {
					expandedParameters, err := expandParameterDefinitionsValueForPolicyDefinition(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding `parameters`: %+v", err)
					}

					props.Parameters = expandedParameters
				}
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PolicyDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicyDefinitionsClient

			id, err := policydefinitions.ParseProviderPolicyDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			if _, err := client.Delete(ctx, *id); err != nil {
				return fmt.Errorf("deleting %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PolicyDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policydefinitions.ValidateProviderPolicyDefinitionID
}

func (r PolicyDefinitionResource) CustomizeDiff() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 10 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			if metadata.ResourceDiff.HasChange("parameters") {
				oldParametersRaw, newParametersRaw := metadata.ResourceDiff.GetChange("parameters")
				if oldParametersString := oldParametersRaw.(string); oldParametersString != "" {
					newParametersString := newParametersRaw.(string)
					if newParametersString == "" {
						return metadata.ResourceDiff.ForceNew("parameters")
					}

					oldParameters, err := expandParameterDefinitionsValueForPolicyDefinition(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := expandParameterDefinitionsValueForPolicyDefinition(newParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					if len(*newParameters) < len(*oldParameters) {
						return metadata.ResourceDiff.ForceNew("parameters")
					}
				}
			}

			return nil
		},
	}
}
