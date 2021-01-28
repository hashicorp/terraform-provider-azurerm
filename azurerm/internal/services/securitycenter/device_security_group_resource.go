package securitycenter

import (
	"fmt"
	"log"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/preview/security/mgmt/v3.0/security"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/azure"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/tf"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/clients"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/services/securitycenter/parse"
	azSchema "github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/tf/schema"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/internal/timeouts"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"
)

func resourceDeviceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceDeviceSecurityGroupCreateUpdate,
		Read:   resourceDeviceSecurityGroupRead,
		Update: resourceDeviceSecurityGroupCreateUpdate,
		Delete: resourceDeviceSecurityGroupDelete,

		Importer: azSchema.ValidateResourceIDPriorToImport(func(id string) error {
			_, err := parse.DeviceSecurityGroupID(id)
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

			"target_resource_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: azure.ValidateResourceID,
			},

			"allow_list_rule": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(security.RuleTypeConnectionToIPNotAllowed),
								string(security.RuleTypeLocalUserNotAllowed),
								string(security.RuleTypeProcessNotAllowed),
							}, false),
						},

						"values": {
							Type:     schema.TypeSet,
							Required: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},

			"time_window_rule": {
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

						"max_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"min_threshold": {
							Type:         schema.TypeInt,
							Required:     true,
							ValidateFunc: validation.IntAtLeast(0),
						},

						"time_window_size": {
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

func resourceDeviceSecurityGroupCreateUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForCreateUpdate(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id := parse.NewDeviceSecurityGroupId(d.Get("target_resource_id").(string), d.Get("name").(string))
	if d.IsNewResource() {
		server, err := client.Get(ctx, id.TargetResourceID, id.Name)
		if err != nil {
			if !utils.ResponseWasNotFound(server.Response) {
				return fmt.Errorf("checking for presence of existing Device Security Group for %q: %+v", id.ID(), err)
			}
		}

		if !utils.ResponseWasNotFound(server.Response) {
			return tf.ImportAsExistsError("azurerm_device_security_group", id.ID())
		}
	}

	timeWindowRules, err := expandDeviceSecurityGroupTimeWindowRule(d.Get("time_window_rule").(*schema.Set).List())
	if err != nil {
		return err
	}
	allowListRules, err := expandDeviceSecurityGroupAllowListRule(d.Get("allow_list_rule").(*schema.Set).List())
	if err != nil {
		return err
	}
	deviceSecurityGroup := security.DeviceSecurityGroup{
		DeviceSecurityGroupProperties: &security.DeviceSecurityGroupProperties{
			TimeWindowRules: timeWindowRules,
			AllowlistRules:  allowListRules,
		},
	}

	if _, err := client.CreateOrUpdate(ctx, id.TargetResourceID, id.Name, deviceSecurityGroup); err != nil {
		return fmt.Errorf("creating/updating %s: %+v", id, err)
	}

	d.SetId(id.ID())
	return resourceDeviceSecurityGroupRead(d, meta)
}

func resourceDeviceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForRead(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DeviceSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	resp, err := client.Get(ctx, id.TargetResourceID, id.Name)
	if err != nil {
		if utils.ResponseWasNotFound(resp.Response) {
			log.Printf("Device Security Group not found for %q: %+v", id.ID(), err)
			d.SetId("")
			return nil
		}

		return fmt.Errorf("retrieving %s: %+v", id, err)
	}

	d.Set("target_resource_id", id.TargetResourceID)
	d.Set("name", id.Name)
	if prop := resp.DeviceSecurityGroupProperties; prop != nil {
		if err := d.Set("allow_list_rule", flattenDeviceSecurityGroupAllowListRule(prop.AllowlistRules)); err != nil {
			return fmt.Errorf("setting `allow_list_rule`: %+v", err)
		}
		if err := d.Set("time_window_rule", flattenDeviceSecurityGroupTimeWindowRule(prop.TimeWindowRules)); err != nil {
			return fmt.Errorf("setting `time_window_rule`: %+v", err)
		}
	}

	return nil
}

func resourceDeviceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*clients.Client).SecurityCenter.DeviceSecurityGroupsClient
	ctx, cancel := timeouts.ForDelete(meta.(*clients.Client).StopContext, d)
	defer cancel()

	id, err := parse.DeviceSecurityGroupID(d.Id())
	if err != nil {
		return err
	}

	if _, err := client.Delete(ctx, id.TargetResourceID, id.Name); err != nil {
		return fmt.Errorf("deleting %s: %+v", id, err)
	}

	return nil
}

func expandDeviceSecurityGroupAllowListRule(input []interface{}) (*[]security.BasicAllowlistCustomAlertRule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make([]security.BasicAllowlistCustomAlertRule, 0)
	ruleTypeMap := make(map[security.RuleTypeBasicCustomAlertRule]struct{})
	for _, item := range input {
		v := item.(map[string]interface{})
		t := security.RuleTypeBasicCustomAlertRule(v["type"].(string))
		values := v["values"].(*schema.Set).List()

		// check duplicate
		if _, ok := ruleTypeMap[t]; ok {
			return nil, fmt.Errorf("rule type duplicate: %q", t)
		}
		ruleTypeMap[t] = struct{}{}

		switch t {
		case security.RuleTypeConnectionToIPNotAllowed:
			result = append(result, security.ConnectionToIPNotAllowed{
				AllowlistValues: utils.ExpandStringSlice(values),
				IsEnabled:       utils.Bool(true),
			})
		case security.RuleTypeLocalUserNotAllowed:
			result = append(result, security.LocalUserNotAllowed{
				AllowlistValues: utils.ExpandStringSlice(values),
				IsEnabled:       utils.Bool(true),
			})
		case security.RuleTypeProcessNotAllowed:
			result = append(result, security.ProcessNotAllowed{
				AllowlistValues: utils.ExpandStringSlice(values),
				IsEnabled:       utils.Bool(true),
			})
		}
	}
	return &result, nil
}

func expandDeviceSecurityGroupTimeWindowRule(input []interface{}) (*[]security.BasicTimeWindowCustomAlertRule, error) {
	if len(input) == 0 {
		return nil, nil
	}
	result := make([]security.BasicTimeWindowCustomAlertRule, 0)
	ruleTypeMap := make(map[security.RuleTypeBasicCustomAlertRule]struct{})
	for _, item := range input {
		v := item.(map[string]interface{})
		t := security.RuleTypeBasicCustomAlertRule(v["type"].(string))
		timeWindowSize := v["time_window_size"].(string)
		minThreshold := int32(v["min_threshold"].(int))
		maxThreshold := int32(v["max_threshold"].(int))

		// check duplicate
		if _, ok := ruleTypeMap[t]; ok {
			return nil, fmt.Errorf("rule type duplicate: %q", t)
		}
		ruleTypeMap[t] = struct{}{}

		switch t {
		case security.RuleTypeActiveConnectionsNotInAllowedRange:
			result = append(result, security.ActiveConnectionsNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpC2DMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttC2DMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPC2DMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.AmqpC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.MqttC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPC2DRejectedMessagesNotInAllowedRange:
			result = append(result, security.HTTPC2DRejectedMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeAmqpD2CMessagesNotInAllowedRange:
			result = append(result, security.AmqpD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeMqttD2CMessagesNotInAllowedRange:
			result = append(result, security.MqttD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeHTTPD2CMessagesNotInAllowedRange:
			result = append(result, security.HTTPD2CMessagesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeDirectMethodInvokesNotInAllowedRange:
			result = append(result, security.DirectMethodInvokesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeFailedLocalLoginsNotInAllowedRange:
			result = append(result, security.FailedLocalLoginsNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeFileUploadsNotInAllowedRange:
			result = append(result, security.FileUploadsNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeQueuePurgesNotInAllowedRange:
			result = append(result, security.QueuePurgesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeTwinUpdatesNotInAllowedRange:
			result = append(result, security.TwinUpdatesNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		case security.RuleTypeUnauthorizedOperationsNotInAllowedRange:
			result = append(result, security.UnauthorizedOperationsNotInAllowedRange{
				TimeWindowSize: utils.String(timeWindowSize),
				MinThreshold:   utils.Int32(minThreshold),
				MaxThreshold:   utils.Int32(maxThreshold),
				IsEnabled:      utils.Bool(true),
			})
		}
	}
	return &result, nil
}

func flattenDeviceSecurityGroupAllowListRule(input *[]security.BasicAllowlistCustomAlertRule) []interface{} {
	if input == nil {
		return []interface{}{}
	}
	result := make([]interface{}, 0)
	for _, v := range *input {
		var isEnabled *bool
		var t string
		var values *[]string
		switch v := v.(type) {
		case security.ConnectionToIPNotAllowed:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			values = v.AllowlistValues
		case security.LocalUserNotAllowed:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			values = v.AllowlistValues
		case security.ProcessNotAllowed:
			isEnabled = v.IsEnabled
			t = string(v.RuleType)
			values = v.AllowlistValues
		default:
			continue
		}
		if isEnabled == nil || !*isEnabled {
			continue
		}
		result = append(result, map[string]interface{}{
			"type":   t,
			"values": utils.FlattenStringSlice(values),
		})

	}
	return result
}

func flattenDeviceSecurityGroupTimeWindowRule(input *[]security.BasicTimeWindowCustomAlertRule) []interface{} {
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

		var timeWindowSize string
		var minThreshold, maxThreshold int
		if timeWindowSizePointer != nil {
			timeWindowSize = *timeWindowSizePointer
		}
		if minThresholdPointer != nil {
			minThreshold = int(*minThresholdPointer)
		}
		if maxThresholdPointer != nil {
			maxThreshold = int(*maxThresholdPointer)
		}
		result = append(result, map[string]interface{}{
			"type":             t,
			"time_window_size": timeWindowSize,
			"min_threshold":    minThreshold,
			"max_threshold":    maxThreshold,
		})

	}
	return result
}
