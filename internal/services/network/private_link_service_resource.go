// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/privatelinkservices"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourcePrivateLinkService() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourcePrivateLinkServiceCreate,
		Read:   resourcePrivateLinkServiceRead,
		Update: resourcePrivateLinkServiceUpdate,
		Delete: resourcePrivateLinkServiceDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := privatelinkservices.ParsePrivateLinkServiceID(id)
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

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

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
					ValidateFunc: validation.Any(validation.IsUUID, validation.StringInSlice([]string{"*"}, false)),
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
								string(privatelinkservices.IPVersionIPvFour),
							}, false),
							Default: string(privatelinkservices.IPVersionIPvFour),
						},
						"subnet_id": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: commonids.ValidateSubnetID,
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

			"tags": commonschema.Tags(),
		},

		CustomizeDiff: pluginsdk.CustomizeDiffShim(func(ctx context.Context, d *pluginsdk.ResourceDiff, v interface{}) error {
			if err := validatePrivateLinkNatIpConfiguration(d); err != nil {
				return err
			}

			return nil
		}),
	}
}

func resourcePrivateLinkServiceCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := privatelinkservices.NewPrivateLinkServiceID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, privatelinkservices.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}
	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_private_link_service", id.ID())
	}

	parameters := privatelinkservices.PrivateLinkService{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Properties: &privatelinkservices.PrivateLinkServiceProperties{
			AutoApproval: &privatelinkservices.ResourceSet{
				Subscriptions: utils.ExpandStringSlice(d.Get("auto_approval_subscription_ids").(*pluginsdk.Set).List()),
			},
			EnableProxyProtocol: pointer.To(d.Get("enable_proxy_protocol").(bool)),
			Visibility: &privatelinkservices.ResourceSet{
				Subscriptions: utils.ExpandStringSlice(d.Get("visibility_subscription_ids").(*pluginsdk.Set).List()),
			},
			IPConfigurations:                     expandPrivateLinkServiceIPConfiguration(d.Get("nat_ip_configuration").([]interface{})),
			LoadBalancerFrontendIPConfigurations: expandPrivateLinkServiceFrontendIPConfiguration(d.Get("load_balancer_frontend_ip_configuration_ids").(*pluginsdk.Set).List()),
			Fqdns:                                utils.ExpandStringSlice(d.Get("fqdns").([]interface{})),
		},
		Tags: tags.Expand(d.Get("tags").(map[string]interface{})),
	}

	if err := client.CreateOrUpdateThenPoll(ctx, id, parameters); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for %s to finish applying", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, id),
		MinTimeout: 15 * time.Second,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutCreate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	d.SetId(id.ID())

	return resourcePrivateLinkServiceRead(d, meta)
}

func resourcePrivateLinkServiceUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, privatelinkservices.DefaultGetOperationOptions())
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

	if d.HasChange("auto_approval_subscription_ids") {
		payload.Properties.AutoApproval = &privatelinkservices.ResourceSet{
			Subscriptions: utils.ExpandStringSlice(d.Get("auto_approval_subscription_ids").(*pluginsdk.Set).List()),
		}
	}

	if d.HasChange("enable_proxy_protocol") {
		payload.Properties.EnableProxyProtocol = pointer.To(d.Get("enable_proxy_protocol").(bool))
	}

	if d.HasChange("visibility_subscription_ids") {
		payload.Properties.Visibility = &privatelinkservices.ResourceSet{
			Subscriptions: utils.ExpandStringSlice(d.Get("visibility_subscription_ids").(*pluginsdk.Set).List()),
		}
	}

	if d.HasChange("fqdns") {
		payload.Properties.Fqdns = utils.ExpandStringSlice(d.Get("fqdns").([]interface{}))
	}

	if d.HasChange("nat_ip_configuration") {
		payload.Properties.IPConfigurations = expandPrivateLinkServiceIPConfiguration(d.Get("nat_ip_configuration").([]interface{}))
	}

	if d.HasChange("load_balancer_frontend_ip_configuration_ids") {
		payload.Properties.LoadBalancerFrontendIPConfigurations = expandPrivateLinkServiceFrontendIPConfiguration(d.Get("load_balancer_frontend_ip_configuration_ids").(*pluginsdk.Set).List())
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if err := client.CreateOrUpdateThenPoll(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	// we can't rely on the use of the Future here due to the resource being successfully completed but now the service is applying those values.
	// currently being tracked with issue #6466: https://github.com/Azure/azure-sdk-for-go/issues/6466
	log.Printf("[DEBUG] Waiting for %s to finish applying", id)
	stateConf := &pluginsdk.StateChangeConf{
		Pending:    []string{"Pending", "Updating", "Creating"},
		Target:     []string{"Succeeded"},
		Refresh:    privateLinkServiceWaitForReadyRefreshFunc(ctx, client, *id),
		MinTimeout: 15 * time.Second,
	}

	stateConf.Timeout = d.Timeout(pluginsdk.TimeoutUpdate)

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return fmt.Errorf("waiting for %s to become available: %s", id, err)
	}

	d.SetId(id.ID())

	return resourcePrivateLinkServiceRead(d, meta)
}

func resourcePrivateLinkServiceRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, privatelinkservices.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[INFO] %s does not exist - removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PrivateLinkServiceName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))
		if props := model.Properties; props != nil {
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
		return tags.FlattenAndSet(d, model.Tags)
	}
	return nil
}

func resourcePrivateLinkServiceDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PrivateLinkServices
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := privatelinkservices.ParsePrivateLinkServiceID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandPrivateLinkServiceIPConfiguration(input []interface{}) *[]privatelinkservices.PrivateLinkServiceIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]privatelinkservices.PrivateLinkServiceIPConfiguration, 0)

	for _, item := range input {
		v := item.(map[string]interface{})
		privateIpAddress := v["private_ip_address"].(string)
		subnetId := v["subnet_id"].(string)
		privateIpAddressVersion := v["private_ip_address_version"].(string)
		name := v["name"].(string)
		primary := v["primary"].(bool)

		result := privatelinkservices.PrivateLinkServiceIPConfiguration{
			Name: pointer.To(name),
			Properties: &privatelinkservices.PrivateLinkServiceIPConfigurationProperties{
				PrivateIPAddress:        pointer.To(privateIpAddress),
				PrivateIPAddressVersion: pointer.To(privatelinkservices.IPVersion(privateIpAddressVersion)),
				Subnet: &privatelinkservices.Subnet{
					Id: pointer.To(subnetId),
				},
				Primary: pointer.To(primary),
			},
		}

		if privateIpAddress != "" {
			result.Properties.PrivateIPAllocationMethod = pointer.To(privatelinkservices.IPAllocationMethodStatic)
		} else {
			result.Properties.PrivateIPAllocationMethod = pointer.To(privatelinkservices.IPAllocationMethodDynamic)
		}

		results = append(results, result)
	}

	return &results
}

func expandPrivateLinkServiceFrontendIPConfiguration(input []interface{}) *[]privatelinkservices.FrontendIPConfiguration {
	if len(input) == 0 {
		return nil
	}

	results := make([]privatelinkservices.FrontendIPConfiguration, 0)

	for _, item := range input {
		result := privatelinkservices.FrontendIPConfiguration{
			Id: pointer.To(item.(string)),
		}

		results = append(results, result)
	}

	return &results
}

func flattenPrivateLinkServiceIPConfiguration(input *[]privatelinkservices.PrivateLinkServiceIPConfiguration) []interface{} {
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

		if props := item.Properties; props != nil {
			if props.PrivateIPAddress != nil {
				privateIpAddress = *props.PrivateIPAddress
			}

			privateIpVersion = string(pointer.From(props.PrivateIPAddressVersion))

			if props.Subnet != nil && props.Subnet.Id != nil {
				subnetId = *props.Subnet.Id
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

func flattenPrivateLinkServiceFrontendIPConfiguration(input *[]privatelinkservices.FrontendIPConfiguration) *pluginsdk.Set {
	results := &pluginsdk.Set{F: pluginsdk.HashString}
	if input == nil {
		return results
	}

	for _, item := range *input {
		if id := item.Id; id != nil {
			results.Add(*id)
		}
	}

	return results
}

func privateLinkServiceWaitForReadyRefreshFunc(ctx context.Context, client *privatelinkservices.PrivateLinkServicesClient, id privatelinkservices.PrivateLinkServiceId) pluginsdk.StateRefreshFunc {
	return func() (interface{}, string, error) {
		res, err := client.Get(ctx, id, privatelinkservices.DefaultGetOperationOptions())
		if err != nil {
			// the API is eventually consistent during recreates..
			if response.WasNotFound(res.HttpResponse) {
				return res, "Pending", nil
			}

			return nil, "Error", fmt.Errorf("retrieving %s: %+v", id, err)
		}

		if model := res.Model; model != nil {
			if props := model.Properties; props != nil && props.ProvisioningState != nil {
				if state := *props.ProvisioningState; state != "" {
					return res, string(state), nil
				}
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
