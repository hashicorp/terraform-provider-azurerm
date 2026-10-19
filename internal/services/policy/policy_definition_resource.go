// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policydefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/custompollers"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourcePolicyDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePolicyDefinitionCreateUpdate,
		Update: resourcePolicyDefinitionCreateUpdate,
		Read:   resourcePolicyDefinitionRead,
		Delete: resourcePolicyDefinitionDelete,

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

		Schema: resourcePolicyDefinitionSchema(),

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v any) error {
			// `parameters` cannot have values removed so we'll ForceNew if there are less parameters between Terraform runs
			if d.HasChange("parameters") {
				oldParametersRaw, newParametersRaw := d.GetChange("parameters")
				if oldParametersString := oldParametersRaw.(string); oldParametersString != "" {
					newParametersString := newParametersRaw.(string)
					if newParametersString == "" {
						return d.ForceNew("parameters")
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
						return d.ForceNew("parameters")
					}
				}
			}

			return nil
		}),
	}
}

func resourcePolicyDefinitionCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)

	var id resourceids.ResourceId
	if v, ok := d.GetOk("management_group_id"); ok {
		managementGroupId, err := mgmtGrpParse.ManagementGroupID(v.(string))
		if err != nil {
			return err
		}
		managementGroupName := managementGroupId.Name
		scopedId := policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, name)
		id = &scopedId
	} else {
		scopedId := policydefinitions.NewProviderPolicyDefinitionID(subscriptionId, name)
		id = &scopedId
	}

	if d.IsNewResource() && !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, _, err := getPolicyDefinitionByID(ctx, client, id)
		if err != nil && !response.WasNotFound(resp) {
			return fmt.Errorf("checking for presence of existing policy definition `%s`: %+v", id.ID(), err)
		}

		if !response.WasNotFound(resp) {
			return tf.ImportAsExistsError("azurerm_policy_definition", id.ID())
		}
	}

	properties := policydefinitions.PolicyDefinitionProperties{
		DisplayName: pointer.To(d.Get("display_name").(string)),
		Description: pointer.To(d.Get("description").(string)),
		PolicyType:  pointer.ToEnum[policydefinitions.PolicyType](d.Get("policy_type").(string)),
		Mode:        pointer.To(d.Get("mode").(string)),
	}

	if policyRuleString := d.Get("policy_rule").(string); policyRuleString != "" {
		policyRule, err := pluginsdk.ExpandJsonFromString(policyRuleString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `policy_rule`: %+v", err)
		}
		var iPolicyRule any = policyRule
		properties.PolicyRule = &iPolicyRule
	}

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
		}
		var iMetadata any = metaData
		properties.Metadata = &iMetadata
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		params, err := expandParameterDefinitionsValueForPolicyDefinition(parametersString)
		if err != nil {
			return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
		}
		properties.Parameters = params
	}

	parameters := policydefinitions.PolicyDefinition{
		Name:       pointer.To(name),
		Properties: &properties,
	}

	switch typedId := id.(type) {
	case *policydefinitions.ProviderPolicyDefinitionId:
		if _, err := client.CreateOrUpdate(ctx, *typedId, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", typedId, err)
		}
	case *policydefinitions.Providers2PolicyDefinitionId:
		if _, err := client.CreateOrUpdateAtManagementGroup(ctx, *typedId, parameters); err != nil {
			return fmt.Errorf("creating/updating %s: %+v", typedId, err)
		}
	default:
		return fmt.Errorf("unsupported ID type: %T", id)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for Policy Definition %q to become available", name)
	poller := custompollers.NewEventualConsistencyPoller(10, func(pollerCtx context.Context) (*http.Response, error) {
		resp, _, err := getPolicyDefinitionByID(pollerCtx, client, id)
		if err == nil && resp == nil {
			return nil, fmt.Errorf("retrieving %s: HTTP response was nil", id)
		}
		return resp, err
	}, &custompollers.EventualConsistencyPollerOptions{
		Interval:              10 * time.Second,
		TargetStatusCode:      pointer.To(http.StatusOK),
		RetryErrorStatusCodes: []int{http.StatusNotFound},
	})
	if err := poller.PollUntilDone(ctx); err != nil {
		return fmt.Errorf("waiting for Policy Definition %q to become available: %+v", name, err)
	}

	resourceId, err := parse.PolicyDefinitionID(id.ID())
	if err != nil {
		return err
	}

	d.SetId(resourceId.Id)

	return resourcePolicyDefinitionRead(d, meta)
}

func resourcePolicyDefinitionRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	oldId, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	var id resourceids.ResourceId
	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	switch scopeId := oldId.PolicyScopeId.(type) {
	case parse.ScopeAtManagementGroup:
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
		scopedId := policydefinitions.NewProviders2PolicyDefinitionID(managementGroupName, oldId.Name)
		id = &scopedId
	case parse.ScopeAtSubscription:
		scopedId := policydefinitions.NewProviderPolicyDefinitionID(scopeId.SubscriptionId, oldId.Name)
		id = &scopedId
	default:
		return fmt.Errorf("unsupported scope for policy definition ID `%s`: expected a subscription or management group scope", d.Id())
	}

	resp, model, err := getPolicyDefinitionByID(ctx, client, id)
	if err != nil {
		if response.WasNotFound(resp) {
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
			d.Set("policy_type", pointer.FromEnum(props.PolicyType))
			d.Set("mode", props.Mode)
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)

			if policyRuleStr := flattenJSON(pointer.From(props.PolicyRule)); policyRuleStr != "" {
				d.Set("policy_rule", policyRuleStr)
				roleIDs, _ := getPolicyRoleDefinitionIDs(policyRuleStr)
				d.Set("role_definition_ids", roleIDs)
			}

			if metadataStr := flattenJSON(pointer.From(props.Metadata)); metadataStr != "" {
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

func resourcePolicyDefinitionDelete(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Policy.PolicyDefinitionsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	oldId, err := parse.PolicyDefinitionID(d.Id())
	if err != nil {
		return err
	}

	switch scopeId := oldId.PolicyScopeId.(type) {
	case parse.ScopeAtManagementGroup:
		id := policydefinitions.NewProviders2PolicyDefinitionID(scopeId.ManagementGroupName, oldId.Name)
		if resp, err := client.DeleteAtManagementGroup(ctx, id); err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil
			}
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	case parse.ScopeAtSubscription:
		id := policydefinitions.NewProviderPolicyDefinitionID(scopeId.SubscriptionId, oldId.Name)
		if resp, err := client.Delete(ctx, id); err != nil {
			if response.WasNotFound(resp.HttpResponse) {
				return nil
			}
			return fmt.Errorf("deleting %s: %+v", id, err)
		}
	default:
		return fmt.Errorf("unsupported scope for policy definition ID `%s`: expected a subscription or management group scope", d.Id())
	}

	return nil
}

func flattenJSON(stringMap any) string {
	if stringMap != nil {
		if v, ok := stringMap.(*any); ok {
			stringMap = *v
		}
		value := stringMap.(map[string]any)
		jsonString, err := pluginsdk.FlattenJsonToString(value)
		if err == nil {
			return jsonString
		}
	}

	return ""
}

// The SDK normalizes parameter type casing, so ignore those differences without ignoring case-sensitive values elsewhere in the JSON for existing resources.
func policyDefinitionParametersDiffSuppressFunc(_, old, new string, _ *pluginsdk.ResourceData) bool {
	var oldParameters, newParameters map[string]any
	if err := json.Unmarshal([]byte(old), &oldParameters); err != nil {
		return false
	}
	if err := json.Unmarshal([]byte(new), &newParameters); err != nil {
		return false
	}

	for _, parameters := range []map[string]any{oldParameters, newParameters} {
		for _, parameter := range parameters {
			definition, ok := parameter.(map[string]any)
			if !ok {
				continue
			}
			parameterType, ok := definition["type"].(string)
			if !ok {
				continue
			}
			for _, knownType := range policydefinitions.PossibleValuesForParameterType() {
				if strings.EqualFold(parameterType, knownType) {
					definition["type"] = knownType
					break
				}
			}
		}
	}

	return reflect.DeepEqual(oldParameters, newParameters)
}

func resourcePolicyDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
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
			DiffSuppressFunc: policyDefinitionParametersDiffSuppressFunc,
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
