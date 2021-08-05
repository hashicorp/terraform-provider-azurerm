package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2020-11-01/network"
	"github.com/hashicorp/go-azure-helpers/response"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	networkValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/network/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tags"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourcePrivateLinkService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateLinkServiceCreateUpdate,
		Read:   resourcePrivateLinkServiceRead,
		Update: resourcePrivateLinkServiceCreateUpdate,
		Delete: resourcePrivateLinkServiceDelete,

		// TODO: replace this with an importer which validates the ID during import
		Importer: pluginsdk.DefaultImporter(),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(60 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(60 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkValidate.PrivateLinkName,
			},

			"location": azure.SchemaLocation(),

			"resource_group_name": azure.SchemaResourceGroupName(),

			"auto_approval_subscription_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
				Set: pluginsdk.HashString,
			},

			"enable_proxy_protocol": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
			},

			"visibility_subscription_ids": {
				Type:     pluginsdk.TypeSet,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.IsUUID,
				},
				Set: pluginsdk.HashString,
			},

			// Required by the API you can't create the resource without at least
			// one ip configuration once primary is set it is set forever unless
			// you destroy the resource and recreate it.
			"nat_ip_configuration": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 8,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: networkValidate.PrivateLinkName,
						},
						"private_ip_address": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: validate.IPv4Address,
						},
						// Only IPv4 is supported by the API, but I am exposing this
						// as they will support IPv6 in a future release.
						"private_ip_address_version": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(network.IPVersionIPv4),
							}, false),
							Default: string(network.IPVersionIPv4),
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: azure.ValidateResourceID,
						},
						"primary": {
							Type:     pluginsdk.TypeBool,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},

			// Required by the API you can't create the resource without at least one load balancer id
			"load_balancer_frontend_ip_configuration_ids": {
				Type:     pluginsdk.TypeSet,
				Required: true,
				ForceNew: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: azure.ValidateResourceID,
				},
				Set: pluginsdk.HashString,
			},

			"alias": {
				Type:     pluginsdk.TypeString,
				Computed: true,
			},

			"tags": tags.Schema(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if err := validatePrivateLinkNatIpConfiguration(d); err != nil {
				return err
			}

			return nil
		}),
	}
}

func resourcePrivateLinkServiceCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)

	if d.IsNewResource() {
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
	autoApproval := d.Get("auto_approval_subscription_ids").(*pluginsdk.Set).List()
	enableProxyProtocol := d.Get("enable_proxy_protocol").(bool)
	primaryIpConfiguration := d.Get("nat_ip_configuration").([]interface{})
	loadBalancerFrontendIpConfigurations := d.Get("load_balancer_frontend_ip_configuration_ids").(*pluginsdk.Set).List()
	visibility := d.Get("visibility_subscription_ids").(*pluginsdk.Set).List()
	t := d.Get("tags").(map[string]interface{})

	parameters := network.PrivateLinkService{
		Location: utils.String(location),
		PrivateLinkServiceProperties: &network.PrivateLinkServiceProperties{
			AutoApproval: &network.PrivateLinkServicePropertiesAutoApproval{
				Subscriptions: utils.ExpandStringSlice(autoApproval),
			},
			EnableProxyProtocol: utils.Bool(enableProxyProtocol),
			Visibility: &network.PrivateLinkServicePropertiesVisibility{
				Subscriptions: utils.ExpandStringSlice(visibility),
			},
			IPConfigurations:                     expandPrivateLinkServiceIPConfiguration(primaryIpConfiguration),
			LoadBalancerFrontendIPConfigurations: expandPrivateLinkServiceFrontendIPConfiguration(loadBalancerFrontendIpConfigurations),
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

	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for Private Link Service to %q (Resource Group %q) to finish applying", name, resourceGroup)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, resourceGroup, name),
		MinTimeout: 15 * time.Second,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("Error waiting for Private Link Service %q (Resource Group %q) to become available: %s", name, resourceGroup, err)
	}

	// TODO: switch over to using an ID parser
	resp, err := client.Get(ctx, resourceGroup, name, "")
	if err != nil {
		return fmt.Errorf("Error retrieving Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}
	if resp.ID == nil || *resp.ID == "" {
		return fmt.Errorf("API returns a nil/empty id on Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
	}

	d.SetId(*resp.ID)

	return resourcePrivateLinkServiceRead(d, meta)
}

func resourcePrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
		d.Set("enable_proxy_protocol", props.EnableProxyProtocol)

		var autoApprovalSub []interface{}
		if autoApproval := props.AutoApproval; autoApproval != nil {
			autoApprovalSub = utils.FlattenStringSlice(autoApproval.Subscriptions)
		}
		if err := d.Set("auto_approval_subscription_ids", autoApprovalSub); err != nil {
			return fmt.Errorf("Error setting `auto_approval_subscription_ids`: %+v", err)
		}

		var subscriptions []interface{}
		if visibility := props.Visibility; visibility != nil {
			subscriptions = utils.FlattenStringSlice(visibility.Subscriptions)
		}
		if err := d.Set("visibility_subscription_ids", subscriptions); err != nil {
			return fmt.Errorf("Error setting `visibility_subscription_ids`: %+v", err)
		}

		if err := d.Set("nat_ip_configuration", flattenPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `nat_ip_configuration`: %+v", err)
		}

		if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
			return fmt.Errorf("Error setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

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
			return fmt.Errorf("Error waiting for deletion of Private Link Service %q (Resource Group %q): %+v", name, resourceGroup, err)
		}
	}

	return nil
}

func expandPrivateLinkServiceIPConfiguration(input []interface{}) *[]network.PrivateLinkServiceIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]network.PrivateLinkServiceIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		privateIpAddress := v["private_ip_address"].(string)
		subnetId := v["subnet_id"].(string)
		privateIpAddressVersion := v["private_ip_address_version"].(string)
		name := v["name"].(string)
		primary := v["primary"].(bool)

		result := network.PrivateLinkServiceIPConfiguration{
			Name: utils.String(name),
			PrivateLinkServiceIPConfigurationProperties: &network.PrivateLinkServiceIPConfigurationProperties{
				PrivateIPAddress:        utils.String(privateIpAddress),
				PrivateIPAddressVersion: network.IPVersion(privateIpAddressVersion),
				Subnet: &network.Subnet{
					ID: utils.String(subnetId),
				},
				Primary: utils.Bool(primary),
			},
		}

		if privateIpAddress != "" {
			result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethodStatic
		} else {
			result.PrivateLinkServiceIPConfigurationProperties.PrivateIPAllocationMethod = network.IPAllocationMethodDynamic
		}

		results = append(results, result)
	}

	return &results
}

func expandPrivateLinkServiceFrontendIPConfiguration(input []interface{}) *[]network.FrontendIPConfiguration {
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

func flattenPrivateLinkServiceIPConfiguration(input *[]network.PrivateLinkServiceIPConfiguration) []interface{} {
	results := make([]interface{}, 0)
	if input == nil {
		return results
	}

	for _, item := range *input {
		name := ""
		if item.Name != nil {
			name = *item.Name
		}

		privateIpAddress := ""
		privateIpVersion := ""
		subnetId := ""
		primary := false

		if props := item.PrivateLinkServiceIPConfigurationProperties; props != nil {
			if props.PrivateIPAddress != nil {
				privateIpAddress = *props.PrivateIPAddress
			}

			privateIpVersion = string(props.PrivateIPAddressVersion)

			if props.Subnet != nil && props.Subnet.ID != nil {
				subnetId = *props.Subnet.ID
			}

			if props.Primary != nil {
				primary = *props.Primary
			}
		}

		results = append(results, map[string]interface{}{
			"name":                       name,
			"primary":                    primary,
			"private_ip_address":         privateIpAddress,
			"private_ip_address_version": privateIpVersion,
			"subnet_id":                  subnetId,
		})
	}

	return results
}

func flattenPrivateLinkServiceFrontendIPConfiguration(input *[]network.FrontendIPConfiguration) *pluginsdk.Set {
	results := &pluginsdk.Set{F: pluginsdk.HashString}
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

func privateLinkServiceWaitForReadyRefreshFunc(ctx context.Context, client *network.PrivateLinkServicesClient, resourceGroupName string, name string) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, resourceGroupName, name, "")
		if err != nil {
			// the API is eventually consistent during recreates..
			if utils.ResponseWasNotFound(res.Response) {
				return res, "Pending", nil
			}

			return nil, "Error", fmt.Errorf("Error issuing read request in privateLinkServiceWaitForReadyRefreshFunc %q (Resource Group %q): %s", name, resourceGroupName, err)
		}
		if props := res.PrivateLinkServiceProperties; props != nil {
			if state := props.ProvisioningState; state != "" {
				return res, string(state), nil
			}
		}

		return res, "Pending", nil
	}
}
func validatePrivateLinkNatIpConfiguration(d *pluginsdk.ResourceDiff) error {
	name := d.Get("name").(string)
	resourceGroup := d.Get("resource_group_name").(string)
	ipConfigurations := d.Get("nat_ip_configuration").([]interface{})

	for i, item := range ipConfigurations {
		v := item.(map[string]interface{})
		p := fmt.Sprintf("nat_ip_configuration.%d.private_ip_address", i)
		s := fmt.Sprintf("nat_ip_configuration.%d.subnet_id", i)
		isPrimary := v["primary"].(bool)
		in := v["name"].(string)

		if d.HasChange(p) {
			o, n := d.GetChange(p)
			if o != "" && n == "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) nat_ip_configuration %q private_ip_address once assigned can not be removed", name, resourceGroup, in)
			}
		}

		if isPrimary && d.HasChange(s) {
			o, _ := d.GetChange(s)
			if o != "" {
				return fmt.Errorf("Private Link Service %q (Resource Group %q) nat_ip_configuration %q primary subnet_id once assigned can not be changed", name, resourceGroup, in)
			}
		}
	}

	return nil
}
