# Example 1: Create a policy definition
resource "azurerm_policy_definition" "test" {
  name         = "test-definition-version"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "Test Definition Version Policy"

  metadata = jsonencode({
    version = "1.0.0"
  })

  policy_rule = jsonencode({
    if = {
      not = {
        field  = "location"
        equals = "westeurope"
      }
    }
    then = {
      effect = "Audit"
    }
  })
}

# Example 2: Test policy assignment with definition_version
data "azurerm_subscription" "current" {}

resource "azurerm_subscription_policy_assignment" "test" {
  name                 = "test-def-version"
  policy_definition_id = azurerm_policy_definition.test.id
  subscription_id      = data.azurerm_subscription.current.id
  
  # NEW FEATURE: Pin to specific version format
  definition_version = "1.0.*"  # Pin to version 1.0.x (any patch)
  
  # Alternatively, you can use:
  # definition_version = "1.*.*"   # Pin to version 1.x.x (any minor/patch)
  # definition_version = "1.0.0"   # Pin to exact version 1.0.0
}

# Example 3: Test data source reading definition_version
data "azurerm_policy_assignment" "test" {
  name     = azurerm_subscription_policy_assignment.test.name
  scope_id = data.azurerm_subscription.current.id
}

# Output to verify the definition_version is being read
output "assignment_definition_version" {
  value = data.azurerm_policy_assignment.test.definition_version
}
