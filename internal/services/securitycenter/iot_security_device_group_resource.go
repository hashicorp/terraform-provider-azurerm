// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security" // nolint: staticcheck
	"github.com/hashicorp/terraform-provider-azurerm/helpers/tf"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	iothubValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/iothub/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/securitycenter/parse"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
	"github.com/hashicorp/terraform-provider-azurerm/internal/timeouts"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

func resourceIotSecurityDeviceGroup() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Create: resourceIotSecurityDeviceGroupCreateUpdate,
		Read:   resourceIotSecurityDeviceGroupRead,
		Update: resourceIotSecurityDeviceGroupCreateUpdate,
		Delete: resourceIotSecurityDeviceGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IotSecurityDeviceGroupID(id)
			return err
		}),

		Timeouts: &pluginsdk.ResourceTimeout{
			Create: pluginsdk.DefaultTimeout(30 * time.Minute),
			Read:   pluginsdk.DefaultTimeout(5 * time.Minute),
			Update: pluginsdk.DefaultTimeout(30 * time.Minute),
			Delete: pluginsdk.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*pluginsdk.Schema{
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
				ValidateFunc: iothubValidate.IotHubID,
			},

			"allow_rule": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &pluginsdk.Resource{
					Schema: map[string]*pluginsdk.Schema{
						"connection_from_ips_not_allowed": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						},

						"connection_to_ips_not_allowed": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type:         pluginsdk.TypeString,
								ValidateFunc: validate.CIDR,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						},

						"local_users_not_allowed": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
						},

						"processes_not_allowed": {
							Type:     pluginsdk.TypeSet,
							Optional: true,
							Elem: &pluginsdk.Schema{
								Type: pluginsdk.TypeString,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_from_ips_not_allowed", "allow_rule.0.connection_to_ips_not_allowed", "allow_rule.0.local_users_not_allowed", "allow_rule.0.processes_not_allowed"},
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
							ValidateFunc: validate.ISO8601Duration,
						},
					},
				},
			},
		},
	}
}

func resourceIotSecurityDeviceGroupCreateUpdate(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewIotSecurityDeviceGroupId(d.Get("iothub_id").(string), d.Get("name").(string))
	if d.IsNewResource() {
		server, err := client.Get(ctx, id.IotHubID, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("checking for presence of existing Device Security Group for %q: %+v", id.ID(), err)
			}
		}

		if !utils.ResponseWasNotFound(server.Response) {
			return tf.ImportAsExistsError("azurerm_iot_security_device_group", id.ID())
		}
	}

	timeWindowRules, err := expandIotSecurityDeviceGroupTimeWindowRule(d.Get("range_rule").(*pluginsdk.Set).List())
	if err != nil {
		return err
	}
	deviceSecurityGroup := security.DeviceSecurityGroup{
		DeviceSecurityGroupProperties: &security.DeviceSecurityGroupProperties{
			TimeWindowRules: timeWindowRules,
			AllowlistRules:  expandIotSecurityDeviceGroupAllowRule(d.Get("allow_rule").([]interface{})),
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.IotHubID, id.Name, deviceSecurityGroup); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceIotSecurityDeviceGroupRead(d, meta)
}

func resourceIotSecurityDeviceGroupRead(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotSecurityDeviceGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.IotHubID, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
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

func resourceIotSecurityDeviceGroupDelete(d *pluginsdk.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.IotSecurityDeviceGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.IotHubID, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandIotSecurityDeviceGroupAllowRule(input []interface{}) *[]security.BasicAllowlistCustomAlertRule {
	if len(input) == 0 || input[0] == nil {
		return nil
	}

	v := input[0].(map[string]interface{})
	result := make([]security.BasicAllowlistCustomAlertRule, 0)

	if connectionFromIPNotAllowed := v["connection_from_ips_not_allowed"].(*pluginsdk.Set).List(); len(connectionFromIPNotAllowed) > 0 {
		result = append(result, security.ConnectionFromIPNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(connectionFromIPNotAllowed),
			IsEnabled:       utils.Bool(true),
		})
	}

	var connectionToIPListNotAllowed *security.ConnectionToIPNotAllowed
	if connectionToIPsNotAllowed := v["connection_to_ips_not_allowed"].(*pluginsdk.Set).List(); len(connectionToIPsNotAllowed) > 0 {
		connectionToIPListNotAllowed = &security.ConnectionToIPNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(connectionToIPsNotAllowed),
			IsEnabled:       utils.Bool(true),
		}
	}
	if connectionToIPListNotAllowed != nil {
		result = append(result, *connectionToIPListNotAllowed)
	}

	var localUserListNotAllowed *security.LocalUserNotAllowed
	if LocalUsersNotAllowed := v["local_users_not_allowed"].(*pluginsdk.Set).List(); len(LocalUsersNotAllowed) > 0 {
		localUserListNotAllowed = &security.LocalUserNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(LocalUsersNotAllowed),
			IsEnabled:       utils.Bool(true),
		}
	}
	if localUserListNotAllowed != nil {
		result = append(result, *localUserListNotAllowed)
	}

	var processListNotAllowed *security.ProcessNotAllowed
	if processesNotAllowed := v["processes_not_allowed"].(*pluginsdk.Set).List(); len(processesNotAllowed) > 0 {
		processListNotAllowed = &security.ProcessNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(processesNotAllowed),
			IsEnabled:       utils.Bool(true),
		}
	}
	if processListNotAllowed != nil {
		result = append(result, *processListNotAllowed)
	}

	return &result
}

// issue to track: https://github.com/Azure/azure-sdk-for-go/issues/14282
// there is a lot of repeated codes here. once the issue is resolved, we should use a more elegant way
func expandIotSecurityDeviceGroupTimeWindowRule(input []interface{}) (*[]security.BasicTimeWindowCustomAlertRule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make([]security.BasicTimeWindowCustomAlertRule, 0)
	ruleTypeMap := make(map[security.RuleTypeBasicCustomAlertRule]struct{})
	for _, item := range input {
		v := item.(map[string]interface{})
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
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpC2DMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttC2DMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPC2DMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpD2CMessagesNotInAllowedRange:
			result = append(result, security.AmqpD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttD2CMessagesNotInAllowedRange:
			result = append(result, security.MqttD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPD2CMessagesNotInAllowedRange:
			result = append(result, security.HTTPD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeDirectMethodInvokesNotInAllowedRange:
			result = append(result, security.DirectMethodInvokesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeFailedLocalLoginsNotInAllowedRange:
			result = append(result, security.FailedLocalLoginsNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeFileUploadsNotInAllowedRange:
			result = append(result, security.FileUploadsNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeQueuePurgesNotInAllowedRange:
			result = append(result, security.QueuePurgesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeTwinUpdatesNotInAllowedRange:
			result = append(result, security.TwinUpdatesNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeUnauthorizedOperationsNotInAllowedRange:
			result = append(result, security.UnauthorizedOperationsNotInAllowedRange{
				TimeWindowSize: utils.String(duration),
				MinThreshold:   utils.Int32(min),
				MaxThreshold:   utils.Int32(max),
				IsEnabled:      utils.Bool(true),
			})
		}
	}
	return &result, nil
}

func flattenIotSecurityDeviceGroupAllowRule(input *[]security.BasicAllowlistCustomAlertRule, d *pluginsdk.ResourceData) []interface{} {
	if input == nil {
		return []interface{}{}
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
				connectionToIPsNotAllowed = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		case security.LocalUserNotAllowed:
			if v, ok := d.GetOk("allow_rule.0.local_users_not_allowed"); ok {
				flag = true
				localUsersNotAllowed = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		case security.ProcessNotAllowed:
			if v, ok := d.GetOk("allow_rule.0.processes_not_allowed"); ok {
				flag = true
				processesNotAllowed = utils.ExpandStringSlice(v.(*pluginsdk.Set).List())
			}
		}
	}
	if !flag {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"connection_from_ips_not_allowed": utils.FlattenStringSlice(connectionFromIPsNotAllowed),
			"connection_to_ips_not_allowed":   utils.FlattenStringSlice(connectionToIPsNotAllowed),
			"local_users_not_allowed":         utils.FlattenStringSlice(localUsersNotAllowed),
			"processes_not_allowed":           utils.FlattenStringSlice(processesNotAllowed),
		},
	}
}

// issue to track: https://github.com/Azure/azure-sdk-for-go/issues/14282
// there is a lot of repeated codes here. once the issue is resolved, we should use a more elegant way
func flattenIotSecurityDeviceGroupTimeWindowRule(input *[]security.BasicTimeWindowCustomAlertRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)

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
		result = append(result, map[string]interface{}{
			"type":     t,
			"duration": duration,
			"min":      min,
			"max":      max,
		})
	}
	return result
}
