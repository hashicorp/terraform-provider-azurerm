// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package virtual_machine_scale_set_packet_capture

import (
	"time"

	"github.com/hashicorp/go-azure-sdk/resource-manager/network/2025-07-01/packetcaptures"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func RegisterResource() *pluginsdk.Resource {
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

		Schema: virtualMachineScaleSetPacketCaptureSchema(),
	}
}
