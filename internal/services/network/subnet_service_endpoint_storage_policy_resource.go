// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/serviceendpointpolicies"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	mgValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceSubnetServiceEndpointStoragePolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetServiceEndpointStoragePolicyCreate,
		Read:   resourceSubnetServiceEndpointStoragePolicyRead,
		Update: resourceSubnetServiceEndpointStoragePolicyUpdate,
		Delete: resourceSubnetServiceEndpointStoragePolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.SubnetServiceEndpointStoragePolicyName,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"location": commonschema.Location(),

			"definition": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 2,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validate.SubnetServiceEndpointStoragePolicyDefinitionName,
						},

						"service_resources": {
							Type:     pluginsdk.TypeSet,
							Required: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
								ValidateFunc: validation.Any(
									azure.ValidateResourceID,
									mgValidate.ManagementGroupID,
									validation.StringInSlice([]string{
										"/services/Azure",
										"/services/Azure/Batch",
										"/services/Azure/DataFactory",
										"/services/Azure/MachineLearning",
										"/services/Azure/ManagedInstance",
										"/services/Azure/WebPI",
									}, false),
								),
							},
						},

						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 140),
						},

						"service": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							Default:  "Microsoft.Storage",
							ValidateFunc: validation.StringInSlice([]string{
								"Microsoft.Storage",
								"Global",
							}, false),
						},
					},
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceSubnetServiceEndpointStoragePolicyCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := serviceendpointpolicies.NewServiceEndpointPolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	resp, err := client.Get(ctx, id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(resp.HttpResponse) {
			return fmt.Errorf("checking for existing %s: %+v", id, err)
		}
	}

	if !response.WasNotFound(resp.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_subnet_service_endpoint_storage_policy", id.ID())
	}

	param := serviceendpointpolicies.ServiceEndpointPolicy{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &serviceendpointpolicies.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, param); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}

func resourceSubnetServiceEndpointStoragePolicyUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	if existing.Model == nil {
		return fmt.Errorf("retrieving %s: `model` was nil", id)
	}
	if existing.Model.Properties == nil {
		return fmt.Errorf("retrieving %s: `properties` was nil", id)
	}

	payload := existing.Model

	if d.HasChange("definition") {
		payload.Properties = &serviceendpointpolicies.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]interface{})),
		}
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}

func resourceSubnetServiceEndpointStoragePolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, serviceendpointpolicies.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", id)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("name", id.ServiceEndpointPolicyName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
			if err := d.Set("definition", flattenServiceEndpointPolicyDefinitions(props.ServiceEndpointPolicyDefinitions)); err != nil {
				return fmt.Errorf("setting `definition`: %v", err)
			}
		}
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourceSubnetServiceEndpointStoragePolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPolicies
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := serviceendpointpolicies.ParseServiceEndpointPolicyID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandServiceEndpointPolicyDefinitions(input []interface{}) *[]serviceendpointpolicies.ServiceEndpointPolicyDefinition {
	if len(input) == 0 {
		return nil
	}

	output := make([]serviceendpointpolicies.ServiceEndpointPolicyDefinition, 0)
	for _, e := range input {
		e := e.(map[string]interface{})
		output = append(output, serviceendpointpolicies.ServiceEndpointPolicyDefinition{
			Name: pointer.To(e["name"].(string)),
			Properties: &serviceendpointpolicies.ServiceEndpointPolicyDefinitionPropertiesFormat{
				Description:      pointer.To(e["description"].(string)),
				Service:          pointer.To(e["service"].(string)),
				ServiceResources: utils.ExpandStringSlice(e["service_resources"].(*pluginsdk.Set).List()),
			},
		})
	}

	return &output
}

func flattenServiceEndpointPolicyDefinitions(input *[]serviceendpointpolicies.ServiceEndpointPolicyDefinition) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	output := make([]interface{}, 0)
	for _, e := range *input {
		name := ""
		if e.Name != nil {
			name = *e.Name
		}

		var (
			description     = ""
			service         = ""
			serviceResource = []interface{}{}
		)
		if b := e.Properties; b != nil {
			if b.Description != nil {
				description = *b.Description
			}
			serviceResource = utils.FlattenStringSlice(b.ServiceResources)
			if b.Service != nil {
				service = *b.Service
			}
		}

		output = append(output, map[string]interface{}{
			"name":              name,
			"description":       description,
			"service_resources": serviceResource,
			"service":           service,
		})
	}

	return output
}
