# Example input file for the scaff "upgrade" command:
#
#   go run ./internal/tools/scaff upgrade -input="internal/tools/scaff/examples/upgrade.hcl"
#
# ---------------------------------------------------------------------------
# Optional global overrides. When omitted these fall back to .scaff.hcl (or the
# tool's built-in defaults), so you usually don't need to set them here.
# ---------------------------------------------------------------------------
# pandora_url = "http://localhost:8080" # Pandora Data API base URL

write     = true # set true to APPLY the changes; false performs a dry run
overwrite = true # set true to replace an existing generated *_resource_list.go

resource "eventgrid_domain" {
  file             = "/Users/mark/gitRepos/scaff-workspace/terraform-provider-azurerm/internal/services/eventgrid/eventgrid_domain_resource.go"
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
#   flatten = true
#   # read_model = "ExampleResourceModel"
#   # identity_properties = "name,resource_group_name" # for the identity test directive
# }
