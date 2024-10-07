// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonids"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/tags"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/networkprofiles"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/locks"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

const azureNetworkProfileResourceName = "azurerm_network_profile"

func resourceNetworkProfile() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceNetworkProfileCreate,
		Read:   resourceNetworkProfileRead,
		Update: resourceNetworkProfileUpdate,
		Delete: resourceNetworkProfileDelete,
		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := networkprofiles.ParseNetworkProfileID(id)
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
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"location": commonschema.Location(),

			"resource_group_name": commonschema.ResourceGroupName(),

			"container_network_interface": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"name": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ValidateFunc: validation.StringIsNotEmpty,
						},
						"ip_configuration": {
							Type:     pluginsdk.TypeList,
							Required: true,
							Elem: &pluginsdk.Resource{
								Schema: map[string]*pluginsdk.Schema{
									"name": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: validation.StringIsNotEmpty,
									},
									"subnet_id": {
										Type:         pluginsdk.TypeString,
										Required:     true,
										ValidateFunc: commonids.ValidateSubnetID,
									},
								},
							},
						},
					},
				},
			},

			"container_network_interface_ids": {
				Type:     pluginsdk.TypeList,
				Computed: true,
				Elem: &pluginsdk.Schema{
					Type: pluginsdk.TypeString,
				},
			},

			"tags": commonschema.Tags(),
		},
	}
}

func resourceNetworkProfileCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	defer cancel()

	id := networkprofiles.NewNetworkProfileID(subscriptionId, d.Get("resource_group_name").(string), d.Get("name").(string))

	existing, err := client.Get(ctx, id, networkprofiles.DefaultGetOperationOptions())
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_profile", id.ID())
	}

	containerNetworkInterfaceConfigurations := expandNetworkProfileContainerNetworkInterface(d.Get("container_network_interface").([]interface{}))
	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(containerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByName(id.NetworkProfileName, azureNetworkProfileResourceName)
	defer locks.UnlockByName(id.NetworkProfileName, azureNetworkProfileResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetsToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetsToLock, SubnetResourceName)

	payload := networkprofiles.NetworkProfile{
		Location: pointer.To(location.Normalize(d.Get("location").(string))),
		Tags:     tags.Expand(d.Get("tags").(map[string]interface{})),
		Properties: &networkprofiles.NetworkProfilePropertiesFormat{
			ContainerNetworkInterfaceConfigurations: containerNetworkInterfaceConfigurations,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkProfileRead(d, meta)
}

func resourceNetworkProfileUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
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

	containerNetworkInterfaceConfigurations := expandNetworkProfileContainerNetworkInterface(d.Get("container_network_interface").([]interface{}))
	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(containerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByName(id.NetworkProfileName, azureNetworkProfileResourceName)
	defer locks.UnlockByName(id.NetworkProfileName, azureNetworkProfileResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetsToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetsToLock, SubnetResourceName)

	if d.HasChange("container_network_interface") {
		payload.Properties.ContainerNetworkInterfaceConfigurations = containerNetworkInterfaceConfigurations
	}

	if d.HasChange("tags") {
		payload.Tags = tags.Expand(d.Get("tags").(map[string]interface{}))
	}

	if _, err := client.CreateOrUpdate(ctx, *id, *payload); err != nil {
		return fmt.Errorf("updating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceNetworkProfileRead(d, meta)
}

func resourceNetworkProfileRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[DEBUG] %s was not found - removing from state!", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("requesting %s: %+v", *id, err)
	}

	d.Set("name", id.NetworkProfileName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		d.Set("location", location.NormalizeNilable(model.Location))

		if props := model.Properties; props != nil {
			cniConfigs := flattenNetworkProfileContainerNetworkInterface(props.ContainerNetworkInterfaceConfigurations)
			if err := d.Set("container_network_interface", cniConfigs); err != nil {
				return fmt.Errorf("setting `container_network_interface`: %+v", err)
			}

			cniIDs := flattenNetworkProfileContainerNetworkInterfaceIDs(props.ContainerNetworkInterfaces)
			if err := d.Set("container_network_interface_ids", cniIDs); err != nil {
				return fmt.Errorf("setting `container_network_interface_ids`: %+v", err)
			}
		}

		if err := tags.FlattenAndSet(d, model.Tags); err != nil {
			return err
		}
	}

	return nil
}

func resourceNetworkProfileDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.NetworkProfiles
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := networkprofiles.ParseNetworkProfileID(d.Id())
	if err != nil {
		return err
	}

	existing, err := client.Get(ctx, *id, networkprofiles.DefaultGetOperationOptions())
	if err != nil {
		return fmt.Errorf("retrieving existing %s: %+v", *id, err)
	}
	if existing.Model == nil || existing.Model.Properties == nil {
		return fmt.Errorf("retrieving existing %s: `model.Properties` was nil", *id)
	}

	subnetsToLock, vnetsToLock, err := expandNetworkProfileVirtualNetworkSubnetNames(existing.Model.Properties.ContainerNetworkInterfaceConfigurations)
	if err != nil {
		return fmt.Errorf("extracting names of Subnet and Virtual Network: %+v", err)
	}

	locks.ByName(id.NetworkProfileName, azureNetworkProfileResourceName)
	defer locks.UnlockByName(id.NetworkProfileName, azureNetworkProfileResourceName)

	locks.MultipleByName(vnetsToLock, VirtualNetworkResourceName)
	defer locks.UnlockMultipleByName(vnetsToLock, VirtualNetworkResourceName)

	locks.MultipleByName(subnetsToLock, SubnetResourceName)
	defer locks.UnlockMultipleByName(subnetsToLock, SubnetResourceName)

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return err
}

func expandNetworkProfileContainerNetworkInterface(input []interface{}) *[]networkprofiles.ContainerNetworkInterfaceConfiguration {
	retCNIConfigs := make([]networkprofiles.ContainerNetworkInterfaceConfiguration, 0)

	for _, cniConfig := range input {
		nciData := cniConfig.(map[string]interface{})
		nciName := nciData["name"].(string)
		ipConfigs := nciData["ip_configuration"].([]interface{})

		retIPConfigs := make([]networkprofiles.IPConfigurationProfile, 0)
		for _, ipConfig := range ipConfigs {
			ipData := ipConfig.(map[string]interface{})
			ipName := ipData["name"].(string)
			subnetId := ipData["subnet_id"].(string)

			retIPConfig := networkprofiles.IPConfigurationProfile{
				Name: &ipName,
				Properties: &networkprofiles.IPConfigurationProfilePropertiesFormat{
					Subnet: &networkprofiles.Subnet{
						Id: &subnetId,
					},
				},
			}

			retIPConfigs = append(retIPConfigs, retIPConfig)
		}

		retCNIConfig := networkprofiles.ContainerNetworkInterfaceConfiguration{
			Name: &nciName,
			Properties: &networkprofiles.ContainerNetworkInterfaceConfigurationPropertiesFormat{
				IPConfigurations: &retIPConfigs,
			},
		}

		retCNIConfigs = append(retCNIConfigs, retCNIConfig)
	}

	return &retCNIConfigs
}

func expandNetworkProfileVirtualNetworkSubnetNames(input *[]networkprofiles.ContainerNetworkInterfaceConfiguration) (*[]string, *[]string, error) {
	subnetNames := make([]string, 0)
	vnetNames := make([]string, 0)

	if input != nil {
		for _, item := range *input {
			if item.Properties == nil || item.Properties.IPConfigurations == nil {
				continue
			}

			for _, config := range *item.Properties.IPConfigurations {
				if config.Properties == nil || config.Properties.Subnet == nil || config.Properties.Subnet.Id == nil {
					continue
				}

				subnetId, err := commonids.ParseSubnetIDInsensitively(*config.Properties.Subnet.Id)
				if err != nil {
					return nil, nil, err
				}

				if !utils.SliceContainsValue(subnetNames, subnetId.SubnetName) {
					subnetNames = append(subnetNames, subnetId.SubnetName)
				}

				if !utils.SliceContainsValue(vnetNames, subnetId.VirtualNetworkName) {
					vnetNames = append(vnetNames, subnetId.VirtualNetworkName)
				}
			}
		}
	}

	return &subnetNames, &vnetNames, nil
}

func flattenNetworkProfileContainerNetworkInterface(input *[]networkprofiles.ContainerNetworkInterfaceConfiguration) []interface{} {
	output := make([]interface{}, 0)
	if input == nil {
		return output
	}

	for _, cniConfig := range *input {
		ipConfigurations := make([]interface{}, 0)
		if props := cniConfig.Properties; props != nil && props.IPConfigurations != nil {
			for _, ipConfig := range *props.IPConfigurations {
				subnetId := ""
				if ipProps := ipConfig.Properties; ipProps != nil && ipProps.Subnet != nil && ipProps.Subnet.Id != nil {
					subnetId = *ipProps.Subnet.Id
				}

				ipConfigurations = append(ipConfigurations, map[string]interface{}{
					"name":      pointer.From(ipConfig.Name),
					"subnet_id": subnetId,
				})
			}
		}

		output = append(output, map[string]interface{}{
			"name":             pointer.From(cniConfig.Name),
			"ip_configuration": ipConfigurations,
		})
	}

	return output
}

func flattenNetworkProfileContainerNetworkInterfaceIDs(input *[]networkprofiles.ContainerNetworkInterface) []string {
	retCNIs := make([]string, 0)

	if input != nil {
		for _, retCNI := range *input {
			if retCNI.Id != nil {
				retCNIs = append(retCNIs, *retCNI.Id)
			}
		}
	}

	return retCNIs
}
