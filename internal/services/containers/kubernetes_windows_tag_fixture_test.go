// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestKubernetesClusterNodePoolWindowsTagFixture(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "eastus"},
	}
	resource := KubernetesClusterNodePoolResource{}
	for _, tagValue := range []string{"dev", "prod"} {
		t.Run(tagValue, func(t *testing.T) {
			config := resource.windowsNodePoolWithTags(data, tagValue)
			if !strings.Contains(config, resource.templateWindowsConfig(data)) {
				t.Fatal("Windows tag fixture must retain the Windows cluster prerequisites")
			}
			if !regexp.MustCompile(`os_sku\s*=\s*"Windows2022"`).MatchString(config) {
				t.Fatal("Windows tag fixture must select the tested Windows2022 image")
			}
			if !strings.Contains(config, fmt.Sprintf("Environment = %q", tagValue)) {
				t.Fatal("Windows tag fixture must render the requested tag value")
			}
		})
	}
}
