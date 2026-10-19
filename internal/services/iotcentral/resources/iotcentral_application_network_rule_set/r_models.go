// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iotcentral_application_network_rule_set

import (
	"github.com/hashicorp/go-azure-sdk/resource-manager/iotcentral/2021-11-01-preview/apps"
)

type IotCentralApplicationNetworkRuleSetModel struct {
	IotCentralApplicationId string             `tfschema:"iotcentral_application_id"`
	ApplyToDevice           bool               `tfschema:"apply_to_device"`
	DefaultAction           apps.NetworkAction `tfschema:"default_action"`
	IPRule                  []IPRule           `tfschema:"ip_rule"`
}

func (r Resource) ModelObject() any {
	return &IotCentralApplicationNetworkRuleSetModel{}
}
