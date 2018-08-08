data "azurerm_subscription" "example" {}

data "azurerm_client_config" "example" {}

resource "azurerm_role_definition" "example" {
  name  = "example-role-definition"
  scope = "${data.azurerm_subscription.example.id}"

  permissions {
    actions     = ["Microsoft.Resources/subscriptions/resourceGroups/read"]
    not_actions = []
  }

  assignable_scopes = [
    "${data.azurerm_subscription.example.id}",
  ]
}

resource "azurerm_role_assignment" "example" {
  scope              = "${data.azurerm_subscription.example.id}"
  role_definition_id = "${azurerm_role_definition.example.id}"
  principal_id       = "${data.azurerm_client_config.example.service_principal_object_id}"
}
