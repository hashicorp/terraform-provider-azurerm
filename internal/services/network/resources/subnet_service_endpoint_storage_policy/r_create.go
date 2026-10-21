// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package subnet_service_endpoint_storage_policy

import (
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetServiceEndpointStoragePolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := serviceendpointpolicies.NewServiceEndpointPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !meta.(*clients.Client).Features.SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
		resp, err := client.Get(ctx, id, serviceendpointpolicies.DefaultGetOperationOptions())
		if err != nil {
			if !response.WasNotFound(resp.HttpResponse) {
				return fmt.Errorf("checking for existing %s: %+v", id, err)
			}
		}

		if !response.WasNotFound(resp.HttpResponse) {
			return tf.ImportAsExistsError("azurerm_subnet_service_endpoint_storage_policy", id.ID())
		}
	}

	param := serviceendpointpolicies.ServiceEndpointPolicy{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &serviceendpointpolicies.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]any)),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]any)),
	}

	if err := client.CreateOrUpdateCallbackThenPoll(ctx, id, param, sdk.SetIDAndIdentityCallback(meta, &id, d)); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	if err := pluginsdk.SetResourceIdentityData(d, &id); err != nil {
		return err
	}

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}

func expandServiceEndpointPolicyDefinitions(input []any) *[]serviceendpointpolicies.ServiceEndpointPolicyDefinition {
	if len(input) == 0 {
		return nil
	}

	output := make([]serviceendpointpolicies.ServiceEndpointPolicyDefinition, 0)
	for _, e := range input {
		e := e.(map[string]any)
		output = append(output, serviceendpointpolicies.ServiceEndpointPolicyDefinition{
			Name: pointer.To(e["name"].(string)),
			Properties: &serviceendpointpolicies.ServiceEndpointPolicyDefinitionPropertiesFormat{
				Description:      pointer.To(e["description"].(string)),
				Service:          pointer.To(e["service"].(string)),
				ServiceResources: pluginsdk.ExpandStringSlice(e["service_resources"].(*pluginsdk.Set).List()),
			},
		})
	}

	return &output
}
