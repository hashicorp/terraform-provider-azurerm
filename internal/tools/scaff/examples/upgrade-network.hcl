# Batch input for the scaff "upgrade" command — network resources for List:
#
#   go run ./internal/tools/scaff upgrade -input="internal/tools/scaff/examples/upgrade-network.hcl"
#
# Every resource below already has Resource Identity, so upgrading for List
# refactors Read into a reusable flatten method (where missing), generates the
# {name}_resource_list.go (+ test), and registers it in the service package's
# registration.go ListResources() method.
#
# NOTE — support status (the batch continues past any that fail and prints a
# summary at the end):
#
#   * Both [typed] and [untyped] resources are supported. All but one of these
#     are [untyped] native Plugin SDK resources; azurerm_virtual_hub_routing_intent
#     is [typed]. All already have Resource Identity.
#
#   * subnet and subnet_service_endpoint_storage_policy have their own Pandora
#     resource keys ("Subnets" / "ServiceEndpointPolicies") and upgrade cleanly.
#
#   * The virtual_hub_* CHILD resources (bgp_connection, connection, ip,
#     route_table) are parent-scoped: they generate a single parent-scoped list
#     call, e.g. client.VirtualHubIPConfigurationListComplete(ctx, *virtualHubID).
#     Because the go-azure-sdk `virtualwans` package is one Pandora "VirtualWANs"
#     resource, their read model is supplied via `read_model` and derived entirely
#     from source (no Pandora call). `list_method` overrides the derived SDK list
#     method name where it pluralises irregularly (e.g. bgp_connection).
#
#   * azurerm_virtual_hub itself is a TOP-LEVEL resource in the same package; its
#     subscription/resource-group list operations are not yet source-derived, so
#     Pandora resolution of "VirtualWANs" may produce an incorrect list for it.
#
# A running Pandora Data API is only required for resources that resolve list
# operations / read models via Pandora (i.e. not the parent-scoped children).

pandora_url = "http://localhost:8080"
write       = true # set true to APPLY the changes; false performs a dry run
overwrite   = false # set true to replace an existing generated *_resource_list.go

# azurerm_subnet [untyped] — child of virtual networks (Microsoft.Network/virtualNetworks/subnets).
resource "subnet" {
  file        = "internal/services/network/subnet_resource.go"
  arm_type    = "Microsoft.Network/virtualNetworks/subnets"
  api_version = "2025-01-01"
  list        = true
}

# azurerm_subnet_service_endpoint_storage_policy [untyped] — a service endpoint policy
# (Microsoft.Network/serviceEndpointPolicies), despite the terraform name.
resource "subnet_service_endpoint_storage_policy" {
  file        = "internal/services/network/subnet_service_endpoint_storage_policy_resource.go"
  arm_type    = "Microsoft.Network/serviceEndpointPolicies"
  api_version = "2025-01-01"
  list        = true
}

# azurerm_virtual_hub [untyped] — Microsoft.Network/virtualHubs (virtualwans SDK package).
resource "virtual_hub" {
  file             = "internal/services/network/virtual_hub_resource.go"
  service          = "Network"
  pandora_resource = "VirtualWANs"
  arm_type         = "Microsoft.Network/virtualHubs"
  api_version      = "2025-01-01"
  list             = true
}

# azurerm_virtual_hub_bgp_connection [untyped] — Microsoft.Network/virtualHubs/bgpConnections.
# Parent-scoped child of virtual_hub. The read model is derived from the vendored
# SDK Get method; list_method overrides the SDK List method name because it
# pluralises irregularly (VirtualHubBgpConnectionGet -> VirtualHubBgpConnectionsList).
resource "virtual_hub_bgp_connection" {
  file        = "internal/services/network/virtual_hub_bgp_connection_resource.go"
  list_method = "VirtualHubBgpConnectionsList"
  list        = true
}

# azurerm_virtual_hub_connection [untyped] — Microsoft.Network/virtualHubs/hubVirtualNetworkConnections.
# Parent-scoped child of virtual_hub (read model + list method derived from source).
resource "virtual_hub_connection" {
  file = "internal/services/network/virtual_hub_connection_resource.go"
  list = true
}

# azurerm_virtual_hub_ip [untyped] — Microsoft.Network/virtualHubs/ipConfigurations.
# Parent-scoped child of virtual_hub.
resource "virtual_hub_ip" {
  file = "internal/services/network/virtual_hub_ip_resource.go"
  list = true
}

# azurerm_virtual_hub_route_table [untyped] — Microsoft.Network/virtualHubs/hubRouteTables.
# Parent-scoped child of virtual_hub.
resource "virtual_hub_route_table" {
  file = "internal/services/network/virtual_hub_route_table_resource.go"
  list = true
}

# azurerm_virtual_hub_routing_intent [typed] — Microsoft.Network/virtualHubs/routingIntents.
# Parent-scoped child of virtual_hub. Typed resources are detected the same way as
# untyped (the parent ID that constructs the child ID), so this is source-derived.
resource "virtual_hub_routing_intent" {
  file = "internal/services/network/virtual_hub_routing_intent_resource.go"
  list = true
}
