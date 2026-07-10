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
overwrite   = false                   # set true to replace existing generated files

# RedHat OpenShift cluster (Microsoft.RedHatOpenShift) — resource + list resource
# + data source, pinned to a specific API version. go_name is set because the
# snake->camel derivation ("RedhatOpenshiftCluster") wouldn't match the canonical
# "RedHatOpenShift" casing.
resource "redhat_openshift_cluster" {
  arm_type       = "Microsoft.RedHatOpenShift/openShiftClusters"
  go_name        = "RedHatOpenShiftCluster"
  servicepackage = "redhatopenshift"
  api_version    = "2025-07-25"
  list           = true
  data_source    = true
}

# Storage Mover (Microsoft.StorageMover) — minimal block. api_version defaults to
# the latest non-preview version, and go_name is derived from the name
# ("storage_mover" -> "StorageMover").
resource "storage_mover" {
  arm_type       = "Microsoft.StorageMover/storageMovers"
  servicepackage = "storagemover"
}

# NetApp account (Microsoft.NetApp) — resource + data source. go_name overridden
# to preserve the "NetApp" casing.
resource "netapp_account" {
  arm_type       = "Microsoft.NetApp/netAppAccounts"
  go_name        = "NetAppAccount"
  servicepackage = "netapp"
  data_source    = true
}

# Automation account (Microsoft.Automation) — resource + list resource.
resource "automation_account" {
  arm_type       = "Microsoft.Automation/automationAccounts"
  go_name        = "AutomationAccount"
  servicepackage = "automation"
  api_version    = "2024-10-23"
  list           = true
}

# Container Registry (Microsoft.ContainerRegistry) — the Pandora service name
# (ContainerRegistry) differs from the provider's service package directory
# (containers), so servicepackage places the files in internal/services/containers.
# An explicit output directory override is shown (commented) via "path".
resource "container_registry" {
  arm_type       = "Microsoft.ContainerRegistry/registries"
  go_name        = "ContainerRegistry"
  servicepackage = "containers"
  api_version    = "2025-11-01"
  data_source    = true
  # path         = "internal/services/containers/registry"
}

# Load Test (Microsoft.LoadTestService) — addressed here by arm_type, but a
# resource can also be addressed by its explicit Pandora service + resource key
# instead (either form is accepted). The equivalent is shown commented below.
resource "load_test" {
  arm_type       = "Microsoft.LoadTestService/loadTests"
  go_name        = "LoadTest"
  servicepackage = "loadtestservice"
  api_version    = "2022-12-01"
  # service          = "LoadTestService"
  # pandora_resource = "LoadTests"
}

# List-only regeneration — generate ONLY the list resource for a resource that
# already exists (assumes the resource file already provides a flatten method).
resource "storage_mover_agent" {
  arm_type       = "Microsoft.StorageMover/agents"
  go_name        = "StorageMoverAgent"
  servicepackage = "storagemover"
  gen_resource   = false
  list           = true
}
