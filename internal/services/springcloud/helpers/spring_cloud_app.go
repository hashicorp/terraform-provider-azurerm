// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package helpers

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	appplatform_rm "github.com/hashicorp/go-azure-sdk/resource-manager/appplatform/2024-01-01-preview/appplatform"
)

func ExpandSpringCloudAppAddon(input string) (map[string]any, error) {
	var addonConfig map[string]any
	if len(input) != 0 {
		if err := json.Unmarshal([]byte(input), &addonConfig); err != nil {
			return nil, fmt.Errorf("unable to unmarshal `addon_json`: %+v", err)
		}
	}
	return addonConfig, nil
}

func FlattenSpringCloudAppAddon(configs map[string]any) *string {
	if len(configs) == 0 {
		return nil
	}
	// The returned value has inconsistent casing
	// TODO: Remove the normalization codes once the following issue is fixed.
	// Issue: https://github.com/Azure/azure-rest-api-specs/issues/22481
	if raw, ok := configs["applicationConfigurationService"]; ok && raw != nil {
		if applicationConfigurationService, ok := raw.(map[string]any); ok && len(applicationConfigurationService) != 0 {
			if resourceId, ok := applicationConfigurationService["resourceId"]; ok && resourceId != nil {
				applicationConfigurationServiceId, err := appplatform_rm.ParseConfigurationServiceIDInsensitively(resourceId.(string))
				if err == nil {
					applicationConfigurationService["resourceId"] = applicationConfigurationServiceId.ID()
					configs["applicationConfigurationService"] = applicationConfigurationService
				}
			}
		}
	}
	if raw, ok := configs["serviceRegistry"]; ok && raw != nil {
		if serviceRegistry, ok := raw.(map[string]any); ok && len(serviceRegistry) != 0 {
			if resourceId, ok := serviceRegistry["resourceId"]; ok && resourceId != nil {
				serviceRegistryId, err := appplatform_rm.ParseServiceRegistryIDInsensitively(resourceId.(string))
				if err == nil {
					serviceRegistry["resourceId"] = serviceRegistryId.ID()
					configs["serviceRegistry"] = serviceRegistry
				}
			}
		}
	}
	addonConfig, _ := json.Marshal(configs)
	return pointer.To(string(addonConfig))
}
