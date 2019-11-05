package azurerm

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2019-06-01/network"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/features"
	aznet "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network"
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
				ValidateFunc: aznet.ValidatePrivateLinkServiceName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"auto_approval_subscription_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.GUID,
				},
			},

			"visibility_subscription_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validate.GUID,
				},
			},

			// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
			// "fqdns": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Elem: &schema.Schema{
			// 		Type:         schema.TypeString,
			// 		ValidateFunc: validate.NoEmptyStrings,
			// 	},
			// },

			// Required by the API you can't create the resource without at least one ip configuration
			// I had to split the schema because if you attempt to change the primary attribute it
			// succeeds but doesn't change the attribute. Once primary is set it is set forever unless
			// you destroy the resource and recreate it.
			"primary_nat_ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: aznet.ValidatePrivateLinkServiceName,
						},
						"private_ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.IPv4Address,
						},
						// Only IPv4 is supported by the API, but I am exposing this
						// as they will support IPv6 in a future release.
						"private_ip_address_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPv4),
							}, false),
							Default: string(network.IPv4),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			"secondary_nat_ip_configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 7,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:         schema.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: aznet.ValidatePrivateLinkServiceName,
						},
						"private_ip_address": {
							Type:         schema.TypeString,
							Optional:     true,
							ValidateFunc: validate.IPv4Address,
						},
						// Only IPv4 is supported by the API, but I am exposing this
						// as they will support IPv6 in a future release.
						"private_ip_address_version": {
							Type:     schema.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPv4),
							}, false),
							Default: string(network.IPv4),
						},
						"subnet_id": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
					},
				},
			},

			// private_endpoint_connections have been removed but maybe useful for the end user to
			// understand the state of endpoints connected to this service, create a datasource?

			// Required by the API you can't create the resource without at least one load balancer id
			"load_balancer_frontend_ip_configuration_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
			},

			"alias": {
				Type:     schema.TypeString,
				Computed: true,
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

		CustomizeDiff: func(d *schema.ResourceDiff, v interface{}) error {
			if err := aznet.ValidatePrivateLinkNatIpConfiguration(d); err != nil {
				return err
			}

			return nil
		},
	}
}

func resourceArmPrivateLinkServiceCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateLinkServiceClient
	ctx := meta.(*ArmClient).StopContext

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if features.ShouldResourcesBeImported() && d.IsNewResource() {
		existing, err := client.Get(ctx, resourceGroup, name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("Error checking for presence of existing Private Link Service %q (Resource Group %q): %s", name, resourceGroup, err)
			}
		}
		if existing.ID != nil && *existing.ID != "" {
			return tf.ImportAsExistsError("azurerm_private_link_service", *existing.ID)
		}
	}

	location := azure.NormalizeLocation(d.Get("location").(string))
	autoApproval := d.Get("auto_approval_subscription_ids").([]interface{})
	// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
	//fqdns := d.Get("fqdns").([]interface{})
	primaryIpConfiguration := d.Get("primary_nat_ip_configuration").([]interface{})
	secondaryIpConfigurations := d.Get("secondary_nat_ip_configuration").([]interface{})
	loadBalancerFrontendIpConfigurations := d.Get("load_balancer_frontend_ip_configuration_ids").([]interface{})
	visibility := d.Get("visibility_subscription_ids").([]interface{})
	t := d.Get("tags").(map[string]interface{})

	parameters := network.PrivateLinkService{
		Location: utils.String(location),
		PrivateLinkServiceProperties: &network.PrivateLinkServiceProperties{
			AutoApproval: &network.PrivateLinkServicePropertiesAutoApproval{
				Subscriptions: utils.ExpandStringSlice(autoApproval),
			},
			Visibility: &network.PrivateLinkServicePropertiesVisibility{
				Subscriptions: utils.ExpandStringSlice(visibility),
			},
			IPConfigurations:                     expandArmPrivateLinkServiceIPConfiguration(primaryIpConfiguration, secondaryIpConfigurations),
			LoadBalancerFrontendIPConfigurations: expandArmPrivateLinkServiceFrontendIPConfiguration(loadBalancerFrontendIpConfigurations),
			//Fqdns:                                utils.ExpandStringSlice(fqdns),
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
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	d.SetId(*resp.ID)

	return resourceArmPrivateLinkServiceRead(d, meta)
}

func resourceArmPrivateLinkServiceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateLinkServiceClient
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
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		d.Set("alias", props.Alias)
		if props.AutoApproval.Subscriptions != nil {
			if err := d.Set("auto_approval_subscription_ids", utils.FlattenStringSlice(props.AutoApproval.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `auto_approval_subscription_ids`: %+v", err)
			}
		}
		if props.Visibility.Subscriptions != nil {
			if err := d.Set("visibility_subscription_ids", utils.FlattenStringSlice(props.Visibility.Subscriptions)); err != nil {
				return fmt.Errorf("Error setting `visibility_subscription_ids`: %+v", err)
			}
		}
		// currently not implemented yet, timeline unknown, exact purpose unknown, maybe coming to a future API near you
		// if props.Fqdns != nil {
		// 	if err := d.Set("fqdns", utils.FlattenStringSlice(props.Fqdns)); err != nil {
		// 		return fmt.Errorf("Error setting `fqdns`: %+v", err)
		// 	}
		// }
		if props.IPConfigurations != nil {
			primaryIpConfig, secondaryIpConfig := flattenArmPrivateLinkServiceIPConfiguration(props.IPConfigurations)
			if err := d.Set("primary_nat_ip_configuration", primaryIpConfig); err != nil {
				return fmt.Errorf("Error setting `primary_nat_ip_configuration`: %+v", err)
			}
			if err := d.Set("secondary_nat_ip_configuration", secondaryIpConfig); err != nil {
				return fmt.Errorf("Error setting `secondary_nat_ip_configuration`: %+v", err)
			}
		}
		if props.LoadBalancerFrontendIPConfigurations != nil {
			if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenArmPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
				return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
			}
		}
		if props.NetworkInterfaces != nil {
			if err := d.Set("network_interface_ids", flattenArmPrivateLinkServiceInterface(props.NetworkInterfaces)); err != nil {
				return fmt.Errorf("Error setting `network_interface_ids`: %+v", err)
			}
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourceArmPrivateLinkServiceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*ArmClient).Network.PrivateLinkServiceClient
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

func expandArmPrivateLinkServiceIPConfiguration(primaryInput []interface{}, secondaryInput []interface{}) *[]network.PrivateLinkServiceIPConfiguration {
	if len(primaryInput) == 0 {
		return nil
	}

	results := make([]network.PrivateLinkServiceIPConfiguration, 0)

	for _, item := range primaryInput {
		v := item.(map[string]interface{})
		privateIpAddress := v["private_ip_address"].(string)
		subnetId := v["subnet_id"].(string)
		privateIpAddressVersion := v["private_ip_address_version"].(string)
		name := v["name"].(string)

		result := network.PrivateLinkServiceIPConfiguration{
			Name: utils.String(name),
			PrivateLinkServiceIPConfigurationProperties: &network.PrivateLinkServiceIPConfigurationProperties{
				PrivateIPAddress:        utils.String(privateIpAddress),
				PrivateIPAddressVersion: network.IPVersion(privateIpAddressVersion),
				Subnet: &network.Subnet{
					ID: utils.String(subnetId),
				},
				Primary: utils.Bool(true),
			},
		}

		if privateIpAddress != "" {
			result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethod("Static")
		} else {
			result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethod("Dynamic")
		}

		results = append(results, result)
	}

	if len(secondaryInput) != 0 {
		for _, item := range secondaryInput {
			v := item.(map[string]interface{})
			privateIpAddress := v["private_ip_address"].(string)
			subnetId := v["subnet_id"].(string)
			privateIpAddressVersion := v["private_ip_address_version"].(string)
			name := v["name"].(string)

			result := network.PrivateLinkServiceIPConfiguration{
				Name: utils.String(name),
				PrivateLinkServiceIPConfigurationProperties: &network.PrivateLinkServiceIPConfigurationProperties{
					PrivateIPAddress:        utils.String(privateIpAddress),
					PrivateIPAddressVersion: network.IPVersion(privateIpAddressVersion),
					Subnet: &network.Subnet{
						ID: utils.String(subnetId),
					},
					Primary: utils.Bool(false),
				},
			}

			if privateIpAddress != "" {
				result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethod("Static")
			} else {
				result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethod("Dynamic")
			}

			results = append(results, result)
		}
	}

	return &results
}

func expandArmPrivateLinkServiceFrontendIPConfiguration(input []interface{}) *[]network.FrontendIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.FrontendIPConfiguration, 0)

	for _, item := range input {
		result := network.FrontendIPConfiguration{
			ID: utils.String(item.(string)),
		}

		results = append(results, result)
	}

	return &results
}

func flattenArmPrivateLinkServiceIPConfiguration(input *[]network.PrivateLinkServiceIPConfiguration) ([]interface{}, []interface{}) {
	primaryResults := make([]interface{}, 0)
	secondaryResults := make([]interface{}, 0)
	primary := false
	if input == nil {
		return primaryResults, secondaryResults
	}

	for _, item := range *input {
		b := make(map[string]interface{})

		if name := item.Name; name != nil {
			b["name"] = *name
		}
		if props := item.PrivateLinkServiceIPConfigurationProperties; props != nil {
			if v := props.PrivateIPAddress; v != nil {
				b["private_ip_address"] = *v
			}
			b["private_ip_address_version"] = string(props.PrivateIPAddressVersion)
			if v := props.Subnet; v != nil {
				if i := v.ID; i != nil {
					b["subnet_id"] = *i
				}
			}
			if v := props.Primary; v != nil {
				primary = *v
			}
		}

		if primary {
			primaryResults = append(primaryResults, b)
		} else {
			secondaryResults = append(secondaryResults, b)
		}
	}

	return primaryResults, secondaryResults
}

func flattenArmPrivateLinkServiceFrontendIPConfiguration(input *[]network.FrontendIPConfiguration) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.ID; id != nil {
			results = append(results, *id)
		}
	}

	return results
}

func flattenArmPrivateLinkServiceInterface(input *[]network.Interface) []string {
	results := make([]string, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.ID; id != nil {
			results = append(results, *id)
		}
	}

	return results
}
