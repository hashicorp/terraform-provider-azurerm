// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package synapse_linked_service

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/synapse/2021-06-01/workspaces"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/synapse/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	artifacts "github.com/jackofallops/kermit/sdk/synapse/2021-06-01-preview/synapse"
)

func resourceSynapseLinkedServiceRead(d *pluginsdk.ResourceData, meta any) error {
	synapseClient := meta.(*clients.Client).Synapse
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()
	environment := meta.(*clients.Client).Account.Environment
	synapseDomainSuffix, ok := environment.Synapse.DomainSuffix()
	if !ok {
		return fmt.Errorf("could not determine Synapse domain suffix for environment %q", environment.Name)
	}

	id, err := parse.LinkedServiceID(d.Id())
	if err != nil {
		return err
	}

	client, err := synapseClient.LinkedServiceClient(id.WorkspaceName, *synapseDomainSuffix)
	if err != nil {
		return err
	}

	resp, err := client.GetLinkedService(ctx, id.Name, "")
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.Name)
	d.Set("synapse_workspace_id", workspaces.NewWorkspaceID(id.SubscriptionId, id.ResourceGroup, id.WorkspaceName).ID())

	byteArr, err := json.Marshal(resp.Properties)
	if err != nil {
		return err
	}

	var m map[string]*json.RawMessage
	if err = json.Unmarshal(byteArr, &m); err != nil {
		return err
	}

	description := ""
	if v, ok := m["description"]; ok && v != nil {
		if err := json.Unmarshal(*v, &description); err != nil {
			return err
		}
		delete(m, "description")
	}
	d.Set("description", description)

	t := ""
	if v, ok := m["type"]; ok && v != nil {
		if err := json.Unmarshal(*v, &t); err != nil {
			return err
		}
		delete(m, "type")
	}
	d.Set("type", t)

	annotations := make([]any, 0)
	if v, ok := m["annotations"]; ok && v != nil {
		if err := json.Unmarshal(*v, &annotations); err != nil {
			return err
		}
		delete(m, "annotations")
	}
	d.Set("annotations", annotations)

	parameters := make(map[string]*artifacts.ParameterSpecification)
	if v, ok := m["parameters"]; ok && v != nil {
		if err := json.Unmarshal(*v, &parameters); err != nil {
			return err
		}
		delete(m, "parameters")
	}
	if err := d.Set("parameters", flattenSynapseParameters(parameters)); err != nil {
		return fmt.Errorf("setting `parameters`: %+v", err)
	}

	var integrationRuntime *artifacts.IntegrationRuntimeReference
	if v, ok := m["connectVia"]; ok && v != nil {
		integrationRuntime = &artifacts.IntegrationRuntimeReference{}
		if err := json.Unmarshal(*v, &integrationRuntime); err != nil {
			return err
		}
		delete(m, "connectVia")
	}
	if err := d.Set("integration_runtime", flattenSynapseLinkedServiceIntegrationRuntimeV2(integrationRuntime)); err != nil {
		return fmt.Errorf("setting `integration_runtime`: %+v", err)
	}

	delete(m, "typeProperties")

	// set "additional_properties"
	additionalProperties := make(map[string]any)
	bytes, err := json.Marshal(m)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, &additionalProperties); err != nil {
		return err
	}
	d.Set("additional_properties", additionalProperties)

	return nil
}

func flattenSynapseParameters(input map[string]*artifacts.ParameterSpecification) map[string]any {
	output := make(map[string]any)

	for k, v := range input {
		if v != nil {
			// we only support string parameters at this time
			val, ok := v.DefaultValue.(string)
			if !ok {
				log.Printf("[DEBUG] Skipping parameter %q since it's not a string", k)
			}

			output[k] = val
		}
	}

	return output
}

func flattenSynapseLinkedServiceIntegrationRuntimeV2(input *artifacts.IntegrationRuntimeReference) []any {
	if input == nil {
		return []any{}
	}

	name := pointer.From(input.ReferenceName)

	return []any{
		map[string]any{
			"name":       name,
			"parameters": input.Parameters,
		},
	}
}

func suppressJsonOrderingDifference(_, old, new string, _ *pluginsdk.ResourceData) bool {
	return pluginsdk.NormalizeJson(old) == pluginsdk.NormalizeJson(new)
}

func checkLinkedServiceResponse(response *http.Response) error {
	respBody, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("reading status response body: %+v", err)
	}
	defer response.Body.Close()

	body := make(map[string]any)
	if err = json.Unmarshal(respBody, &body); err != nil {
		return fmt.Errorf("could not parse status response: %+v", err)
	}

	if statusRaw, ok := body["status"]; ok && statusRaw != nil {
		if status, ok := statusRaw.(string); ok {
			if status == "Failed" {
				if errorRaw, ok := body["error"]; ok && errorRaw != nil {
					if responseError, ok := errorRaw.(map[string]any); ok {
						if messageRaw, ok := responseError["message"]; ok && messageRaw != nil {
							if message, ok := messageRaw.(string); ok {
								return fmt.Errorf("creating/updating Linked Service: %s", message)
							}
						}
					}
				}
				// we are specifically checking for `error` in the payload but if the status is Failed, we should return what we know
				return fmt.Errorf("creating/updating Linked Service: %+v", body)
			}
		}
	}

	return nil
}
