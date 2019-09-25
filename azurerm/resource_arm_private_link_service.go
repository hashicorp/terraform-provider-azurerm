package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceArmPrivateLinkService() *schema.Resource {
	return &schema.Resource{
		Create: resourceArmPrivateLinkServiceCreateUpdate,
		Read:   resourceArmPrivateLinkServiceRead,
		Update: resourceArmPrivateLinkServiceCreateUpdate,
		Delete: resourceArmPrivateLinkServiceDelete,

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

			"auto_approval": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"fqdns": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},

			"ip_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"private_ip_allocation_method": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.Static),
								string(network.Dynamic),
							}, false),
							Default: string(network.Static),
						},
						"private_ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
						"private_ip_address_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPv4),
								string(network.IPv6),
							}, false),
							Default: string(network.IPv4),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"load_balancer_frontend_ip_configurations": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.NoEmptyStrings,
						},
					},
				},
			},

			"private_endpoint_connections": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"private_endpoint": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"location": azure.SchemaLocation(),
									"tags":     tags.Schema(),
								},
							},
						},
						"private_link_service_connection_state": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"action_required": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"description": {
										Type:     schema.TypeString,
										Optional: true,
									},
									"status": {
										Type:     schema.TypeString,
										Optional: true,
									},
								},
							},
						},
					},
				},
			},

			"visibility": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"subscriptions": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"alias": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network_interfaces": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},

			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},
	}
}

func resourceArmPrivateLinkServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		resp, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(resp.Response) {
				return fmt.Errorf("Error checking for present of existing Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
			}
		}
		if !utils.ResponseWasNotFound(resp.Response) {
			return tf.ImportAsExistsError("azurerm_private_link_service", *resp.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	autoApproval := d.Get("auto_approval").([]interface{})
	fqdns := d.Get("fqdns").([]interface{})
	ipConfigurations := d.Get("ip_configurations").([]interface{})
	loadBalancerFrontendIpConfigurations := d.Get("load_balancer_frontend_ip_configurations").([]interface{})
	privateEndpointConnections := d.Get("private_endpoint_connections").([]interface{})
	visibility := d.Get("visibility").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.PrivateLinkService{
		Location: utils.String(location),
		PrivateLinkServiceProperties: &network.PrivateLinkServiceProperties{
			AutoApproval:                         expandArmPrivateLinkServicePrivateLinkServicePropertiesAutoApproval(autoApproval),
			Fqdns:                                utils.ExpandStringSlice(fqdns),
			IPConfigurations:                     expandArmPrivateLinkServicePrivateLinkServiceIPConfiguration(ipConfigurations),
			LoadBalancerFrontendIPConfigurations: expandArmPrivateLinkServiceFrontendIPConfiguration(loadBalancerFrontendIpConfigurations),
			PrivateEndpointConnections:           expandArmPrivateLinkServicePrivateEndpointConnection(privateEndpointConnections),
			Visibility:                           expandArmPrivateLinkServicePrivateLinkServicePropertiesVisibility(visibility),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, resourceGroup, name, parameters)
	if err != nil {
		return fmt.Errorf("Error creating Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("Error waiting for creation of Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil {
		return fmt.Errorf("Cannot read Private Link Service %q (Resource Group %q) ID", name, resourceGroup)
	}
	d.SetId(*resp.ID)

	return resourceArmPrivateLinkServiceRead(d, meta)
}

func resourceArmPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateLinkServices"]

	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Link Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.Set("name", resp.Name)
	d.Set("resource_group_name", resourceGroup)
	if location := resp.Location; location != nil {
		d.Set("location", azure.NormalizeLocation(*location))
	}
	if privateLinkServiceProperties := resp.PrivateLinkServiceProperties; privateLinkServiceProperties != nil {
		d.Set("alias", privateLinkServiceProperties.Alias)
		if err := d.Set("auto_approval", flattenArmPrivateLinkServicePrivateLinkServicePropertiesAutoApproval(privateLinkServiceProperties.AutoApproval)); err != nil {
			return fmt.Errorf("Error setting `auto_approval`: %+v", err)
		}
		d.Set("fqdns", utils.FlattenStringSlice(privateLinkServiceProperties.Fqdns))
		if err := d.Set("ip_configurations", flattenArmPrivateLinkServicePrivateLinkServiceIPConfiguration(privateLinkServiceProperties.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `ip_configurations`: %+v", err)
		}
		if err := d.Set("load_balancer_frontend_ip_configurations", flattenArmPrivateLinkServiceFrontendIPConfiguration(privateLinkServiceProperties.LoadBalancerFrontendIPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `load_balancer_frontend_ip_configurations`: %+v", err)
		}
		if err := d.Set("network_interfaces", flattenArmPrivateLinkServiceInterface(privateLinkServiceProperties.NetworkInterfaces)); err != nil {
			return fmt.Errorf("Error setting `network_interfaces`: %+v", err)
		}
		if err := d.Set("private_endpoint_connections", flattenArmPrivateLinkServicePrivateEndpointConnection(privateLinkServiceProperties.PrivateEndpointConnections)); err != nil {
			return fmt.Errorf("Error setting `private_endpoint_connections`: %+v", err)
		}
		if err := d.Set("visibility", flattenArmPrivateLinkServicePrivateLinkServicePropertiesVisibility(privateLinkServiceProperties.Visibility)); err != nil {
			return fmt.Errorf("Error setting `visibility`: %+v", err)
		}
	}
	d.Set("type", resp.Type)

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPrivateLinkServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	id, err := azure.ParseAzureResourceID(d.Id())
	if err != nil {
		return err
	}
	resourceGroup := id.ResourceGroup
	name := id.Path["privateLinkServices"]

	future, err := client.Delete(ctx, resourceGroup, name)
	if err != nil {
		if response.WasNotFound(future.Response()) {
			return nil
		}
		return fmt.Errorf("Error deleting Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		if !response.WasNotFound(future.Response()) {
			return fmt.Errorf("Error waiting for deleting Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandArmPrivateLinkServicePrivateLinkServicePropertiesAutoApproval(input []interface{}) *network.PrivateLinkServicePropertiesAutoApproval {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	subscriptions := v["subscriptions"].([]interface{})

	result := network.PrivateLinkServicePropertiesAutoApproval{
		Subscriptions: utils.ExpandStringSlice(subscriptions),
	}
	return &result
}

func expandArmPrivateLinkServicePrivateLinkServiceIPConfiguration(input []interface{}) *[]network.PrivateLinkServiceIPConfiguration {
	results := make([]network.PrivateLinkServiceIPConfiguration, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		privateIpAddress := v["private_ip_address"].(string)
		privateIPAllocationMethod := v["private_ip_allocation_method"].(string)
		subnetId := v["subnet_id"].(string)
		privateIpAddressVersion := v["private_ip_address_version"].(string)
		name := v["name"].(string)

		result := network.PrivateLinkServiceIPConfiguration{
			Name: utils.String(name),
			PrivateLinkServiceIPConfigurationProperties: &network.PrivateLinkServiceIPConfigurationProperties{
				PrivateIPAddress:          utils.String(privateIpAddress),
				PrivateIPAddressVersion:   network.IPVersion(privateIpAddressVersion),
				PrivateIPAllocationMethod: network.IPAllocationMethod(privateIPAllocationMethod),
				Subnet: &network.Subnet{
					ID: utils.String(subnetId),
				},
			},
		}

		results = append(results, result)
	}
	return &results
}

func expandArmPrivateLinkServiceFrontendIPConfiguration(input []interface{}) *[]network.FrontendIPConfiguration {
	results := make([]network.FrontendIPConfiguration, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		id := v["id"].(string)

		result := network.FrontendIPConfiguration{
			ID: utils.String(id),
		}

		results = append(results, result)
	}
	return &results
}

func expandArmPrivateLinkServicePrivateEndpointConnection(input []interface{}) *[]network.PrivateEndpointConnection {
	results := make([]network.PrivateEndpointConnection, 0)
	for _, item := range input {
		v := item.(map[string]interface{})
		id := v["id"].(string)
		privateEndpoint := v["private_endpoint"].([]interface{})
		privateLinkServiceConnectionState := v["private_link_service_connection_state"].([]interface{})
		name := v["name"].(string)

		result := network.PrivateEndpointConnection{
			ID:   utils.String(id),
			Name: utils.String(name),
			PrivateEndpointConnectionProperties: &network.PrivateEndpointConnectionProperties{
				PrivateEndpoint:                   expandArmPrivateLinkServicePrivateEndpoint(privateEndpoint),
				PrivateLinkServiceConnectionState: expandArmPrivateLinkServicePrivateLinkServiceConnectionState(privateLinkServiceConnectionState),
			},
		}

		results = append(results, result)
	}
	return &results
}

func expandArmPrivateLinkServicePrivateLinkServicePropertiesVisibility(input []interface{}) *network.PrivateLinkServicePropertiesVisibility {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	subscriptions := v["subscriptions"].([]interface{})

	result := network.PrivateLinkServicePropertiesVisibility{
		Subscriptions: utils.ExpandStringSlice(subscriptions),
	}
	return &result
}

func expandArmPrivateLinkServicePrivateEndpoint(input []interface{}) *network.PrivateEndpoint {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	id := v["id"].(string)
	location := azure.NormalizeLocation(v["location"].(string))
	t := v["tags"].(map[string]interface{})

	result := network.PrivateEndpoint{
		ID:       utils.String(id),
		Location: utils.String(location),
		Tags:     tags.Expand(t),
	}
	return &result
}

func expandArmPrivateLinkServicePrivateLinkServiceConnectionState(input []interface{}) *network.PrivateLinkServiceConnectionState {
	if len(input) == 0 {
		return nil
	}
	v := input[0].(map[string]interface{})

	status := v["status"].(string)
	description := v["description"].(string)
	actionRequired := v["action_required"].(string)

	result := network.PrivateLinkServiceConnectionState{
		ActionRequired: utils.String(actionRequired),
		Description:    utils.String(description),
		Status:         utils.String(status),
	}
	return &result
}

func flattenArmPrivateLinkServicePrivateLinkServicePropertiesAutoApproval(input *network.PrivateLinkServicePropertiesAutoApproval) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["subscriptions"] = utils.FlattenStringSlice(input.Subscriptions)

	return []interface{}{result}
}

func flattenArmPrivateLinkServicePrivateLinkServiceIPConfiguration(input *[]network.PrivateLinkServiceIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if name := item.Name; name != nil {
			v["name"] = *name
		}
		if privateLinkServiceIPConfigurationProperties := item.PrivateLinkServiceIPConfigurationProperties; privateLinkServiceIPConfigurationProperties != nil {
			v["private_ip_allocation_method"] = string(privateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod)
			if privateIpAddress := privateLinkServiceIPConfigurationProperties.PrivateIPAddress; privateIpAddress != nil {
				v["private_ip_address"] = *privateIpAddress
			}
			v["private_ip_address_version"] = string(privateLinkServiceIPConfigurationProperties.PrivateIPAddressVersion)
			if subnet := privateLinkServiceIPConfigurationProperties.Subnet; subnet != nil {
				if subnetId := subnet.ID; subnetId != nil {
					v["subnet_id"] = *subnetId
				}
			}
		}

		results = append(results, v)
	}

	return results
}

func flattenArmPrivateLinkServiceFrontendIPConfiguration(input *[]network.FrontendIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if id := item.ID; id != nil {
			v["id"] = *id
		}

		results = append(results, v)
	}

	return results
}

func flattenArmPrivateLinkServiceInterface(input *[]network.Interface) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if id := item.ID; id != nil {
			v["id"] = *id
		}

		results = append(results, v)
	}

	return results
}

func flattenArmPrivateLinkServicePrivateEndpointConnection(input *[]network.PrivateEndpointConnection) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		v := make(map[string]interface{})

		if id := item.ID; id != nil {
			v["id"] = *id
		}
		if privateEndpointConnectionProperties := item.PrivateEndpointConnectionProperties; privateEndpointConnectionProperties != nil {
			v["private_endpoint"] = flattenArmPrivateLinkServicePrivateEndpoint(privateEndpointConnectionProperties.PrivateEndpoint)
			v["private_link_service_connection_state"] = flattenArmPrivateLinkServicePrivateLinkServiceConnectionState(privateEndpointConnectionProperties.PrivateLinkServiceConnectionState)
		}

		results = append(results, v)
	}

	return results
}

func flattenArmPrivateLinkServicePrivateLinkServicePropertiesVisibility(input *network.PrivateLinkServicePropertiesVisibility) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	result["subscriptions"] = utils.FlattenStringSlice(input.Subscriptions)

	return []interface{}{result}
}

func flattenArmPrivateLinkServicePrivateEndpoint(input *network.PrivateEndpoint) []interface{} {
	if input == nil {
		return make([]interface{}, 0)
	}

	result := make(map[string]interface{})

	if id := input.ID; id != nil {
		result["id"] = *id
	}
	if location := input.Location; location != nil {
		result["location"] = azure.NormalizeLocation(*location)
	}

	return []interface{}{result}
}

func flattenArmPrivateLinkServicePrivateLinkServiceConnectionState(input *network.PrivateLinkServiceConnectionState) []interface{} {
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
