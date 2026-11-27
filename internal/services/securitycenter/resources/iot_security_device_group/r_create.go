// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package iot_security_device_group

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" //nolint:staticcheck
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients/account"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	securitycenterClient "github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/client"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
)

func resourceIotSecurityDeviceGroupCreateUpdate(d *pluginsdk.ResourceData, meta any) error {
	client := securitycenterClient.From(meta).DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(account.StopContext(meta), d)
	defer cancel()

	id := parse.NewIotSecurityDeviceGroupId(d.Get("iothub_id").(string), d.Get("name").(string))
	if d.IsNewResource() {
		if !features.From(meta).SkipImportCheckOnCreateAndAllowOverwritingExistingResources {
			server, err := client.Get(ctx, id.IotHubID, id.Name)
			if err != nil {
				if !response.WasNotFound(server.Response.Response) {
					return fmt.Errorf("checking for presence of existing Device Security Group for %q: %+v", id.ID(), err)
				}
			}

			if !response.WasNotFound(server.Response.Response) {
				return tf.ImportAsExistsError("azurerm_iot_security_device_group", id.ID())
			}
		}
	}

	timeWindowRules, err := expandIotSecurityDeviceGroupTimeWindowRule(d.Get("range_rule").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	deviceSecurityGroup := security.DeviceSecurityGroup{
		DeviceSecurityGroupProperties: &security.DeviceSecurityGroupProperties{
			TimeWindowRules: timeWindowRules,
			AllowlistRules:  expandIotSecurityDeviceGroupAllowRule(d.Get("allow_rule").([]any)),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.IotHubID, id.Name, deviceSecurityGroup); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceIotSecurityDeviceGroupRead(d, meta)
}

func expandIotSecurityDeviceGroupAllowRule(input []any) *[]security.BasicAllowlistCustomAlertRule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]any)
	result := make([]security.BasicAllowlistCustomAlertRule, 0)

	if connectionFromIPNotAllowed := v["connection_from_ips_not_allowed"].(*pluginsdk.Set).List(); len(connectionFromIPNotAllowed) > 0 {
		result = append(result, security.ConnectionFromIPNotAllowed{
			AllowlistValues: pluginsdk.ExpandStringSlice(connectionFromIPNotAllowed),
			IsEnabled:       pointer.To(true),
		})
	}

	var connectionToIPListNotAllowed *security.ConnectionToIPNotAllowed
	if connectionToIPsNotAllowed := v["connection_to_ips_not_allowed"].(*pluginsdk.Set).List(); len(connectionToIPsNotAllowed) > 0 {
		connectionToIPListNotAllowed = &security.ConnectionToIPNotAllowed{
			AllowlistValues: pluginsdk.ExpandStringSlice(connectionToIPsNotAllowed),
			IsEnabled:       pointer.To(true),
		}
	}
	if connectionToIPListNotAllowed != nil {
		result = append(result, *connectionToIPListNotAllowed)
	}

	var localUserListNotAllowed *security.LocalUserNotAllowed
	if LocalUsersNotAllowed := v["local_users_not_allowed"].(*pluginsdk.Set).List(); len(LocalUsersNotAllowed) > 0 {
		localUserListNotAllowed = &security.LocalUserNotAllowed{
			AllowlistValues: pluginsdk.ExpandStringSlice(LocalUsersNotAllowed),
			IsEnabled:       pointer.To(true),
		}
	}
	if localUserListNotAllowed != nil {
		result = append(result, *localUserListNotAllowed)
	}

	var processListNotAllowed *security.ProcessNotAllowed
	if processesNotAllowed := v["processes_not_allowed"].(*pluginsdk.Set).List(); len(processesNotAllowed) > 0 {
		processListNotAllowed = &security.ProcessNotAllowed{
			AllowlistValues: pluginsdk.ExpandStringSlice(processesNotAllowed),
			IsEnabled:       pointer.To(true),
		}
	}
	if processListNotAllowed != nil {
		result = append(result, *processListNotAllowed)
	}

	return &result
}

// issue to track: https://github.com/Azure/azure-sdk-for-go/issues/14282
// there is a lot of repeated codes here. once the issue is resolved, we should use a more elegant way
func expandIotSecurityDeviceGroupTimeWindowRule(input []any) (*[]security.BasicTimeWindowCustomAlertRule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make([]security.BasicTimeWindowCustomAlertRule, 0)
	ruleTypeMap := make(map[security.RuleTypeBasicCustomAlertRule]struct{})
	for _, item := range input {
		v := item.(map[string]any)
		t := security.RuleTypeBasicCustomAlertRule(v["type"].(string))
		duration := v["duration"].(string)
		min := int32(v["min"].(int))
		max := int32(v["max"].(int))

		// check duplicate
		if _, ok := ruleTypeMap[t]; ok {
			return nil, fmt.Errorf("rule type duplicate: %q", t)
		}
		ruleTypeMap[t] = struct{}{}

		switch t {
		case security.RuleTypeActiveConnectionsNotInAllowedRange:
			result = append(result, security.ActiveConnectionsNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeAmqpC2DMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeMqttC2DMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeHTTPC2DMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeAmqpC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeMqttC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeHTTPC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeAmqpD2CMessagesNotInAllowedRange:
			result = append(result, security.AmqpD2CMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeMqttD2CMessagesNotInAllowedRange:
			result = append(result, security.MqttD2CMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeHTTPD2CMessagesNotInAllowedRange:
			result = append(result, security.HTTPD2CMessagesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeDirectMethodInvokesNotInAllowedRange:
			result = append(result, security.DirectMethodInvokesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeFailedLocalLoginsNotInAllowedRange:
			result = append(result, security.FailedLocalLoginsNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeFileUploadsNotInAllowedRange:
			result = append(result, security.FileUploadsNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeQueuePurgesNotInAllowedRange:
			result = append(result, security.QueuePurgesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeTwinUpdatesNotInAllowedRange:
			result = append(result, security.TwinUpdatesNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		case security.RuleTypeUnauthorizedOperationsNotInAllowedRange:
			result = append(result, security.UnauthorizedOperationsNotInAllowedRange{
				TimeWindowSize: pointer.To(duration),
				MinThreshold:   pointer.To(min),
				MaxThreshold:   pointer.To(max),
				IsEnabled:      pointer.To(true),
			})
		}
	}
	return &result, nil
}
