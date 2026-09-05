// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package function_test

import (
	"fmt"
)

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
