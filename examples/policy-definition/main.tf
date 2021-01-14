provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "example" {
  name         = "${var.prefix}-policy"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "${var.prefix}-policy"

  policy_rule = <<POLICY_RULE
	{
    "if": {
      "not": {
        "field": "location",
        "in": "[parameters('allowedLocations')]"
      }
    },
    "then": {
      "effect": "audit"
    }
  }
	POLICY_RULE

  parameters = <<PARAMETERS
	{
    "allowedLocations": {
      "type": "Array",
      "metadata": {
        "description": "The list of allowed locations for resources.",
        "displayName": "Allowed locations",
        "strongType": "location"
      }
    }
  }
  	PARAMETERS
}
