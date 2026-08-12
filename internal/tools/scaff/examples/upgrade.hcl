# Example input file for the scaff "upgrade" command:
#
#   go run ./internal/tools/scaff upgrade -input="internal/tools/scaff/examples/upgrade.hcl"
#
# ---------------------------------------------------------------------------
# Optional global overrides. When omitted these fall back to .scaff.hcl (or the
# tool's built-in defaults), so you usually don't need to set them here.
# ---------------------------------------------------------------------------

pandora_url = "http://localhost:8080"
write       = true # set true to APPLY the changes; false performs a dry run
overwrite   = true # set true to replace an existing generated *_resource_list.go

resource "eventgrid_domain" {
  file             = "/Users/mark/gitRepos/scaff-workspace/terraform-provider-azurerm/internal/services/eventgrid/eventgrid_domain_resource.go"
}
# azurerm_subnet [untyped] — parent-scoped child of virtual network
# (Microsoft.Network/virtualNetworks/subnets); its list is derived from the
# vendored SDK (subnets.ListComplete(ctx, commonids.VirtualNetworkId)).
resource "subnet" {
  file = "internal/services/network/subnet_resource.go"
}

# azurerm_subnet_service_endpoint_storage_policy [untyped] — a top-level service
# endpoint policy (Microsoft.Network/serviceEndpointPolicies), despite the
# terraform name; its subscription/resource-group list methods are SDK-derived.
resource "subnet_service_endpoint_storage_policy" {
  file = "internal/services/network/subnet_service_endpoint_storage_policy_resource.go"
}

# azurerm_virtual_hub [untyped] — Microsoft.Network/virtualHubs (virtualwans SDK package).
# TOP-LEVEL resource: its read model and subscription/resource-group list methods
# are derived from the vendored SDK, so no Pandora attributes are required.
resource "virtual_hub" {
  file = "internal/services/network/virtual_hub_resource.go"
}

# azurerm_virtual_hub_bgp_connection [untyped] — Microsoft.Network/virtualHubs/bgpConnections.
# Parent-scoped child of virtual_hub. The read model is derived from the vendored
# SDK Get method; list_method overrides the SDK List method name because it
# pluralises irregularly (VirtualHubBgpConnectionGet -> VirtualHubBgpConnectionsList).
resource "virtual_hub_bgp_connection" {
  file        = "internal/services/network/virtual_hub_bgp_connection_resource.go"
}

# azurerm_virtual_hub_connection [untyped] — Microsoft.Network/virtualHubs/hubVirtualNetworkConnections.
# Parent-scoped child of virtual_hub (read model + list method derived from source).
resource "virtual_hub_connection" {
  file = "internal/services/network/virtual_hub_connection_resource.go"
}

# azurerm_virtual_hub_ip [untyped] — Microsoft.Network/virtualHubs/ipConfigurations.
# Parent-scoped child of virtual_hub.
resource "virtual_hub_ip" {
  file = "internal/services/network/virtual_hub_ip_resource.go"
}

# azurerm_virtual_hub_route_table [untyped] — Microsoft.Network/virtualHubs/hubRouteTables.
# Parent-scoped child of virtual_hub.
resource "virtual_hub_route_table" {
  file = "internal/services/network/virtual_hub_route_table_resource.go"
}

# azurerm_virtual_hub_routing_intent [typed] — Microsoft.Network/virtualHubs/routingIntents.
# Parent-scoped child of virtual_hub. Typed resources are detected the same way as
# untyped (the parent ID that constructs the child ID), so this is source-derived.
resource "virtual_hub_routing_intent" {
  file = "internal/services/network/virtual_hub_routing_intent_resource.go"
}
