# Example input file for the scaff "generate" command:
#
#   go run ./internal/tools/scaff generate -input="internal/tools/scaff/examples/generate.hcl"
#
# Each "resource" block scaffolds a typed (internal/sdk) resource from the
# Pandora Data API, and optionally a list resource and/or a data source. The
# block label is the terraform resource name (snake_case, without the provider
# prefix, e.g. "redhat_openshift_cluster" -> azurerm_redhat_openshift_cluster).
#
# A running Pandora Data API is required (see -pandora-url / pandora_url below).
# ---------------------------------------------------------------------------
# Optional global overrides. When omitted these fall back to .scaff.hcl (or the
# tool's built-in defaults), so you usually don't need to set them here.
# ---------------------------------------------------------------------------
# pandora_url = "http://localhost:8080" # Pandora Data API base URL
# output_path = "internal/services"     # base output dir, joined with each servicepackage
# provider    = "azurerm"               # provider name (terraform type prefix + import paths)
# org         = "hashicorp"             # provider GitHub org (import paths)
overwrite   = true                   # set true to replace existing generated files

# resource "redhat_openshift_cluster" {
#   arm_type       = "Microsoft.RedHatOpenShift/openShiftClusters"
#   go_name        = "RedHatOpenShiftCluster"
#   servicepackage = "redhatopenshift"
#   api_version    = "2025-07-25"
#   list           = true
#   data_source    = true
# }
#
# resource "storage_mover" {
#   arm_type       = "Microsoft.StorageMover/storageMovers"
#   servicepackage = "storagemover"
# }
#
# resource "netapp_account" {
#   arm_type       = "Microsoft.NetApp/netAppAccounts"
#   go_name        = "NetAppAccount"
#   servicepackage = "netapp"
#   data_source    = true
# }

resource "automation_account" {
  arm_type       = "Microsoft.Automation/automationAccounts"
  go_name        = "AutomationAccount"
  servicepackage = "automation"
  api_version    = "2024-10-23"
  list           = true
  data_source    = true
}
#
# resource "container_registry" {
#   arm_type       = "Microsoft.ContainerRegistry/registries"
#   go_name        = "ContainerRegistry"
#   servicepackage = "containers"
#   api_version    = "2025-11-01"
#   data_source    = true
#   # path         = "internal/services/containers/registry"
# }
#
# resource "load_test" {
#   arm_type       = "Microsoft.LoadTestService/loadTests"
#   go_name        = "LoadTest"
#   servicepackage = "loadtestservice"
#   api_version    = "2022-12-01"
#   # service          = "LoadTestService"
#   # pandora_resource = "LoadTests"
# }
#
# resource "storage_mover_agent" {
#   arm_type       = "Microsoft.StorageMover/agents"
#   go_name        = "StorageMoverAgent"
#   servicepackage = "storagemover"
#   gen_resource   = false
#   list           = true
# }
