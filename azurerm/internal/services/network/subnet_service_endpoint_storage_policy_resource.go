package network

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/location"
	mgValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/managementgroup/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceSubnetServiceEndpointStoragePolicy() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceSubnetServiceEndpointStoragePolicyCreateUpdate,
		Read:   resourceSubnetServiceEndpointStoragePolicyRead,
		Update: resourceSubnetServiceEndpointStoragePolicyCreateUpdate,
		Delete: resourceSubnetServiceEndpointStoragePolicyDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.SubnetServiceEndpointStoragePolicyID(id)
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

			"resource_group_name": azure.SchemaResourceGroupName(),

			"location": location.Schema(),

			"definition": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MinItems: 1,
				MaxItems: 1,
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
								),
							},
						},

						"description": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(0, 140),
						},
					},
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceSubnetServiceEndpointStoragePolicyCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	resourceId := parse.NewSubnetServiceEndpointStoragePolicyID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))
	if d.IsNewResource() {
		resp, err := client.Get(ctx, resourceId.ResourceGroup, resourceId.ServiceEndpointPolicyName, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("checking for existing %s: %+v", resourceId, err)
			}
		}

		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_subnet_service_endpoint_storage_policy", resourceId.ID())
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	param := network.ServiceEndpointPolicy{
		Location: &location,
		ServiceEndpointPolicyPropertiesFormat: &network.ServiceEndpointPolicyPropertiesFormat{
			ServiceEndpointPolicyDefinitions: expandServiceEndpointPolicyDefinitions(d.Get("definition").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	future, err := client.CreateOrUpdate(ctx, resourceId.ResourceGroup, resourceId.ServiceEndpointPolicyName, param)
	if err != nil {
		return fmt.Errorf("creating Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", resourceId.ServiceEndpointPolicyName, resourceId.ResourceGroup, err)
	}

	if err := future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", resourceId.ServiceEndpointPolicyName, resourceId.ResourceGroup, err)
	}

	d.SetId(resourceId.ID())

	return resourceSubnetServiceEndpointStoragePolicyRead(d, meta)
}

func resourceSubnetServiceEndpointStoragePolicyRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubnetServiceEndpointStoragePolicyID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[DEBUG] Subnet Service Endpoint Storage Policy %q was not found in Resource Group %q - removing from state!", id.ServiceEndpointPolicyName, id.ResourceGroup)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", id.ServiceEndpointPolicyName, id.ResourceGroup, err)
	}

	d.Set("name", id.ServiceEndpointPolicyName)
	d.Set("resource_group_name", id.ResourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if prop := resp.ServiceEndpointPolicyPropertiesFormat; prop != nil {
		if err := d.Set("definition", flattenServiceEndpointPolicyDefinitions(prop.ServiceEndpointPolicyDefinitions)); err != nil {
			return fmt.Errorf("setting `definition`: %v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceSubnetServiceEndpointStoragePolicyDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.ServiceEndpointPoliciesClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.SubnetServiceEndpointStoragePolicyID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.ResourceGroup, id.ServiceEndpointPolicyName); err != nil {
		return fmt.Errorf("deleting Subnet Service Endpoint Storage Policy %q (Resource Group %q): %+v", id.ServiceEndpointPolicyName, id.ResourceGroup, err)
	}

	return nil
}

func expandServiceEndpointPolicyDefinitions(input []interface{}) *[]network.ServiceEndpointPolicyDefinition {
	if len(input) == 0 {
		return nil
	}

	output := make([]network.ServiceEndpointPolicyDefinition, 0)
	for _, e := range input {
		e := e.(map[string]interface{})
		output = append(output, network.ServiceEndpointPolicyDefinition{
			Name: utils.String(e["name"].(string)),
			ServiceEndpointPolicyDefinitionPropertiesFormat: &network.ServiceEndpointPolicyDefinitionPropertiesFormat{
				Description:      utils.String(e["description"].(string)),
				Service:          utils.String("Microsoft.Storage"),
				ServiceResources: utils.ExpandStringSlice(e["service_resources"].(*pluginsdk.Set).List()),
			},
		})
	}

	return &output
}

func flattenServiceEndpointPolicyDefinitions(input *[]network.ServiceEndpointPolicyDefinition) []interface{} {
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
			serviceResource = []interface{}{}
		)
		if b := e.ServiceEndpointPolicyDefinitionPropertiesFormat; b != nil {
			if b.Description != nil {
				description = *b.Description
			}
			serviceResource = utils.FlattenStringSlice(b.ServiceResources)
		}

		output = append(output, map[string]interface{}{
			"name":              name,
			"description":       description,
			"service_resources": serviceResource,
		})
	}

	return output
}
