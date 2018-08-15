resource "azurerm_policy_definition" "policy" {
  name         = "${var.policy_definition_name}"
  policy_type  = "${var.policy_type}"
  mode         = "${var.mode}"
  display_name = "${var.display_name}"

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
