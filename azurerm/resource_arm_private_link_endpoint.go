package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateLinkEndpoint() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateLinkEndpointCreateUpdate,
		Read:   resourceArmPrivateLinkEndpointRead,
		Update: resourceArmPrivateLinkEndpointCreateUpdate,
		Delete: resourceArmPrivateLinkEndpointDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validate.NoEmptyStrings,
			},

			"manual_private_link_service_connection": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"private_link_service_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state_action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_message": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Please approve my connection",
						},
					},
				},
			},

			"private_link_service_connection": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"private_link_service_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"group_ids": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"state_action_required": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_message": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "Please approve my connection",
						},
					},
				},
			},

			"network_interface_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateLinkEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_private_link_endpoint", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	manualPrivateLinkServiceConnections := d.Get("manual_private_link_service_connection").([]interface{})
	privateLinkServiceConnections := d.Get("private_link_service_connection").([]interface{})
	subnetId := d.Get("subnet_id").(string)
	t := d.Get("tags").(map[string]interface{})

	parameters := network.PrivateEndpoint{
		Location: utils.String(location),
		PrivateEndpointProperties: &network.PrivateEndpointProperties{
			ManualPrivateLinkServiceConnections: expandArmPrivateEndpointPrivateLinkServiceConnection(manualPrivateLinkServiceConnections),
			PrivateLinkServiceConnections:       expandArmPrivateEndpointPrivateLinkServiceConnection(privateLinkServiceConnections),
			Subnet: &network.Subnet{
				ID: utils.String(subnetId),
			},
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Private Endpoint %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmPrivateLinkEndpointRead(d, meta)
}

func resourceArmPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if privateEndpointProperties := resp.PrivateEndpointProperties; privateEndpointProperties != nil {
		if err := d.Set("manual_private_link_service_connection", flattenArmPrivateEndpointPrivateLinkServiceConnection(privateEndpointProperties.ManualPrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `manual_private_link_service_connection`: %+v", err)
		}
		if err := d.Set("network_interface_ids", flattenArmPrivateEndpointInterface(privateEndpointProperties.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
		}
		if err := d.Set("private_link_service_connection", flattenArmPrivateEndpointPrivateLinkServiceConnection(privateEndpointProperties.PrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `private_link_service_connection`: %+v", err)
		}
		if subnet := privateEndpointProperties.Subnet; subnet != nil {
			d.Set("subnet_id", subnet.ID)
		}
	}
	tags.FlattenAndSet(d, resp.Tags)

	return nil
}

func resourceArmPrivateLinkEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateEndpoints"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Private Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmPrivateEndpointPrivateLinkServiceConnection(input []interface{}) *[]network.PrivateLinkServiceConnection {
	results := make([]network.PrivateLinkServiceConnection, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		privateLinkServiceID := v["private_link_service_id"].(string)
		groupIds := v["group_ids"].([]interface{})
		requestMessage := v["request_message"].(string)
		name := v["name"].(string)

		result := network.PrivateLinkServiceConnection{
			Name: utils.String(name),
			PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
				GroupIds:             utils.ExpandStringSlice(groupIds),
				PrivateLinkServiceID: utils.String(privateLinkServiceID),
				RequestMessage:       utils.String(requestMessage),
			},
		}

		results = append(results, result)
	}
	return &results
}

func flattenArmPrivateEndpointPrivateLinkServiceConnection(input *[]network.PrivateLinkServiceConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		if privateLinkServiceConnectionProperties := item.PrivateLinkServiceConnectionProperties; privateLinkServiceConnectionProperties != nil {
			v["group_ids"] = utils.FlattenStringSlice(privateLinkServiceConnectionProperties.GroupIds)
			v["private_link_service_connection_state"] = flattenArmPrivateLinkServiceConnectionState(privateLinkServiceConnectionProperties.PrivateLinkServiceConnectionState)
			if privateLinkServiceId := privateLinkServiceConnectionProperties.PrivateLinkServiceID; privateLinkServiceId != nil {
				v["private_link_service_id"] = *privateLinkServiceId
			}
			if requestMessage := privateLinkServiceConnectionProperties.RequestMessage; requestMessage != nil {
				v["request_message"] = *requestMessage
			}
		}

		results = append(results, v)
	}

	return results
}

func flattenArmPrivateEndpointInterface(input *[]network.Interface) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	results := make([]string, 0)

	for _, item := range *input {
		if id := item.ID; id != nil {
			results = append(results, *id)
		}
	}

	return []interface{}{results}
}

func flattenArmPrivateLinkServiceConnectionState(input *network.PrivateLinkServiceConnectionState) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if actionRequired := input.ActionRequired; actionRequired != nil {
		result["action_required"] = *actionRequired
	}
	if description := input.Description; description != nil {
		result["description"] = *description
	}
	if status := input.Status; status != nil {
		result["status"] = *status
	}

	return []interface{}{result}
}
