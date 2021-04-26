package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	iothubValidate "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/iothub/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/pluginsdk"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceIotSecurityDeviceGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceIotSecurityDeviceGroupCreateUpdate,
		Read:   resourceIotSecurityDeviceGroupRead,
		Update: resourceIotSecurityDeviceGroupCreateUpdate,
		Delete: resourceIotSecurityDeviceGroupDelete,

		Importer: pluginsdk.ImporterValidatingResourceId(func(id string) error {
			_, err := parse.IotSecurityDeviceGroupID(id)
			return err
		}),

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringIsNotEmpty,
			},

			"iothub_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: iothubValidate.IotHubID,
			},

			"allow_rule": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"connection_to_ip_not_allowed": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type:         schema.TypeString,
								ValidateFunc: validate.CIDR,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_to_ip_not_allowed", "allow_rule.0.local_user_not_allowed", "allow_rule.0.process_not_allowed"},
						},

						"local_user_not_allowed": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_to_ip_not_allowed", "allow_rule.0.local_user_not_allowed", "allow_rule.0.process_not_allowed"},
						},

						"process_not_allowed": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: []string{"allow_rule.0.connection_to_ip_not_allowed", "allow_rule.0.local_user_not_allowed", "allow_rule.0.process_not_allowed"},
						},
					},
				},
			},

			"range_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
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
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"min": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"duration": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: validate.ISO8601Duration,
						},
					},
				},
			},
		},
	}
}

func resourceIotSecurityDeviceGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
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

	timeWindowRules, err := expandIotSecurityDeviceGroupTimeWindowRule(d.Get("range_rule").(*schema.Set).List())
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

func resourceIotSecurityDeviceGroupRead(d *schema.ResourceData, meta interface{}) error {
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
		if err := d.Set("allow_rule", flattenIotSecurityDeviceGroupAllowRule(prop.AllowlistRules)); err != nil {
			return fmt.Errorf("setting `allow_rule`: %+v", err)
		}
		if err := d.Set("range_rule", flattenIotSecurityDeviceGroupTimeWindowRule(prop.TimeWindowRules)); err != nil {
			return fmt.Errorf("setting `range_rule`: %+v", err)
		}
	}

	return nil
}

func resourceIotSecurityDeviceGroupDelete(d *schema.ResourceData, meta interface{}) error {
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

	if connectionToIPNotAllowed := v["connection_to_ip_not_allowed"].(*schema.Set).List(); len(connectionToIPNotAllowed) > 0 {
		result = append(result, security.ConnectionToIPNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(connectionToIPNotAllowed),
			IsEnabled:       utils.Bool(true),
		})
	}
	if LocalUserNotAllowed := v["local_user_not_allowed"].(*schema.Set).List(); len(LocalUserNotAllowed) > 0 {
		result = append(result, security.LocalUserNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(LocalUserNotAllowed),
			IsEnabled:       utils.Bool(true),
		})
	}
	if processNotAllowed := v["process_not_allowed"].(*schema.Set).List(); len(processNotAllowed) > 0 {
		result = append(result, security.ProcessNotAllowed{
			AllowlistValues: utils.ExpandStringSlice(processNotAllowed),
			IsEnabled:       utils.Bool(true),
		})
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

func flattenIotSecurityDeviceGroupAllowRule(input *[]security.BasicAllowlistCustomAlertRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}

	var flag bool
	var connectionToIpNotAllowed, localUserNotAllowed, processNotAllowed *[]string
	for _, v := range *input {
		switch v := v.(type) {
		case security.ConnectionToIPNotAllowed:
			if v.IsEnabled != nil && *v.IsEnabled {
				flag = true
				connectionToIpNotAllowed = v.AllowlistValues
			}
		case security.LocalUserNotAllowed:
			if v.IsEnabled != nil && *v.IsEnabled {
				flag = true
				localUserNotAllowed = v.AllowlistValues
			}
		case security.ProcessNotAllowed:
			if v.IsEnabled != nil && *v.IsEnabled {
				flag = true
				processNotAllowed = v.AllowlistValues
			}
		}
	}
	if !flag {
		return []interface{}{}
	}
	return []interface{}{
		map[string]interface{}{
			"connection_to_ip_not_allowed": utils.FlattenStringSlice(connectionToIpNotAllowed),
			"local_user_not_allowed":       utils.FlattenStringSlice(localUserNotAllowed),
			"process_not_allowed":          utils.FlattenStringSlice(processNotAllowed),
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
