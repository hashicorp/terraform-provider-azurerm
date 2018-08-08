data "azurerm_subscription" "example" {}

resource "azurerm_management_lock" "example" {
  name       = "subscription-level"
  scope      = "${data.azurerm_subscription.example.id}"
  lock_level = "CanNotDelete"
  notes      = "Items can't be deleted in this subscription!"
}
