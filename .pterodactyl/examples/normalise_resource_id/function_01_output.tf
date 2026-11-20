# result: /subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.ApiManagement/service/service1/gateways/gateway1/hostnameConfigurations/config1

output "test" {
  value = provider::azurerm::normalise_resource_id("/Subscriptions/12345678-1234-9876-4563-123456789012/ResourceGroups/resGroup1/PROVIDERS/microsoft.apimanagement/service/service1/gateWays/gateway1/hostnameconfigurations/config1")
}
