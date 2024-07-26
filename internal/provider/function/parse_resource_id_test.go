// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package function_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/features"
	"github.com/hashicorp/terraform-provider-azurerm/internal/provider/framework"
)

func TestProviderFunctionParseResourceID_basic(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skipf("skipping test due to missing feature flag")
	}
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: testParseResourceIdOutput("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1"),
				Check: acceptance.ComposeTestCheckFunc(
					acceptance.TestCheckOutput("subscription_id", "12345678-1234-9876-4563-123456789012"),
					acceptance.TestCheckOutput("resource_provider", "Microsoft.ApiManagement"),
					acceptance.TestCheckOutput("resource_group_name", "resGroup1"),
					acceptance.TestCheckOutput("resource_type", "hostnameConfigurations"),
					acceptance.TestCheckOutput("resource_name", "config1"),
					acceptance.TestCheckOutput("resource_scope", ""),
					acceptance.TestCheckOutput("full_resource_type", "Microsoft.ApiManagement/service/gateways/hostnameConfigurations"),
					acceptance.TestCheckOutput("service_name", "service1"),
					acceptance.TestCheckOutput("gateway_name", "gateway1"),
				),
			},
		},
	})
}

func TestProviderFunctionParseResourceID_scopedAtSubscription(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skipf("skipping test due to missing feature flag")
	}
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: testParseScopedResourceIdOutput("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Chaos/targets/target1"),
				Check: acceptance.ComposeTestCheckFunc(
					acceptance.TestCheckOutput("resource_provider", "Microsoft.Chaos"),
					acceptance.TestCheckOutput("resource_type", "targets"),
					acceptance.TestCheckOutput("resource_name", "target1"),
					acceptance.TestCheckOutput("full_resource_type", "Microsoft.Chaos/targets"),
					acceptance.TestCheckOutput("resource_scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1"),
				),
			},
		},
	})
}

func TestProviderFunctionParseResourceID_scopedAtResource(t *testing.T) {
	if !features.FourPointOhBeta() {
		t.Skipf("skipping test due to missing feature flag")
	}
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(version.Must(version.NewVersion("1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: framework.ProtoV5ProviderFactoriesInit(context.Background(), "azurerm"),
		Steps: []resource.TestStep{
			{
				Config: testParseScopedResourceIdOutput("/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/mystorageaccount/providers/Microsoft.EventGrid/eventSubscriptions/event1"),
				Check: acceptance.ComposeTestCheckFunc(
					acceptance.TestCheckOutput("resource_provider", "Microsoft.EventGrid"),
					acceptance.TestCheckOutput("resource_type", "eventSubscriptions"),
					acceptance.TestCheckOutput("resource_name", "event1"),
					acceptance.TestCheckOutput("full_resource_type", "Microsoft.EventGrid/eventSubscriptions"),
					acceptance.TestCheckOutput("resource_scope", "/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.Storage/storageAccounts/mystorageaccount"),
				),
			},
		},
	})
}

func testParseResourceIdOutput(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  parsed_id             = provider::azurerm::parse_resource_id("%s")
  parent_resource_name1 = local.parsed_id["parent_resources"]["service"]
  parent_resource_name2 = local.parsed_id["parent_resources"]["gateways"]
}

output "resource_name" {
  value = local.parsed_id["resource_name"]
}

output "resource_provider" {
  value = local.parsed_id["resource_provider"]
}
output "resource_scope" {
  value = local.parsed_id["resource_scope"]
}

output "resource_group_name" {
  value = local.parsed_id["resource_group_name"]
}

output "resource_type" {
  value = local.parsed_id["resource_type"]
}

output "service_name" {
  value = local.parent_resource_name1
}

output "gateway_name" {
  value = local.parent_resource_name2
}

output "subscription_id" {
  value = local.parsed_id["subscription_id"]
}

output "full_resource_type" {
  value = local.parsed_id["full_resource_type"]
}


`, id)
}

func testParseScopedResourceIdOutput(id string) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

locals {
  parsed_id = provider::azurerm::parse_resource_id("%s")
}

output "resource_name" {
  value = local.parsed_id["resource_name"]
}

output "resource_provider" {
  value = local.parsed_id["resource_provider"]
}

output "resource_scope" {
  value = local.parsed_id["resource_scope"]
}

output "resource_group_name" {
  value = local.parsed_id["resource_group_name"]
}

output "resource_type" {
  value = local.parsed_id["resource_type"]
}

output "subscription_id" {
  value = local.parsed_id["subscription_id"]
}

output "full_resource_type" {
  value = local.parsed_id["full_resource_type"]
}


`, id)
}
