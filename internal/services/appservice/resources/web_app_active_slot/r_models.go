// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package web_app_active_slot

type WebAppActiveSlotModel struct {
	SlotID              string `tfschema:"slot_id"`
	OverwriteNetworking bool   `tfschema:"overwrite_network_config"` // Note: This setting controls the ambiguously named `PreserveVnet`
	LastSwap            string `tfschema:"last_successful_swap"`
}

func (r Resource) ModelObject() any {
	return &WebAppActiveSlotModel{}
}
