package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/network/mgmt/2021-08-01/network"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tags"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateLinkService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateLinkServiceCreateUpdate,
		Read:   resourcePrivateLinkServiceRead,
		Update: resourcePrivateLinkServiceCreateUpdate,
		Delete: resourcePrivateLinkServiceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.PrivateLinkServiceID(id)
			return err
		}),

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

			// TODO 4.0: change this from enable_* to *_enabled
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

			"fqdns": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
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
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewPrivateLinkServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	if d.IsNewResource() {
		existing, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
		if err != nil {
			if !utils.ResponseWasNotFound(existing.Response) {
				return fmt.Errorf("checking for presence of existing %s: %s", id, err)
			}
		}
		if !utils.ResponseWasNotFound(existing.Response) {
			return tf.ImportAsExistsError("azurerm_private_link_service", id.ID())
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
			Fqdns:                                utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
		},
		Tags: tags.Expand(t),
	}

	future, err := client.CreateOrUpdate(ctx, id.ResourceGroup, id.Name, parameters)
	if err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}
	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for creation of %s: %+v", id, err)
	}

	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for %s to finish applying", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, id.ResourceGroup, id.Name),
		MinTimeout: 15 * time.Second,
	}

	if d.IsNewResource() {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)
	} else {
		stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	d.SetId(id.ID())

	return resourcePrivateLinkServiceRead(d, meta)
}

func resourcePrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.ResourceGroup, id.Name, "")
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("[INFO] Private Link Service %q does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.Name)
	d.Set("resource_group_name", id.ResourceGroup)
	d.Set("location", azure.NormalizeLocation(*resp.Location))

	if props := resp.PrivateLinkServiceProperties; props != nil {
		d.Set("alias", props.Alias)
		d.Set("enable_proxy_protocol", props.EnableProxyProtocol)

		var autoApprovalSub []interface{}
		if autoApproval := props.AutoApproval; autoApproval != nil {
			autoApprovalSub = utils.FlattenStringSlice(autoApproval.Subscriptions)
		}
		if err := d.Set("auto_approval_subscription_ids", autoApprovalSub); err != nil {
			return fmt.Errorf("setting `auto_approval_subscription_ids`: %+v", err)
		}

		var subscriptions []interface{}
		if visibility := props.Visibility; visibility != nil {
			subscriptions = utils.FlattenStringSlice(visibility.Subscriptions)
		}
		if err := d.Set("visibility_subscription_ids", subscriptions); err != nil {
			return fmt.Errorf("setting `visibility_subscription_ids`: %+v", err)
		}

		if err := d.Set("fqdns", utils.FlattenStringSlice(props.Fqdns)); err != nil {
			return fmt.Errorf("setting `fqdns`: %+v", err)
		}

		if err := d.Set("nat_ip_configuration", flattenPrivateLinkServiceIPConfiguration(props.IPConfigurations)); err != nil {
			return fmt.Errorf("setting `nat_ip_configuration`: %+v", err)
		}

		if err := d.Set("load_balancer_frontend_ip_configuration_ids", flattenPrivateLinkServiceFrontendIPConfiguration(props.LoadBalancerFrontendIPConfigurations)); err != nil {
			return fmt.Errorf("setting `load_balancer_frontend_ip_configuration_ids`: %+v", err)
		}
	}

	return tags.FlattenAndSet(d, resp.Tags)
}

func resourcePrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServiceClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.PrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	future, err := client.Delete(ctx, id.ResourceGroup, id.Name)
	if err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	if err = future.WaitForCompletionRef(ctx, client.Client); err != nil {
		return fmt.Errorf("waiting for the deletion of %s: %+v", *id, err)
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

			return nil, "Error", fmt.Errorf("issuing read request in privateLinkServiceWaitForReadyRefreshFunc %q (Resource Group %q): %s", name, resourceGroupName, err)
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
