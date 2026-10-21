// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package data_factory_integration_runtime_self_hosted

import (
	"bytes"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/factories"
	"github.com/hashicorp/go-azure-sdk/resource-manager/datafactory/2018-06-01/integrationruntimes"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceDataFactoryIntegrationRuntimeSelfHostedRead(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).DataFactory.IntegrationRuntimesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := integrationruntimes.ParseIntegrationRuntimeID(d.Id())
	if err != nil {
		return err
	}

	dataFactoryId := factories.NewFactoryID(id.SubscriptionId, id.ResourceGroupName, id.FactoryName)

	resp, err := client.Get(ctx, *id, integrationruntimes.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.IntegrationRuntimeName)
	d.Set("data_factory_id", dataFactoryId.ID())

	if model := resp.Model; model != nil {
		runTime, ok := model.Properties.(integrationruntimes.SelfHostedIntegrationRuntime)
		if !ok {
			return fmt.Errorf("asserting `IntegrationRuntime` as `SelfHostedIntegrationRuntime` for %s", *id)
		}

		d.Set("description", runTime.Description)

		if props := runTime.TypeProperties; props != nil {
			d.Set("self_contained_interactive_authoring_enabled", pointer.From(props.SelfContainedInteractiveAuthoringEnabled))
			rbacAuthorization, ok := props.LinkedInfo.(integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization)
			if ok {
				if err := d.Set("rbac_authorization", pluginsdk.NewSet(resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash, flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(rbacAuthorization))); err != nil {
					return fmt.Errorf("setting `rbac_authorization`: %#v", err)
				}
			}

			// The ListAuthenticationKeys on integration runtime type Linked is not supported.
			// Only skip the call to ListAuthKeys if the linkedInfo is valid.
			if _, ok := props.LinkedInfo.(integrationruntimes.RawLinkedIntegrationRuntimeTypeImpl); !ok {
				return nil
			}
		}
	}

	keyResp, keyErr := client.ListAuthKeys(ctx, *id)
	if keyErr != nil {
		if response.WasNotFound(keyResp.HttpResponse) {
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Authorization Keys for %s: %+v", *id, keyErr)
	}

	if keyResp.Model == nil {
		return fmt.Errorf("retrieving Authorization Keys for %s: `model` was nil", *id)
	}
	keyModel := keyResp.Model

	d.Set("primary_authorization_key", keyModel.AuthKey1)
	d.Set("secondary_authorization_key", keyModel.AuthKey2)

	return nil
}

func flattenAzureRmDataFactoryIntegrationRuntimeSelfHostedTypePropertiesRbacAuthorization(input integrationruntimes.LinkedIntegrationRuntimeRbacAuthorization) []any {
	result := make(map[string]any)
	result["resource_id"] = input.ResourceId

	return []any{result}
}

func resourceDataFactoryIntegrationRuntimeSelfHostedRbacAuthorizationHash(v any) int {
	var buf bytes.Buffer

	if m, ok := v.(map[string]any); ok {
		if v, ok := m["resource_id"]; ok {
			fmt.Fprintf(&buf, "%s-", v.(string))
		}
	}

	return pluginsdk.HashString(buf.String())
}
