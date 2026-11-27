// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_device_group

import (
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotSecurityDeviceGroupRead(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForRead(account.StopContext(meta), d)
	defer cancel()

	id, err := parse.IotSecurityDeviceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.IotHubID, id.Name)
	if err != nil {
		if response.WasNotFound(resp.Response.Response) {
			log.Printf("Device Security Group not found for %q: %+v", id.ID(), err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("iothub_id", id.IotHubID)
	d.Set("name", id.Name)
	if prop := resp.DeviceSecurityGroupProperties; prop != nil {
		if err := d.Set("allow_rule", flattenIotSecurityDeviceGroupAllowRule(prop.AllowlistRules, d)); err != nil {
			return fmt.Errorf("setting `allow_rule`: %+v", err)
		}
		if err := d.Set("range_rule", flattenIotSecurityDeviceGroupTimeWindowRule(prop.TimeWindowRules)); err != nil {
			return fmt.Errorf("setting `range_rule`: %+v", err)
		}
	}

	return nil
}

func flattenIotSecurityDeviceGroupAllowRule(input *[]security.BasicAllowlistCustomAlertRule, d *pluginsdk.ResourceData) []any {
	if input == nil {
		return []any{}
	}

	var flag bool
	var connectionFromIPsNotAllowed, connectionToIPsNotAllowed, localUsersNotAllowed, processesNotAllowed *[]string
	for _, v := range *input {
		switch v := v.(type) {
		case security.ConnectionFromIPNotAllowed:
			if v.IsEnabled != nil && *v.IsEnabled {
				flag = true
				connectionFromIPsNotAllowed = v.AllowlistValues
			}
		case security.ConnectionToIPNotAllowed:
			if v, ok := d.GetOk("allow_rule.0.connection_to_ips_not_allowed"); ok {
				flag = true
				connectionToIPsNotAllowed = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		case security.LocalUserNotAllowed:
			if v, ok := d.GetOk("allow_rule.0.local_users_not_allowed"); ok {
				flag = true
				localUsersNotAllowed = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		case security.ProcessNotAllowed:
			if v, ok := d.GetOk("allow_rule.0.processes_not_allowed"); ok {
				flag = true
				processesNotAllowed = pluginsdk.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		}
	}
	if !flag {
		return []any{}
	}
	return []any{
		map[string]any{
			"connection_from_ips_not_allowed": pluginsdk.FlattenSlice(connectionFromIPsNotAllowed),
			"connection_to_ips_not_allowed":   pluginsdk.FlattenSlice(connectionToIPsNotAllowed),
			"local_users_not_allowed":         pluginsdk.FlattenSlice(localUsersNotAllowed),
			"processes_not_allowed":           pluginsdk.FlattenSlice(processesNotAllowed),
		},
	}
}

// issue to track: https://github.com/Azure/azure-sdk-for-go/issues/14282
// there is a lot of repeated codes here. once the issue is resolved, we should use a more elegant way
func flattenIotSecurityDeviceGroupTimeWindowRule(input *[]security.BasicTimeWindowCustomAlertRule) []any {
	if input == nil {
		return []any{}
	}
	result := make([]any, 0)

	for _, v := range *input {
		var isEnabled *bool
		var t string
		var timeWindowSizePointer *string
		var minThresholdPointer, maxThresholdPointer *int32
		switch v := v.(type) {
		case security.ActiveConnectionsNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.AmqpC2DMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.MqttC2DMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.HTTPC2DMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.AmqpC2DRejectedMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.MqttC2DRejectedMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.HTTPC2DRejectedMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.AmqpD2CMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.MqttD2CMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.HTTPD2CMessagesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.DirectMethodInvokesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.FailedLocalLoginsNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.FileUploadsNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.QueuePurgesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.TwinUpdatesNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		case security.UnauthorizedOperationsNotInAllowedRange:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			timeWindowSizePointer = v.TimeWindowSize
			minThresholdPointer = v.MinThreshold
			maxThresholdPointer = v.MaxThreshold
		default:
			continue
		}
		if isEnabled == nil || !*isEnabled {
			continue
		}

		var duration string
		var min, max int
		if timeWindowSizePointer != nil {
			duration = *timeWindowSizePointer
		}
		if minThresholdPointer != nil {
			min = int(*minThresholdPointer)
		}
		if maxThresholdPointer != nil {
			max = int(*maxThresholdPointer)
		}
		result = append(result, map[string]any{
			"type":     t,
			"duration": duration,
			"min":      min,
			"max":      max,
		})
	}
	return result
}
