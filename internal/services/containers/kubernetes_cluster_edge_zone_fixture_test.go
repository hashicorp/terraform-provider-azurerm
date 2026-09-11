// Copyright IBM Corp. 2014, 2026
// SPDX-License-Identifier: MPL-2.0

package containers_test

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
)

func TestKubernetesClusterEdgeZoneFixture(t *testing.T) {
	data := acceptance.TestData{
		RandomInteger: 12345,
		Locations:     acceptance.Regions{Primary: "westus"},
	}
	config := KubernetesClusterResource{}.edgeZone(data, "1.35.0", "Test1")
	if !regexp.MustCompile(`vm_size\s*=\s*"Standard_D2s_v3"`).MatchString(config) {
		t.Fatal("edge-zone fixture must use Standard_D2s_v3")
	}
}
