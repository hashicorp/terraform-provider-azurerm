// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_device_group

import (
	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func iotSecurityDeviceGroupSchema() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Description:  "The name of the Device Security Group.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"iothub_id": {
			Description:  "The ID of the IoT Hub which to link the Security Device Group to.",
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.IotHubID,
		},

		"allow_rule": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"connection_from_ips_not_allowed": {
						Description:  "Which IP is not allowed to be connected to in current device group for inbound connection.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"connection_to_ips_not_allowed": {
						Description:  "Which IP is not allowed to be connected to in current device group for outbound connection.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"local_users_not_allowed": {
						Description:  "Which local user is not allowed to login in current device group.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"processes_not_allowed": {
						Description:  "Which process is not allowed to be executed in current device group.",
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},
				},
			},
		},

		"range_rule": {
			Type:     pluginsdk.TypeSet,
			Optional: true,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"type": {
						Description: "The type of supported rule type.",
						Type:        pluginsdk.TypeString,
						Required:    true,
						ValidateFunc: validation.StringInSlice([]string{
							string(security.RuleTypeActiveConnectionsNotInAllowedRange),
							string(security.RuleTypeAmqpC2DMessagesNotInAllowedRange),
							string(security.RuleTypeMqttC2DMessagesNotInAllowedRange),
							string(security.RuleTypeHTTPC2DMessagesNotInAllowedRange),
							string(security.RuleTypeAmqpC2DRejectedMessagesNotInAllowedRange),
							string(security.RuleTypeMqttC2DRejectedMessagesNotInAllowedRange),
							string(security.RuleTypeHTTPC2DRejectedMessagesNotInAllowedRange),
							string(security.RuleTypeAmqpD2CMessagesNotInAllowedRange),
							string(security.RuleTypeMqttD2CMessagesNotInAllowedRange),
							string(security.RuleTypeHTTPD2CMessagesNotInAllowedRange),
							string(security.RuleTypeDirectMethodInvokesNotInAllowedRange),
							string(security.RuleTypeFailedLocalLoginsNotInAllowedRange),
							string(security.RuleTypeFileUploadsNotInAllowedRange),
							string(security.RuleTypeQueuePurgesNotInAllowedRange),
							string(security.RuleTypeTwinUpdatesNotInAllowedRange),
							string(security.RuleTypeUnauthorizedOperationsNotInAllowedRange),
						}, false),
					},

					"max": {
						Description:  "The maximum threshold in the given time window.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"min": {
						Description:  "The minimum threshold in the given time window.",
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"duration": {
						Description:  "The time range. Represented in ISO 8601 duration format.",
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.ISO8601Duration,
					},
				},
			},
		},
	}
}
