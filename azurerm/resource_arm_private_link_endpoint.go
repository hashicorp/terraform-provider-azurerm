package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-09-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznet "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
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
				ValidateFunc: aznet.ValidatePrivateLinkName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupNameDiffSuppress(),

			"subnet_id": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"private_service_connection": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: aznet.ValidatePrivateLinkName,
						},
						"is_manual_connection": {
							Type:     schema.TypeBool,
							Required: true,
							ForceNew: true,
						},
						"private_connection_resource_id": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"subresource_names": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.NoEmptyStrings,
							},
						},
						"request_message": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validation.StringLenBetween(1, 140),
						},
					},
				},
			},

			"network_interface_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: schema.HashString,
			},

			// tags has been removed
			// API Issue "Unable to remove Tags from Private Link Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467
		},
	}
}

func resourceArmPrivateLinkEndpointCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if err := aznet.ValidatePrivateLinkEndpointSettings(d); err != nil {
		return fmt.Errorf("Error validating the configuration for the Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}

		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_link_endpoint", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	privateServiceConnections := d.Get("private_service_connection").([]interface{})
	subnetId := d.Get("subnet_id").(string)

	parameters := network.PrivateEndpoint{
		Location: utils.String(location),
		PrivateEndpointProperties: &network.PrivateEndpointProperties{
			PrivateLinkServiceConnections:       expandArmPrivateLinkEndpointServiceConnection(privateServiceConnections, false),
			ManualPrivateLinkServiceConnections: expandArmPrivateLinkEndpointServiceConnection(privateServiceConnections, true),
			Subnet: &network.Subnet{
				ID: utils.String(subnetId),
			},
		},
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmPrivateLinkEndpointRead(d, meta)
}

func resourceArmPrivateLinkEndpointRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
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
			log.Printf("[INFO] Private Link Endpoint %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if props := resp.PrivateEndpointProperties; props != nil {
		if subnet := props.Subnet; subnet != nil {
			d.Set("subnet_id", subnet.ID)
		}
		if err := d.Set("private_service_connection", flattenArmPrivateLinkEndpointServiceConnection(props.PrivateLinkServiceConnections, props.ManualPrivateLinkServiceConnections)); err != nil {
			return fmt.Errorf("Error setting `private_service_connection`: %+v", err)
		}
		if err := d.Set("network_interface_ids", flattenArmPrivateLinkEndpointInterface(props.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
		}
	}

	// API Issue "Unable to remove Tags from Private Link Endpoint": https://github.com/Azure/azure-sdk-for-go/issues/6467
	return nil
}

func resourceArmPrivateLinkEndpointDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateEndpointClient
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
		return fmt.Errorf("Error deleting Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Private Link Endpoint %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmPrivateLinkEndpointServiceConnection(input []interface{}, parseManual bool) *[]network.PrivateLinkServiceConnection {
	results := make([]network.PrivateLinkServiceConnection, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		privateConnectonResourceId := v["private_connection_resource_id"].(string)
		subresourceNames := v["subresource_names"].([]interface{})
		requestMessage := v["request_message"].(string)
		isManual := v["is_manual_connection"].(bool)
		name := v["name"].(string)

		if isManual == parseManual {
			result := network.PrivateLinkServiceConnection{
				Name: utils.String(name),
				PrivateLinkServiceConnectionProperties: &network.PrivateLinkServiceConnectionProperties{
					GroupIds:             utils.ExpandStringSlice(subresourceNames),
					PrivateLinkServiceID: utils.String(privateConnectonResourceId),
				},
			}

			if requestMessage != "" {
				result.PrivateLinkServiceConnectionProperties.RequestMessage = utils.String(requestMessage)
			}

			results = append(results, result)
		}
	}

	return &results
}

func flattenArmPrivateLinkEndpointServiceConnection(serviceConnections *[]network.PrivateLinkServiceConnection, manualServiceConnections *[]network.PrivateLinkServiceConnection) []interface{} {
	results := make([]interface{}, 0)
	if serviceConnections == nil && manualServiceConnections == nil {
		return results
	}

	if serviceConnections != nil {
		for _, item := range *serviceConnections {
			result := make(map[string]interface{})

			result["is_manual_connection"] = false

			if v := item.Name; v != nil {
				result["name"] = *v
			}
			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.GroupIds; v != nil {
					result["subresource_names"] = utils.FlattenStringSlice(v)
				}
				if v := props.PrivateLinkServiceID; v != nil {
					result["private_connection_resource_id"] = *v
				}
			}
			results = append(results, result)
		}
	}

	if manualServiceConnections != nil {
		for _, item := range *manualServiceConnections {
			result := make(map[string]interface{})

			result["is_manual_connection"] = true

			if v := item.Name; v != nil {
				result["name"] = *v
			}
			if props := item.PrivateLinkServiceConnectionProperties; props != nil {
				if v := props.GroupIds; v != nil {
					result["subresource_names"] = utils.FlattenStringSlice(v)
				}
				if v := props.PrivateLinkServiceID; v != nil {
					result["private_connection_resource_id"] = *v
				}
				if v := props.RequestMessage; v != nil {
					result["request_message"] = *v
				}
			}
			results = append(results, result)
		}
	}

	return results
}

func flattenArmPrivateLinkEndpointInterface(input *[]network.Interface) *schema.Set {
	results := &schema.Set{F: schema.HashString}
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.ID; id != nil {
			results.Add(*id)
		}
	}

	return results
}
