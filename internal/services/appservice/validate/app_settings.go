// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package validate

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
)

var UnmanagedSettings = []string{
	"DOCKER_REGISTRY_SERVER_URL",
	"DOCKER_REGISTRY_SERVER_USERNAME",
	"DOCKER_REGISTRY_SERVER_PASSWORD",
	"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
	"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
	"WEBSITE_HTTPLOGGING_CONTAINER_URL",
	"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
	"WEBSITE_VNET_ROUTE_ALL",
	"spring.datasource.password",
	"spring.datasource.url",
	"spring.datasource.username",
	"WEBSITE_HEALTHCHECK_MAXPINGFAILURES",
}
var UnmanagedSettingsDeprecated = []string{
	"DIAGNOSTICS_AZUREBLOBCONTAINERSASURL",
	"DIAGNOSTICS_AZUREBLOBRETENTIONINDAYS",
	"WEBSITE_HTTPLOGGING_CONTAINER_URL",
	"WEBSITE_HTTPLOGGING_RETENTION_DAYS",
	"WEBSITE_VNET_ROUTE_ALL",
	"spring.datasource.password",
	"spring.datasource.url",
	"spring.datasource.username",
	"WEBSITE_HEALTHCHECK_MAXPINGFAILURES",
}

func AppSettings(input interface{}, key string) (warnings []string, errors []error) {
	if appSettings, ok := input.(map[string]interface{}); ok {
		for k := range appSettings {
			if !features.FourPointOhBeta() {
				for _, f := range UnmanagedSettingsDeprecated {
					if strings.EqualFold(k, f) {
						errors = append(errors, fmt.Errorf("cannot set a value for %s in %s", k, key))
					}
				}
			} else {
				for _, f := range UnmanagedSettings {
					if strings.EqualFold(k, f) {
						errors = append(errors, fmt.Errorf("cannot set a value for %s in %s", k, key))
					}
				}
			}
		}
	} else {
		errors = append(errors, fmt.Errorf("expected %s to be a map of strings", key))
		return
	}
	return
}
