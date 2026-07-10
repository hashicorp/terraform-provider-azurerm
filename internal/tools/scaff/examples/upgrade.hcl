# Example input file for the scaff "upgrade" command:
#
#   go run ./internal/tools/scaff upgrade -input="internal/tools/scaff/examples/upgrade.hcl"
#
# The "upgrade" command makes an EXISTING typed (internal/sdk) resource
# list-ready: it adds Resource Identity and refactors Read into a reusable
# flatten method when those are missing, then optionally generates the list
# resource (and its acceptance test) alongside it.
#
# Each "resource" block targets one existing resource file. The block label is
# the terraform resource name (snake_case, without the provider prefix, e.g.
# "monitor_workspace" -> azurerm_monitor_workspace) and is used for the
# resource-identity test go:generate directive.
#
# A running Pandora Data API is required whenever a block resolves list
# operations or the SDK read model (i.e. list = true, or flatten with no
# read_model). See -pandora-url / pandora_url below.
#
# IMPORTANT: this command MODIFIES hand-written resource files. It only writes
# when write = true; the default below is a dry run that prints a unified diff of
# the proposed changes and previews any generated files. Review the diff, then
# set write = true (or pass -input with the file's write flag) to apply.

# ---------------------------------------------------------------------------
# Optional global overrides. When omitted these fall back to .scaff.hcl (or the
# tool's built-in defaults), so you usually don't need to set them here.
# ---------------------------------------------------------------------------
# pandora_url = "http://localhost:8080" # Pandora Data API base URL
# provider    = "azurerm"               # provider name (terraform type prefix + import paths)
# org         = "hashicorp"             # provider GitHub org (import paths)
write     = true # set true to APPLY the changes; false performs a dry run
overwrite = true # set true to replace an existing generated *_resource_list.go

# Azure Monitor Workspace (Microsoft.Monitor/accounts) — a typed resource that
# currently has neither Resource Identity nor a flatten method. With list = true
# the command adds both and generates monitor_workspace_resource_list.go.
#
# The ARM segment ("accounts") does not match the Pandora resource key, so the
# resource is addressed explicitly via service + pandora_resource. The API
# version is pinned because the resource key differs in newer versions.
resource "monitor_workspace" {
  file             = "internal/services/monitor/monitor_workspace_resource.go"
  service          = "Monitor"
  pandora_resource = "azuremonitorworkspaces"
  api_version      = "2023-04-03"
  list             = true
}

# Add ONLY Resource Identity to an existing resource (no list). No Pandora call
# is needed for identity alone — the ID type is read from the resource's own
# Read/Importer. Uncomment and point "file" at a real resource to try it.
#
# resource "storage_mover" {
#   file     = "internal/services/storagemover/storage_mover_resource.go"
#   identity = true
# }

# Refactor Read into a reusable flatten method only. The SDK read model is
# resolved from Pandora via arm_type (or supply read_model to skip the call).
#
# resource "example_resource" {
#   file     = "internal/services/example/example_resource.go"
#   arm_type = "Microsoft.Example/examples"
#   flatten  = true
#   # read_model = "ExampleResourceModel"
#   # identity_properties = "name,resource_group_name" # for the identity test directive
# }
