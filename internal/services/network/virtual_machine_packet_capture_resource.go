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
	networkValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/network/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceVirtualMachinePacketCapture() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceVirtualMachinePacketCaptureCreate,
		Read:   resourceVirtualMachinePacketCaptureRead,
		Delete: resourceVirtualMachinePacketCaptureDelete,

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

			"virtual_machine_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.Any(
					commonids.ValidateVirtualMachineID,
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
		},
	}
}

func resourceVirtualMachinePacketCaptureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	watcherId, err := networkwatchers.ParseNetworkWatcherID(d.Get("network_watcher_id").(string))
	if err != nil {
		return err
	}

	id := packetcaptures.NewPacketCaptureID(subscriptionId, watcherId.ResourceGroupName, watcherId.NetworkWatcherName, d.Get("name").(string))

	targetResourceId := d.Get("virtual_machine_id").(string)
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
		return tf.ImportAsExistsError("azurerm_virtual_machine_packet_capture", id.ID())
	}

	storageLocation := expandVirtualMachinePacketCaptureStorageLocation(d.Get("storage_location").([]interface{}))
	payload := packetcaptures.PacketCapture{
		Properties: packetcaptures.PacketCaptureParameters{
			Target:                  targetResourceId,
			TargetType:              pointer.To(packetcaptures.PacketCaptureTargetTypeAzureVM),
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: utils.Int64(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      utils.Int64(int64(timeLimitInSeconds)),
			TotalBytesPerSession:    utils.Int64(int64(totalBytesPerSession)),
			Filters:                 expandVirtualMachinePacketCaptureFilters(d.Get("filter").([]interface{})),
		},
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())

	return resourceVirtualMachinePacketCaptureRead(d, meta)
}

func resourceVirtualMachinePacketCaptureRead(d *pluginsdk.ResourceData, meta interface{}) error {
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

		return fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	d.Set("name", id.PacketCaptureName)
	d.Set("network_watcher_id", networkwatchers.NewNetworkWatcherID(id.SubscriptionId, id.ResourceGroupName, id.NetworkWatcherName).ID())

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("virtual_machine_id", props.Target)
			d.Set("maximum_bytes_per_packet", int(*props.BytesToCapturePerPacket))
			d.Set("maximum_bytes_per_session", int(*props.TotalBytesPerSession))
			d.Set("maximum_capture_duration_in_seconds", int(*props.TimeLimitInSeconds))

			location := flattenVirtualMachinePacketCaptureStorageLocation(props.StorageLocation)
			if err := d.Set("storage_location", location); err != nil {
				return fmt.Errorf("setting `storage_location`: %+v", err)
			}

			filters := flattenVirtualMachinePacketCaptureFilters(props.Filters)
			if err := d.Set("filter", filters); err != nil {
				return fmt.Errorf("setting `filter`: %+v", err)
			}
		}
	}

	return nil
}

func resourceVirtualMachinePacketCaptureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandVirtualMachinePacketCaptureStorageLocation(input []interface{}) packetcaptures.PacketCaptureStorageLocation {
	location := input[0].(map[string]interface{})

	storageLocation := packetcaptures.PacketCaptureStorageLocation{}

	if v := location["file_path"]; v != "" {
		storageLocation.FilePath = pointer.To(v.(string))
	}
	if v := location["storage_account_id"]; v != "" {
		storageLocation.StorageId = pointer.To(v.(string))
	}

	return storageLocation
}

func flattenVirtualMachinePacketCaptureStorageLocation(input packetcaptures.PacketCaptureStorageLocation) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"file_path":          pointer.From(input.FilePath),
			"storage_account_id": pointer.From(input.StorageId),
			"storage_path":       pointer.From(input.StoragePath),
		},
	}
}

func expandVirtualMachinePacketCaptureFilters(input []interface{}) *[]packetcaptures.PacketCaptureFilter {
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

		filter := packetcaptures.PacketCaptureFilter{
			LocalIPAddress:  utils.String(localIPAddress),
			LocalPort:       utils.String(localPort),
			Protocol:        pointer.To(packetcaptures.PcProtocol(protocol)),
			RemoteIPAddress: utils.String(remoteIPAddress),
			RemotePort:      utils.String(remotePort),
		}
		filters = append(filters, filter)
	}

	return &filters
}

func flattenVirtualMachinePacketCaptureFilters(input *[]packetcaptures.PacketCaptureFilter) []interface{} {
	filters := make([]interface{}, 0)

	if input != nil {
		for _, v := range *input {
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
