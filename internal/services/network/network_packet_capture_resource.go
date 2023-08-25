// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package network

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/commonschema"
	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2023-04-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/network/migration"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceNetworkPacketCapture() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create:             resourceNetworkPacketCaptureCreate,
		Read:               resourceNetworkPacketCaptureRead,
		Delete:             resourceNetworkPacketCaptureDelete,
		DeprecationMessage: "The \"azurerm_network_packet_capture\" resource is deprecated and will be removed in favour of the `azurerm_virtual_machine_packet_capture` and `azurerm_virtual_machine_scale_set_packet_capture` resources in version 4.0 of the AzureRM Provider.",

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := packetcaptures.ParsePacketCaptureID(id)
			return err
		}),

		SchemaVersion: 1,
		StateUpgraders: pluginsdk.StateUpgrades(map[int]pluginsdk.StateUpgrade{
			0: migration.NetworkPacketCaptureV0ToV1{},
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"resource_group_name": commonschema.ResourceGroupName(),

			"network_watcher_name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
			},

			"target_resource_id": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ForceNew: true,
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

			"maximum_capture_duration": {
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
							AtLeastOneOf: []string{"storage_location.0.file_path", "storage_location.0.storage_account_id"},
						},
						"storage_account_id": {
							Type:         pluginsdk.TypeString,
							Optional:     true,
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

func resourceNetworkPacketCaptureCreate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).Network.PacketCaptures
	subscriptionId := meta.(*clients.Client).Account.SubscriptionId
	ctx, cancel := timeouts.ForCreate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := packetcaptures.NewPacketCaptureID(subscriptionId, d.Get("resource_group_name").(string), d.Get("network_watcher_name").(string), d.Get("name").(string))

	targetResourceId := d.Get("target_resource_id").(string)
	bytesToCapturePerPacket := d.Get("maximum_bytes_per_packet").(int)
	totalBytesPerSession := d.Get("maximum_bytes_per_session").(int)
	timeLimitInSeconds := d.Get("maximum_capture_duration").(int)

	existing, err := client.Get(ctx, id)
	if err != nil {
		if !response.WasNotFound(existing.HttpResponse) {
			return fmt.Errorf("checking for presence of existing %s: %s", id, err)
		}
	}

	if !response.WasNotFound(existing.HttpResponse) {
		return tf.ImportAsExistsError("azurerm_network_packet_capture", id.ID())
	}

	storageLocation := expandNetworkPacketCaptureStorageLocation(d.Get("storage_location").([]interface{}))
	payload := packetcaptures.PacketCapture{
		Properties: packetcaptures.PacketCaptureParameters{
			Target:                  targetResourceId,
			StorageLocation:         storageLocation,
			BytesToCapturePerPacket: pointer.To(int64(bytesToCapturePerPacket)),
			TimeLimitInSeconds:      pointer.To(int64(timeLimitInSeconds)),
			TotalBytesPerSession:    pointer.To(int64(totalBytesPerSession)),
			Filters:                 expandNetworkPacketCaptureFilters(d.Get("filter").([]interface{})),
		},
	}

	if err := client.CreateThenPoll(ctx, id, payload); err != nil {
		return fmt.Errorf("creating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceNetworkPacketCaptureRead(d, meta)
}

func resourceNetworkPacketCaptureRead(d *pluginsdk.ResourceData, meta interface{}) error {
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
	d.Set("network_watcher_name", id.NetworkWatcherName)
	d.Set("resource_group_name", id.ResourceGroupName)

	if model := resp.Model; model != nil {
		if props := model.Properties; props != nil {
			d.Set("target_resource_id", props.Target)
			d.Set("maximum_bytes_per_packet", int(pointer.From(props.BytesToCapturePerPacket)))
			d.Set("maximum_bytes_per_session", int(pointer.From(props.TotalBytesPerSession)))
			d.Set("maximum_capture_duration", int(pointer.From(props.TimeLimitInSeconds)))

			location := flattenNetworkPacketCaptureStorageLocation(props.StorageLocation)
			if err := d.Set("storage_location", location); err != nil {
				return fmt.Errorf("setting `storage_location`: %+v", err)
			}

			filters := flattenNetworkPacketCaptureFilters(props.Filters)
			if err := d.Set("filter", filters); err != nil {
				return fmt.Errorf("setting `filter`: %+v", err)
			}
		}
	}

	return nil
}

func resourceNetworkPacketCaptureDelete(d *pluginsdk.ResourceData, meta interface{}) error {
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

func expandNetworkPacketCaptureStorageLocation(input []interface{}) packetcaptures.PacketCaptureStorageLocation {
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

func flattenNetworkPacketCaptureStorageLocation(input packetcaptures.PacketCaptureStorageLocation) []interface{} {
	return []interface{}{
		map[string]interface{}{
			"file_path":          pointer.From(input.FilePath),
			"storage_account_id": pointer.From(input.StorageId),
			"storage_path":       pointer.From(input.StoragePath),
		},
	}
}

func expandNetworkPacketCaptureFilters(input []interface{}) *[]packetcaptures.PacketCaptureFilter {
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
			LocalIPAddress:  utils.String(localIPAddress),
			LocalPort:       utils.String(localPort),
			Protocol:        pointer.To(packetcaptures.PcProtocol(protocol)),
			RemoteIPAddress: utils.String(remoteIPAddress),
			RemotePort:      utils.String(remotePort),
		})
	}

	return &filters
}

func flattenNetworkPacketCaptureFilters(input *[]packetcaptures.PacketCaptureFilter) []interface{} {
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
