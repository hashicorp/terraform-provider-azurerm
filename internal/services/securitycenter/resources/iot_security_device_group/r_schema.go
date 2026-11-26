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
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.StringIsNotEmpty,
		},

		"iothub_id": {
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
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"connection_to_ips_not_allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type:         pluginsdk.TypeString,
							ValidateFunc: validation.IsCIDRIPv4,
						},
					},

					"local_users_not_allowed": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						Elem: &pluginsdk.Schema{
							Type: pluginsdk.TypeString,
						},
					},

					"processes_not_allowed": {
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
						Type:     pluginsdk.TypeString,
						Required: true,
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
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"min": {
						Type:         pluginsdk.TypeInt,
						Required:     true,
						ValidateFunc: validation.IntAtLeast(0),
					},

					"duration": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validation.ISO8601Duration,
					},
				},
			},
		},
	}
}
