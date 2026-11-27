// Copyright IBM Corp. 2014, 2026
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
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/sdk"
	networkClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceSubnetServiceEndpointStoragePolicyCreate(d *pluginsdk.ResourceData, meta any) error {
	client := networkClient.From(meta).ServiceEndpointPolicies
	subscriptionId := account.From(meta).SubscriptionId
	ctx, cancel := timeouts.ForCreate(account.StopContext(meta), d)
	defer cancel()

	id := serviceendpointpolicies.NewServiceEndpointPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
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
