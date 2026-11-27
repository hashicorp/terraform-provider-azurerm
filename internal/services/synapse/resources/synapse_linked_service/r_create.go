// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_linked_service

import (
	"encoding/json"
	"fmt"
	"maps"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	synapseClient2 "github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/jackofallops/kermit/sdk/synapse/2021-06-01-preview/synapse"
)

func resourceSynapseLinkedServiceCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	synapseClient := synapseClient2.From(meta)
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()
	environment := account.From(meta).Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	// todo 6.0 - move to the case-sensitive parser when validation.AsGeneratedID is removed: this parses a config
	// value which the paired AsGeneratedID validator accepts with legacy casing, and configs cannot be migrated.
	workspaceId, err := workspaces.ParseWorkspaceIDInsensitively(d.Get("synapse_workspace_id").(string))
	if err != nil {
		return err
	}

	client, err := synapseClient.LinkedServiceClient(workspaceId.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	id := parse.NewLinkedServiceID(workspaceId.SubscriptionId, workspaceId.ResourceGroupName, workspaceId.WorkspaceName, d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			existing, err := client.GetLinkedService(ctx, id.Name, "")
			if err != nil {
				if !response.WasNotFound(existing.Response.Response) {
					return fmt.Errorf("checking for presence of existing %s: %+v", id, err)
				}
			}
			if !response.WasNotFound(existing.Response.Response) {
				return tf.ImportAsExistsError("azurerm_synapse_linked_service", id.ID())
			}
		}
	}

	props := map[string]any{
		"type":       d.Get("type").(string),
		"connectVia": expandSynapseLinkedServiceIntegrationRuntimeV2(d.Get("integration_runtime").([]any)),
	}

	jsonDataStr := fmt.Sprintf(`{ "typeProperties": %s }`, d.Get("type_properties_json").(string))
	if err = json.Unmarshal([]byte(jsonDataStr), &props); err != nil {
		return err
	}

	if v, ok := d.GetOk("description"); ok {
		props["description"] = v.(string)
	}

	if v, ok := d.GetOk("parameters"); ok {
		props["parameters"] = expandSynapseParameters(v.(map[string]any))
	}

	if v, ok := d.GetOk("annotations"); ok {
		props["annotations"] = v.([]any)
	}

	additionalProperties := d.Get("additional_properties").(map[string]any)
	maps.Copy(props, additionalProperties)

	jsonData, err := json.Marshal(map[string]any{
		"properties": props,
	})
	if err != nil {
		return err
	}

	linkedService := &synapse.LinkedServiceResource{}
	if err := linkedService.UnmarshalJSON(jsonData); err != nil {
		return err
	}

	future, err := client.CreateOrUpdateLinkedService(ctx, id.Name, *linkedService, "")
	if err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting on creation for %s: %+v", id, err)
	}

	// Sometimes this resource fails to create but Azure is returning a 200. We'll check if the last response failed or not before moving on
	// todo remove this once https://github.com/hashicorp/go-azure-sdk/pull/122 is merged
	if err = checkLinkedServiceResponse(future.Response()); err != nil {
		return err
	}

	if d.IsNewResource() {
		d.SetId(id.ID())
	}

	return resourceSynapseLinkedServiceRead(d, meta)
}

func expandSynapseParameters(input map[string]any) map[string]*synapse.ParameterSpecification {
	output := make(map[string]*synapse.ParameterSpecification)

	for k, v := range input {
		output[k] = &synapse.ParameterSpecification{
			Type:         synapse.ParameterTypeString,
			DefaultValue: v.(string),
		}
	}

	return output
}

func expandSynapseLinkedServiceIntegrationRuntimeV2(input []any) *synapse.IntegrationRuntimeReference {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	return &synapse.IntegrationRuntimeReference{
		ReferenceName: pointer.To(v["name"].(string)),
		Type:          pointer.To("IntegrationRuntimeReference"),
		Parameters:    v["parameters"].(map[string]any),
	}
}
