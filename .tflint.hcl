# Copyright IBM Corp. 2014, 2026
# SPDX-License-Identifier: MPL-2.0

# tflint config for the terraform embedded in acceptance tests / website docs and the examples,
# run via 'make hcl-lint' (scripts/checks/tflint-embedded.sh). The tflint binary version is
# pinned in the GNUmakefile.
plugin "azurerm" {
  enabled = true
  version = "0.32.0"
  source  = "github.com/terraform-linters/tflint-ruleset-azurerm"
}

rule "terraform_comment_syntax" {
  enabled = true
}

rule "terraform_required_providers" {
  enabled = false
}

rule "terraform_required_version" {
  enabled = false
}

# Terraform rules - https://github.com/terraform-linters/tflint-ruleset-terraform/blob/main/docs/rules
# test configs and doc examples routinely declare data sources / variables they never reference,
# and leave variables untyped
rule "terraform_unused_declarations" {
  enabled = false
}

rule "terraform_typed_variables" {
  enabled = false
}

# azurerm rules - https://github.com/terraform-linters/tflint-ruleset-azurerm/tree/master/docs/rules
# opinionated rules that do not apply to test configs / doc examples
rule "azurerm_resources_missing_prevent_destroy" {
  enabled = false
}

rule "azurerm_app_service_missing_auto_heal_setting" {
  enabled = false
}

# tests deliberately use small/old vm sizes
rule "azurerm_linux_virtual_machine_retired_size" {
  enabled = false
}

rule "azurerm_windows_virtual_machine_retired_size" {
  enabled = false
}

# stale in ruleset 0.32.0: the provider accepts 1.1 and StandardV2
rule "azurerm_stream_analytics_job_invalid_compatibility_level" {
  enabled = false
}

rule "azurerm_public_ip_invalid_sku" {
  enabled = false
}
