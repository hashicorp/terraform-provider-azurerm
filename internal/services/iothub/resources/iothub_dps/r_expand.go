// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package iothub_dps

import (
	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-helpers/resourcemanager/location"
	"github.com/hashicorp/go-azure-sdk/resource-manager/deviceprovisioningservices/2022-02-05/iotdpsresource"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

func expandIoTHubDPSSku(d *pluginsdk.ResourceData) iotdpsresource.IotDpsSkuInfo {
	skuList := d.Get("sku").([]any)
	skuMap := skuList[0].(map[string]any)

	return iotdpsresource.IotDpsSkuInfo{
		Name:     pointer.ToEnum[iotdpsresource.IotDpsSku](skuMap["name"].(string)),
		Capacity: pointer.To(int64(skuMap["capacity"].(int))),
	}
}

func expandIoTHubDPSIoTHubs(input []any) *[]iotdpsresource.IotHubDefinitionDescription {
	linkedHubs := make([]iotdpsresource.IotHubDefinitionDescription, 0)

	for _, attr := range input {
		linkedHubConfig := attr.(map[string]any)
		linkedHub := iotdpsresource.IotHubDefinitionDescription{
			ConnectionString:      linkedHubConfig["connection_string"].(string),
			AllocationWeight:      pointer.To(int64(linkedHubConfig["allocation_weight"].(int))),
			ApplyAllocationPolicy: pointer.To(linkedHubConfig["apply_allocation_policy"].(bool)),
			Location:              location.Normalize(linkedHubConfig["location"].(string)),
		}

		linkedHubs = append(linkedHubs, linkedHub)
	}

	return &linkedHubs
}

func expandDpsIPFilterRules(d *pluginsdk.ResourceData) *[]iotdpsresource.IPFilterRule {
	ipFilterRuleList := d.Get("ip_filter_rule").([]any)
	if len(ipFilterRuleList) == 0 {
		return nil
	}

	rules := make([]iotdpsresource.IPFilterRule, 0)

	for _, r := range ipFilterRuleList {
		rawRule := r.(map[string]any)
		rule := &iotdpsresource.IPFilterRule{
			FilterName: rawRule["name"].(string),
			Action:     iotdpsresource.IPFilterActionType(rawRule["action"].(string)),
			IPMask:     rawRule["ip_mask"].(string),
			Target:     pointer.ToEnum[iotdpsresource.IPFilterTargetType](rawRule["target"].(string)),
		}

		rules = append(rules, *rule)
	}
	return &rules
}
