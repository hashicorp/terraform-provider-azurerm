// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package policy

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/resources/2025-01-01/policysetdefinitions"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	mgmtGrpParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/policy/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

// BEGIN
// TODO: Remove from here until the `END` comment on ln836 post 5.0
func resourceArmPolicySetDefinition() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceArmPolicySetDefinitionCreate,
		Update: resourceArmPolicySetDefinitionUpdate,
		Read:   resourceArmPolicySetDefinitionRead,
		Delete: resourceArmPolicySetDefinitionDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PolicySetDefinitionID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: resourcePolicySetDefinitionSchema(),

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

func resourcePolicySetDefinitionSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"policy_type": {
			Type:     pluginsdk.TypeString,
			Required: true,
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(policysetdefinitions.PolicyTypeBuiltIn),
				string(policysetdefinitions.PolicyTypeCustom),
				string(policysetdefinitions.PolicyTypeNotSpecified),
				string(policysetdefinitions.PolicyTypeStatic),
			}, false),
		},

		"management_group_id": {
			Type:       pluginsdk.TypeString,
			Optional:   true,
			ForceNew:   true,
			Deprecated: "`management_group_id` has been deprecated in favour of the `azurerm_management_group_policy_set_definition` resource and will be removed in v5.0 of the AzureRM Provider.",
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"metadata": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: policySetDefinitionsMetadataDiffSuppressFunc,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		// lintignore: S013
		"policy_definition_reference": {
			Type:     pluginsdk.TypeList,
			Required: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"policy_definition_id": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.PolicyDefinitionID,
					},

					"parameter_values": {
						Type:             pluginsdk.TypeString,
						Optional:         true,
						ValidateFunc:     validation.StringIsJSON,
						DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
					},

					"reference_id": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						Computed: true,
					},

					"policy_group_names": {
						Type:     pluginsdk.TypeSet,
						Optional: true,
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.StringIsNotEmpty,
						},
					},

					"version": {
						Type:     pluginsdk.TypeString,
						Optional: true,
						// Note: O+C because Azure sets and returns the latest version if not specified in the request.
						Computed: true,
						ValidateFunc: validation.StringMatch(
							regexp.MustCompile(`^([1-9]\d*)\.(\d+|\*)(\.\*(-preview)?)?$`),
							"version must match `{major}.{minor}[.*][-preview]` where each segment is a number or an asterisk. The major version number must be greater than zero and `-preview` is optional.",
						),
					},
				},
			},
		},

		"policy_definition_group": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"name": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"display_name": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"category": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"description": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},

					"additional_metadata_resource_id": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						ValidateFunc: validation.StringIsNotEmpty,
					},
				},
			},
			Set: policySetDefinitionPolicyDefinitionGroupHash,
		},
	}
}

func resourceArmPolicySetDefinitionCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicySetDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	if v, ok := d.GetOk("management_group_id"); ok {
		return createForManagementGroup(ctx, client, d, meta, v.(string))
	}

	id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, d.Get("name").(string))

	resp, _, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		if !response.WasNotFound(resp) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(resp) {
		return tf.ImportAsExistsError("azurerm_policy_set_definition", id.ID())
	}

	parameters := policysetdefinitions.PolicySetDefinition{
		Properties: &policysetdefinitions.PolicySetDefinitionProperties{
			DisplayName: pointer.To(d.Get("display_name").(string)),
			Description: pointer.To(d.Get("description").(string)),
			PolicyType:  pointer.To(policysetdefinitions.PolicyType(d.Get("policy_type").(string))),
		},
	}
	props := parameters.Properties

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding `metadata`: %+v", err)
		}

		var iMetadata interface{} = metaData

		props.Metadata = &iMetadata
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		params, err := expandParameterDefinitionsValue(parametersString)
		if err != nil {
			return fmt.Errorf("expanding `parameters`: %+v", err)
		}
		props.Parameters = params
	}

	if v, ok := d.GetOk("policy_definition_reference"); ok {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(v.([]interface{}), d)
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		props.PolicyDefinitions = definitions
	}

	if v, ok := d.GetOk("policy_definition_group"); ok {
		props.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(v.(*pluginsdk.Set).List())
	}

	if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for %s to become available", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, id),
		MinTimeout:                10 * time.Second,
		Timeout:                   d.Timeout(pluginsdk.TimeoutCreate),
		ContinuousTargetOccurence: 10,
	}

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	resourceId, err := parse.PolicySetDefinitionID(id.ID())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", id.ID(), err)
	}

	d.SetId(resourceId.Id)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func createForManagementGroup(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, d *pluginsdk.ResourceData, meta any, managementGroupIdString string) error {
	managementGroupId, err := mgmtGrpParse.ManagementGroupID(managementGroupIdString)
	if err != nil {
		return err
	}

	id := policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupId.Name, d.Get("name").(string))

	resp, _, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		if !response.WasNotFound(resp) {
			return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(resp) {
		return tf.ImportAsExistsError("azurerm_policy_set_definition", id.ID())
	}

	parameters := policysetdefinitions.PolicySetDefinition{
		Properties: &policysetdefinitions.PolicySetDefinitionProperties{
			DisplayName: pointer.To(d.Get("display_name").(string)),
			Description: pointer.To(d.Get("description").(string)),
			PolicyType:  pointer.To(policysetdefinitions.PolicyType(d.Get("policy_type").(string))),
		},
	}
	props := parameters.Properties

	if metaDataString := d.Get("metadata").(string); metaDataString != "" {
		metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
		if err != nil {
			return fmt.Errorf("expanding `metadata`: %+v", err)
		}

		var iMetadata interface{} = metaData

		props.Metadata = &iMetadata
	}

	if parametersString := d.Get("parameters").(string); parametersString != "" {
		params, err := expandParameterDefinitionsValue(parametersString)
		if err != nil {
			return fmt.Errorf("expanding `parameters`: %+v", err)
		}
		props.Parameters = params
	}

	if v, ok := d.GetOk("policy_definition_reference"); ok {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(v.([]interface{}), d)
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		props.PolicyDefinitions = definitions
	}

	if v, ok := d.GetOk("policy_definition_group"); ok {
		props.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(v.(*pluginsdk.Set).List())
	}

	if _, err = client.CreateOrUpdateAtManagementGroup(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// Policy Definitions are eventually consistent; wait for them to stabilize
	log.Printf("[DEBUG] Waiting for %s to become available", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:                   []string{"404"},
		Target:                    []string{"200"},
		Refresh:                   policySetDefinitionRefreshFunc(ctx, client, id),
		MinTimeout:                10 * time.Second,
		ContinuousTargetOccurence: 10,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err = stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %+v", id, err)
	}

	resourceId, err := parse.PolicySetDefinitionID(id.ID())
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", id.ID(), err)
	}

	d.SetId(resourceId.Id)

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicySetDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	if v, ok := resourceId.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupId = mgmtGrpParse.NewManagementGroupId(v.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	if managementGroupName != "" {
		return updateForManagementGroup(ctx, client, d, meta, managementGroupId.ID())
	}

	id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, resourceId.Name)

	_, model, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	props := model.Properties

	if d.HasChange("policy_type") {
		props.PolicyType = pointer.To(policysetdefinitions.PolicyType(d.Get("policy_type").(string)))
	}

	if d.HasChange("display_name") {
		props.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("metadata") {
		metaDataString := d.Get("metadata").(string)
		if metaDataString != "" {
			metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
			}

			var iMetadata interface{} = metaData

			props.Metadata = &iMetadata
		} else {
			props.Metadata = nil
		}
	}

	if d.HasChange("parameters") {
		parametersString := d.Get("parameters").(string)
		if parametersString != "" {
			parameters, err := expandParameterDefinitionsValue(parametersString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
			}
			props.Parameters = parameters
		} else {
			props.Parameters = nil
		}
	}

	if d.HasChange("policy_definition_group") {
		props.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(d.Get("policy_definition_group").(*pluginsdk.Set).List())
	}

	if d.HasChange("policy_definition_reference") {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(d.Get("policy_definition_reference").([]interface{}), d)
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		props.PolicyDefinitions = definitions
	}

	if _, err = client.CreateOrUpdate(ctx, id, *model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func updateForManagementGroup(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, d *pluginsdk.ResourceData, meta any, managementGroupIdString string) error {
	resourceId, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupId, err := mgmtGrpParse.ManagementGroupID(managementGroupIdString)
	if err != nil {
		return fmt.Errorf("parsing %s: %+v", managementGroupIdString, err)
	}

	id := policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupId.Name, resourceId.Name)

	_, model, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}

	if model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}
	props := model.Properties

	if d.HasChange("policy_type") {
		props.PolicyType = pointer.To(policysetdefinitions.PolicyType(d.Get("policy_type").(string)))
	}

	if d.HasChange("display_name") {
		props.DisplayName = pointer.To(d.Get("display_name").(string))
	}

	if d.HasChange("description") {
		props.Description = pointer.To(d.Get("description").(string))
	}

	if d.HasChange("metadata") {
		metaDataString := d.Get("metadata").(string)
		if metaDataString != "" {
			metaData, err := pluginsdk.ExpandJsonFromString(metaDataString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `metadata`: %+v", err)
			}

			var iMetadata interface{} = metaData

			props.Metadata = &iMetadata
		} else {
			props.Metadata = nil
		}
	}

	if d.HasChange("parameters") {
		parametersString := d.Get("parameters").(string)
		if parametersString != "" {
			parameters, err := expandParameterDefinitionsValue(parametersString)
			if err != nil {
				return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
			}
			props.Parameters = parameters
		} else {
			props.Parameters = nil
		}
	}

	if d.HasChange("policy_definition_group") {
		props.PolicyDefinitionGroups = expandAzureRMPolicySetDefinitionPolicyGroups(d.Get("policy_definition_group").(*pluginsdk.Set).List())
	}

	if d.HasChange("policy_definition_reference") {
		definitions, err := expandAzureRMPolicySetDefinitionPolicyDefinitions(d.Get("policy_definition_reference").([]interface{}), d)
		if err != nil {
			return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
		}
		props.PolicyDefinitions = definitions
	}

	if _, err = client.CreateOrUpdateAtManagementGroup(ctx, id, *model); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	return resourceArmPolicySetDefinitionRead(d, meta)
}

func resourceArmPolicySetDefinitionRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicySetDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	if scopeId, ok := resourceId.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	if managementGroupName != "" {
		return readForManagementGroup(ctx, client, d, managementGroupId.ID())
	}

	id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, resourceId.Name)

	resp, model, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		if response.WasNotFound(resp) {
			log.Printf("[INFO] Error reading Policy Set Definition %q - removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model != nil {
		d.Set("name", model.Name)

		if props := model.Properties; props != nil {
			d.Set("policy_type", string(pointer.From(props.PolicyType)))
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)

			if iMetadata := props.Metadata; iMetadata != nil {
				metadata := *iMetadata
				if v, ok := metadata.(map[string]interface{}); ok {
					metadataStr, err := pluginsdk.FlattenJsonToString(v)
					if err != nil {
						return fmt.Errorf("flattening `metadata`: %+v", err)
					}
					d.Set("metadata", metadataStr)
				}
			}

			if parameters := props.Parameters; parameters != nil {
				parametersStr, err := flattenParameterDefinitionsValue(parameters)
				if err != nil {
					return fmt.Errorf("flattening `parameters`: %+v", err)
				}

				d.Set("parameters", parametersStr)
			}

			references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(props.PolicyDefinitions)
			if err != nil {
				return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
			}
			if err := d.Set("policy_definition_reference", references); err != nil {
				return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
			}

			if err := d.Set("policy_definition_group", flattenAzureRMPolicySetDefinitionPolicyGroups(props.PolicyDefinitionGroups)); err != nil {
				return fmt.Errorf("setting `policy_definition_group`: %+v", err)
			}
		}
	}

	return nil
}

func readForManagementGroup(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, d *pluginsdk.ResourceData, managementGroupIdString string) error {
	resourceId, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupId, err := mgmtGrpParse.ManagementGroupID(managementGroupIdString)
	if err != nil {
		return err
	}

	id := policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupId.Name, resourceId.Name)

	resp, model, err := getPolicySetDefinitionByID(ctx, client, id)
	if err != nil {
		if response.WasNotFound(resp) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if model != nil {
		d.Set("name", model.Name)
		d.Set("management_group_id", managementGroupIdString)

		if props := model.Properties; props != nil {
			d.Set("policy_type", string(pointer.From(props.PolicyType)))
			d.Set("display_name", props.DisplayName)
			d.Set("description", props.Description)

			if iMetadata := props.Metadata; iMetadata != nil {
				metadata := *iMetadata
				if v, ok := metadata.(map[string]interface{}); ok {
					metadataStr, err := pluginsdk.FlattenJsonToString(v)
					if err != nil {
						return fmt.Errorf("flattening `metadata`: %+v", err)
					}
					d.Set("metadata", metadataStr)
				}
			}

			if parameters := props.Parameters; parameters != nil {
				parametersStr, err := flattenParameterDefinitionsValue(parameters)
				if err != nil {
					return fmt.Errorf("flattening `parameters`: %+v", err)
				}

				d.Set("parameters", parametersStr)
			}

			references, err := flattenAzureRMPolicySetDefinitionPolicyDefinitions(props.PolicyDefinitions)
			if err != nil {
				return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
			}
			if err := d.Set("policy_definition_reference", references); err != nil {
				return fmt.Errorf("setting `policy_definition_reference`: %+v", err)
			}

			if err := d.Set("policy_definition_group", flattenAzureRMPolicySetDefinitionPolicyGroups(props.PolicyDefinitionGroups)); err != nil {
				return fmt.Errorf("setting `policy_definition_group`: %+v", err)
			}
		}
	}

	return nil
}

func resourceArmPolicySetDefinitionDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Policy.PolicySetDefinitionsClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId, err := parse.PolicySetDefinitionID(d.Id())
	if err != nil {
		return err
	}

	managementGroupName := ""
	var managementGroupId mgmtGrpParse.ManagementGroupId
	if scopeId, ok := resourceId.PolicyScopeId.(parse.ScopeAtManagementGroup); ok {
		managementGroupId = mgmtGrpParse.NewManagementGroupId(scopeId.ManagementGroupName)
		managementGroupName = managementGroupId.Name
	}

	if managementGroupName != "" {
		return deleteForManagementGroup(ctx, client, policysetdefinitions.NewProviders2PolicySetDefinitionID(managementGroupName, resourceId.Name).ID())
	}

	id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, resourceId.Name)

	if _, err := client.Delete(ctx, id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func deleteForManagementGroup(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, managementGroupIdString string) error {
	id, err := policysetdefinitions.ParseProviders2PolicySetDefinitionID(managementGroupIdString)
	if err != nil {
		return err
	}

	if _, err = client.DeleteAtManagementGroup(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandAzureRMPolicySetDefinitionPolicyDefinitions(input []interface{}, resourceData *pluginsdk.ResourceData) ([]policysetdefinitions.PolicyDefinitionReference, error) {
	result := make([]policysetdefinitions.PolicyDefinitionReference, 0)

	for idx, item := range input {
		v := item.(map[string]interface{})

		var parameters map[string]policysetdefinitions.ParameterValuesValue
		if p, ok := v["parameter_values"].(string); ok && p != "" {
			parameters = make(map[string]policysetdefinitions.ParameterValuesValue)
			if err := json.Unmarshal([]byte(p), &parameters); err != nil {
				return nil, fmt.Errorf("unmarshalling `parameter_values`: %+v", err)
			}
		}

		reference := policysetdefinitions.PolicyDefinitionReference{
			PolicyDefinitionId:          v["policy_definition_id"].(string),
			Parameters:                  pointer.To(parameters),
			PolicyDefinitionReferenceId: pointer.To(v["reference_id"].(string)),
			GroupNames:                  utils.ExpandStringSlice(v["policy_group_names"].(*pluginsdk.Set).List()),
		}

		// The API returns an error if we send an empty string
		if version := v["version"].(string); version != "" {
			path := fmt.Sprintf("policy_definition_reference.%d.version", idx)

			// We need to check the version value using RawConfig due to how Terraform manages blocks with computed nested items.
			// E.g. in a list of 3 `policy_definition_reference` blocks, if the middle (index 1) block is removed
			// the `version` argument contains the value from state as it's considered unchanged. However, due to the "shifted"
			// indexes, the `version` previously computed/returned by Azure may be incorrect for the `policy_definition_id`
			// it is now associated with, leading to a 400 Bad Request error from Azure.
			rawVersion, diags := resourceData.GetRawConfigAt(sdk.ConstructCtyPath(path))
			if diags.HasError() {
				return nil, fmt.Errorf("retrieving value at path `%s`: %+v", path, diags)
			}

			if !rawVersion.IsNull() {
				reference.DefinitionVersion = pointer.To(v["version"].(string))
			}
		}

		result = append(result, reference)
	}

	return result, nil
}

func flattenAzureRMPolicySetDefinitionPolicyDefinitions(input []policysetdefinitions.PolicyDefinitionReference) ([]interface{}, error) {
	result := make([]interface{}, 0)

	for _, definition := range input {
		parameterValues, err := flattenPolicyDefinitionReferenceParameterValues(definition.Parameters)
		if err != nil {
			return nil, fmt.Errorf("flattening `parameter_values`: %+v", err)
		}

		result = append(result, map[string]interface{}{
			"policy_definition_id": definition.PolicyDefinitionId,
			"parameter_values":     parameterValues,
			"reference_id":         pointer.From(definition.PolicyDefinitionReferenceId),
			"policy_group_names":   utils.FlattenStringSlice(definition.GroupNames),
			"version":              pointer.From(definition.DefinitionVersion),
		})
	}
	return result, nil
}

func expandAzureRMPolicySetDefinitionPolicyGroups(input []interface{}) *[]policysetdefinitions.PolicyDefinitionGroup {
	result := make([]policysetdefinitions.PolicyDefinitionGroup, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		group := policysetdefinitions.PolicyDefinitionGroup{}
		if name := v["name"].(string); name != "" {
			group.Name = name
		}
		if displayName := v["display_name"].(string); displayName != "" {
			group.DisplayName = pointer.To(displayName)
		}
		if category := v["category"].(string); category != "" {
			group.Category = pointer.To(category)
		}
		if description := v["description"].(string); description != "" {
			group.Description = pointer.To(description)
		}
		if metadataID := v["additional_metadata_resource_id"].(string); metadataID != "" {
			group.AdditionalMetadataId = pointer.To(metadataID)
		}
		result = append(result, group)
	}

	return &result
}

func flattenAzureRMPolicySetDefinitionPolicyGroups(input *[]policysetdefinitions.PolicyDefinitionGroup) []interface{} {
	result := make([]interface{}, 0)
	if input == nil {
		return result
	}

	for _, group := range *input {
		result = append(result, map[string]interface{}{
			"name":                            group.Name,
			"display_name":                    pointer.From(group.DisplayName),
			"category":                        pointer.From(group.Category),
			"description":                     pointer.From(group.Description),
			"additional_metadata_resource_id": pointer.From(group.AdditionalMetadataId),
		})
	}

	return result
}

// END TODO: Remove post 5.0

type PolicySetDefinitionResource struct{}

type PolicySetDefinitionResourceModel struct {
	Name                      string                           `tfschema:"name"`
	PolicyType                string                           `tfschema:"policy_type"`
	DisplayName               string                           `tfschema:"display_name"`
	Description               string                           `tfschema:"description"`
	Metadata                  string                           `tfschema:"metadata"`
	Parameters                string                           `tfschema:"parameters"`
	PolicyDefinitionReference []PolicyDefinitionReferenceModel `tfschema:"policy_definition_reference"`
	PolicyDefinitionGroup     []PolicyDefinitionGroupModel     `tfschema:"policy_definition_group"`
}

var (
	_ sdk.ResourceWithUpdate         = PolicySetDefinitionResource{}
	_ sdk.ResourceWithStateMigration = PolicySetDefinitionResource{}
	_ sdk.ResourceWithCustomizeDiff  = PolicySetDefinitionResource{}
)

func (r PolicySetDefinitionResource) StateUpgraders() sdk.StateUpgradeData {
	return sdk.StateUpgradeData{
		SchemaVersion: 1,
		Upgraders: map[int]pluginsdk.StateUpgrade{
			0: migration.PolicySetDefinitionV0ToV1{},
		},
	}
}

func (r PolicySetDefinitionResource) Arguments() map[string]*pluginsdk.Schema {
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
			ValidateFunc: validation.StringInSlice(policysetdefinitions.PossibleValuesForPolicyType(), false),
		},

		"display_name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"description": {
			Type:     pluginsdk.TypeString,
			Optional: true,
		},

		"metadata": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			Computed:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: policySetDefinitionsMetadataDiffSuppressFunc,
		},

		"parameters": {
			Type:             pluginsdk.TypeString,
			Optional:         true,
			ValidateFunc:     validation.StringIsJSON,
			DiffSuppressFunc: pluginsdk.SuppressJsonDiff,
		},

		"policy_definition_reference": policyDefinitionReferenceSchema(),

		"policy_definition_group": policyDefinitionGroupSchema(),
	}
}

func (r PolicySetDefinitionResource) Attributes() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{}
}

func (r PolicySetDefinitionResource) ModelObject() interface{} {
	return &PolicySetDefinitionResourceModel{}
}

func (r PolicySetDefinitionResource) ResourceType() string {
	return "azurerm_policy_set_definition"
}

func (r PolicySetDefinitionResource) Create() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient
			subscriptionId := metadata.Client.Account.SubscriptionId

			var model PolicySetDefinitionResourceModel
			if err := metadata.Decode(&model); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			id := policysetdefinitions.NewProviderPolicySetDefinitionID(subscriptionId, model.Name)

			resp, _, err := getPolicySetDefinition(ctx, client, id)
			if err != nil && !response.WasNotFound(resp) {
				return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
			}

			if !response.WasNotFound(resp) {
				return metadata.ResourceRequiresImport(r.ResourceType(), id)
			}

			parameters := policysetdefinitions.PolicySetDefinition{
				Name: pointer.To(model.Name),
				Properties: &policysetdefinitions.PolicySetDefinitionProperties{
					Description: pointer.To(model.Description),
					DisplayName: pointer.To(model.DisplayName),
					PolicyType:  pointer.To(policysetdefinitions.PolicyType(model.PolicyType)),
				},
			}

			props := parameters.Properties
			if model.Metadata != "" {
				expandedMetadata, err := pluginsdk.ExpandJsonFromString(model.Metadata)
				if err != nil {
					return fmt.Errorf("expanding `metadata`: %+v", err)
				}

				var iMetadata interface{} = expandedMetadata

				props.Metadata = &iMetadata
			}

			if model.Parameters != "" {
				expandedParameters, err := expandParameterDefinitionsValue(model.Parameters)
				if err != nil {
					return fmt.Errorf("expanding `parameters`: %+v", err)
				}
				props.Parameters = expandedParameters
			}

			if len(model.PolicyDefinitionReference) > 0 {
				expandedDefinitions, err := expandPolicyDefinitionReference(model.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if len(model.PolicyDefinitionGroup) > 0 {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(model.PolicyDefinitionGroup)
			}

			if _, err = client.CreateOrUpdate(ctx, id, parameters); err != nil {
				return fmt.Errorf("creating %s: %+v", id, err)
			}

			metadata.SetID(id)
			return nil
		},
	}
}

func (r PolicySetDefinitionResource) Read() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 5 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			resp, model, err := getPolicySetDefinition(ctx, client, *id)
			if err != nil {
				if response.WasNotFound(resp) {
					return metadata.MarkAsGone(id)
				}

				return fmt.Errorf("retrieving %s: %+v", *id, err)
			}

			state := PolicySetDefinitionResourceModel{
				Name: id.PolicySetDefinitionName,
			}

			if model != nil {
				if props := model.Properties; props != nil {
					state.Description = pointer.From(props.Description)
					state.DisplayName = pointer.From(props.DisplayName)
					state.PolicyType = string(pointer.From(props.PolicyType))

					if v, ok := pointer.From(props.Metadata).(map[string]interface{}); ok {
						flattenedMetadata, err := pluginsdk.FlattenJsonToString(v)
						if err != nil {
							return fmt.Errorf("flattening `metadata`: %+v", err)
						}
						state.Metadata = flattenedMetadata
					}

					flattenedParameters, err := flattenParameterDefinitionsValue(props.Parameters)
					if err != nil {
						return fmt.Errorf("flattening `parameters`: %+v", err)
					}
					state.Parameters = flattenedParameters

					flattenedDefinitions, err := flattenPolicyDefinitionReference(props.PolicyDefinitions)
					if err != nil {
						return fmt.Errorf("flattening `policy_definition_reference`: %+v", err)
					}
					state.PolicyDefinitionReference = flattenedDefinitions

					state.PolicyDefinitionGroup = flattenPolicyDefinitionGroup(props.PolicyDefinitionGroups)
				}
			}

			return metadata.Encode(&state)
		},
	}
}

func (r PolicySetDefinitionResource) Update() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
			if err != nil {
				return err
			}

			var config PolicySetDefinitionResourceModel
			if err := metadata.Decode(&config); err != nil {
				return fmt.Errorf("decoding: %+v", err)
			}

			resp, model, err := getPolicySetDefinition(ctx, client, *id)
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
					expandedParameters, err := expandParameterDefinitionsValue(config.Parameters)
					if err != nil {
						return fmt.Errorf("expanding `parameters`: %+v", err)
					}
					props.Parameters = expandedParameters
				}
			}

			if metadata.ResourceData.HasChange("policy_definition_reference") {
				expandedDefinitions, err := expandPolicyDefinitionReference(config.PolicyDefinitionReference, metadata)
				if err != nil {
					return fmt.Errorf("expanding `policy_definition_reference`: %+v", err)
				}
				props.PolicyDefinitions = expandedDefinitions
			}

			if metadata.ResourceData.HasChange("policy_definition_group") {
				props.PolicyDefinitionGroups = expandPolicyDefinitionGroup(config.PolicyDefinitionGroup)
			}

			if _, err := client.CreateOrUpdate(ctx, *id, *model); err != nil {
				return fmt.Errorf("updating %s: %+v", *id, err)
			}

			return nil
		},
	}
}

func (r PolicySetDefinitionResource) Delete() sdk.ResourceFunc {
	return sdk.ResourceFunc{
		Timeout: 30 * time.Minute,
		Func: func(ctx context.Context, metadata sdk.ResourceMetaData) error {
			client := metadata.Client.Policy.PolicySetDefinitionsClient

			id, err := policysetdefinitions.ParseProviderPolicySetDefinitionID(metadata.ResourceData.Id())
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

func (r PolicySetDefinitionResource) IDValidationFunc() pluginsdk.SchemaValidateFunc {
	return policysetdefinitions.ValidateProviderPolicySetDefinitionID
}

func getPolicySetDefinition(ctx context.Context, client *policysetdefinitions.PolicySetDefinitionsClient, id policysetdefinitions.ProviderPolicySetDefinitionId) (*http.Response, *policysetdefinitions.PolicySetDefinition, error) {
	resp, err := client.GetBuiltIn(ctx, policysetdefinitions.NewPolicySetDefinitionID(id.PolicySetDefinitionName), policysetdefinitions.DefaultGetBuiltInOperationOptions())
	if response.WasNotFound(resp.HttpResponse) {
		resp, err := client.Get(ctx, id, policysetdefinitions.DefaultGetOperationOptions())
		return resp.HttpResponse, resp.Model, err
	}

	return resp.HttpResponse, resp.Model, err
}

func (r PolicySetDefinitionResource) CustomizeDiff() sdk.ResourceFunc {
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

					oldParameters, err := expandParameterDefinitionsValue(oldParametersString)
					if err != nil {
						return fmt.Errorf("expanding JSON for `parameters`: %+v", err)
					}

					newParameters, err := expandParameterDefinitionsValue(newParametersString)
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
