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
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-09-01/networkwatchers"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-11-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	computeParse "github.com/hashicorp/terraform-provider-azurerm/internal/services/compute/parse"
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualMachineScaleSetPacketCapture() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachineScaleSetPacketCaptureCreate,
		Read:   resourceVirtualMachineScaleSetPacketCaptureRead,
		Delete: resourceVirtualMachineScaleSetPacketCaptureDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := packetcaptures.ParsePacketCaptureID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"network_watcher_id": {
				Type:         pluginsdk.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: networkwatchers.ValidateNetworkWatcherID,
			},

			"virtual_machine_scale_set_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					commonids.ValidateVirtualMachineScaleSetID,
				),
			},

			"maximum_bytes_per_packet": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  0,
			},

			"maximum_bytes_per_session": {
				Type:     pluginsdk.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  1073741824,
			},

			"maximum_capture_duration_in_seconds": {
				Type:         pluginsdk.TypeInt,
				Optional:     true,
				ForceNew:     true,
				Default:      18000,
				ValidateFunc: validation.IntBetween(1, 18000),
			},

			"storage_location": {
				Type:     pluginsdk.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"file_path": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: networkValidate.FilePath,
							AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
						},
						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
							ValidateFunc: commonids.ValidateStorageAccountID,
							AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
						},
						"storage_path": {
							Type:     pluginsdk.TypeString,
							Computed: true,
						},
					},
				},
			},

			"filter": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"local_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"local_port": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"protocol": {
							Type:         pluginsdk.TypeString,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.StringInSlice(packetcaptures.PossibleValuesForPcProtocol(), false),
						},
						"remote_ip_address": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"remote_port": {
							Type:     pluginsdk.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},

			"machine_scope": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"exclude_instance_ids": {
							Type:          pluginsdk.TypeList,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"machine_scope.0.include_instance_ids"},
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},

						"include_instance_ids": {
							Type:          pluginsdk.TypeList,
							Optional:      true,
							ForceNew:      true,
							ConflictsWith: []string{"machine_scope.0.exclude_instance_ids"},
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validation.StringIsNotEmpty,
							},
						},
					},
				},
			},
		},
	}
}

func resourceVirtualMachineScaleSetPacketCaptureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := networkwatchers.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	id := packetcaptures.NewPacketCaptureID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	targetResourceId := d.Get("virtual_machine_scale_set_id").(string)
	bytesToCapturePerPacket := d.Get("maximum_bytes_per_packet").(int)
	totalBytesPerSession := d.Get("maximum_bytes_per_session").(int)
	timeLimitInSeconds := d.Get("maximum_capture_duration_in_seconds").(int)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_virtual_machine_scale_set_packet_capture", id.ID())
	}

	storageLocation := expandVirtualMachineScaleSetPacketCaptureStorageLocation(d.Get("storage_location").([]interface{}))
	payload := packetcaptures.PacketCapture{
		Properties: packetcaptures.PacketCaptureParameters{
			Target:                  targetResourceId,
			TargetType:              pointer.To(packetcaptures.PacketCaptureTargetTypeAzureVMSS),
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: pointer.To(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      pointer.To(int64(timeLimitInSeconds)),
			TotalBytesPerSession:    pointer.To(int64(totalBytesPerSession)),
			Filters:                 expandVirtualMachineScaleSetPacketCaptureFilters(d.Get("filter").([]interface{})),
		},
	}

	if v, ok := d.GetOk("machine_scope"); ok {
		payload.Properties.Scope = expandVirtualMachineScaleSetPacketCaptureMachineScope(v.([]interface{}))
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachineScaleSetPacketCaptureRead(d, meta)
}

func resourceVirtualMachineScaleSetPacketCaptureRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := packetcaptures.ParsePacketCaptureID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			log.Printf("[WARN] %s not found - removing from state", *id)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("reading %s: %+v", *id, err)
	}

	d.Set("name", id.PacketCaptureName)
	d.Set("network_watcher_id", networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("virtual_machine_scale_set_id", props.Target)
			d.Set("maximum_bytes_per_packet", int(*props.BytesToCapturePerPacket))
			d.Set("maximum_bytes_per_session", int(*props.TotalBytesPerSession))
			d.Set("maximum_capture_duration_in_seconds", int(*props.TimeLimitInSeconds))

			location := flattenVirtualMachineScaleSetPacketCaptureStorageLocation(props.StorageLocation)
			if err := d.Set("storage_location", location); err != nil {
				return fmt.Errorf("setting `storage_location`: %+v", err)
			}

			filters := flattenVirtualMachineScaleSetPacketCaptureFilters(props.Filters)
			if err := d.Set("filter", filters); err != nil {
				return fmt.Errorf("setting `filter`: %+v", err)
			}

			scope, err := flattenVirtualMachineScaleSetPacketCaptureMachineScope(props.Scope)
			if err != nil {
				return err
			}
			if err := d.Set("machine_scope", scope); err != nil {
				return fmt.Errorf(`setting "machine_scope": %+v`, err)
			}
		}
	}

	return nil
}

func resourceVirtualMachineScaleSetPacketCaptureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := packetcaptures.ParsePacketCaptureID(d.Id())
	if err != nil {
		return err
	}

	if err := client.DeleteThenPoll(ctx, *id); err != nil {
		return fmt.Errorf("deleting %s: %+v", *id, err)
	}

	return nil
}

func expandVirtualMachineScaleSetPacketCaptureStorageLocation(input []interface{}) packetcaptures.PacketCaptureStorageLocation {
	location := input[0].(map[string]interface{})

	storageLocation := packetcaptures.PacketCaptureStorageLocation{}

	if v := location["file_path"]; v != "" {
		storageLocation.FilePath = utils.String(v.(string))
	}
	if v := location["storage_account_id"]; v != "" {
		storageLocation.StorageId = utils.String(v.(string))
	}

	return storageLocation
}

func flattenVirtualMachineScaleSetPacketCaptureStorageLocation(input packetcaptures.PacketCaptureStorageLocation) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"file_path":          pointer.From(input.FilePath),
			"storage_account_id": pointer.From(input.StorageId),
			"storage_path":       pointer.From(input.StoragePath),
		},
	}
}

func expandVirtualMachineScaleSetPacketCaptureFilters(input []interface{}) *[]packetcaptures.PacketCaptureFilter {
	if len(input) == 0 {
		return nil
	}

	filters := make([]packetcaptures.PacketCaptureFilter, 0)
	for _, v := range input {
		inputFilter := v.(map[string]interface{})

		localIPAddress := inputFilter["local_ip_address"].(string)
		localPort := inputFilter["local_port"].(string) // TODO: should this be an int?
		protocol := inputFilter["protocol"].(string)
		remoteIPAddress := inputFilter["remote_ip_address"].(string)
		remotePort := inputFilter["remote_port"].(string)

		filters = append(filters, packetcaptures.PacketCaptureFilter{
			LocalIPAddress:  pointer.To(localIPAddress),
			LocalPort:       pointer.To(localPort),
			Protocol:        pointer.To(packetcaptures.PcProtocol(protocol)),
			RemoteIPAddress: pointer.To(remoteIPAddress),
			RemotePort:      pointer.To(remotePort),
		})
	}

	return &filters
}

func flattenVirtualMachineScaleSetPacketCaptureFilters(input *[]packetcaptures.PacketCaptureFilter) []interface{} {
	filters := make([]interface{}, 0)

	if inFilter := input; inFilter != nil {
		for _, v := range *inFilter {
			protocol := ""
			if v.Protocol != nil {
				protocol = string(*v.Protocol)
			}

			filters = append(filters, map[string]interface{}{
				"local_ip_address":  pointer.From(v.LocalIPAddress),
				"local_port":        pointer.From(v.LocalPort),
				"protocol":          protocol,
				"remote_ip_address": pointer.From(v.RemoteIPAddress),
				"remote_port":       pointer.From(v.RemotePort),
			})
		}
	}

	return filters
}

func expandVirtualMachineScaleSetPacketCaptureMachineScope(input []interface{}) *packetcaptures.PacketCaptureMachineScope {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	raw := input[0].(map[string]interface{})
	output := &packetcaptures.PacketCaptureMachineScope{}

	if exclude := raw["exclude_instance_ids"].([]interface{}); len(exclude) > 0 {
		output.Exclude = utils.ExpandStringSlice(exclude)
	}

	if include := raw["include_instance_ids"].([]interface{}); len(include) > 0 {
		output.Include = utils.ExpandStringSlice(include)
	}

	return output
}

func flattenVirtualMachineScaleSetPacketCaptureMachineScope(input *packetcaptures.PacketCaptureMachineScope) ([]interface{}, error) {
	outputs := make([]interface{}, 0)
	if input == nil || (input.Exclude == nil && input.Include == nil) || (len(*input.Exclude) == 0 && len(*input.Include) == 0) {
		return outputs, nil
	}

	excludedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Exclude)
	if err != nil {
		return nil, err
	}

	includedInstanceIds, err := flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input.Include)
	if err != nil {
		return nil, err
	}

	outputs = append(outputs, map[string]interface{}{
		"exclude_instance_ids": excludedInstanceIds,
		"include_instance_ids": includedInstanceIds,
	})

	return outputs, nil
}

func flattenVirtualMachineScaleSetPacketCaptureScopeInstanceIds(input *[]string) ([]string, error) {
	instances := make([]string, 0)
	if input == nil {
		return instances, nil
	}

	for _, instance := range *input {
		vmssInstanceId, err := computeParse.VMSSInstanceIDInsensitively(instance)
		if err != nil {
			return nil, err
		}

		instances = append(instances, vmssInstanceId.VirtualMachineName)
	}

	return instances, nil
}
