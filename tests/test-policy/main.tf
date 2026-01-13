provider "azurerm" {
  features {}
}

resource "azurerm_policy_definition" "test" {
  name         = "Acctestpol-Lucia"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "Acctestpol-Lucia"

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

resource "azurerm_management_group" "test" {
  display_name = "acctestmg-lucia"
}

resource "azurerm_management_group_policy_definition" "test2" {
  name                = "acctestpol-lucia-mg"
  policy_type         = "Custom"
  mode                = "Indexed"
  display_name        = "acctestpol-lucia-mg"
  management_group_id = azurerm_management_group.test.id

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

resource "azurerm_policy_definition" "test3" {
  name         = "acctestpol-lucia-test-after-migrate"
  policy_type  = "Custom"
  mode         = "All"
  display_name = "acctestpol-lucia-test-after-migrate"

  policy_rule = <<POLICY_RULE
    {
    "if": {
      "allOf": [
        {
          "not": {
            "field": "location",
            "in": "[parameters('allowedLocations')]"
          }
        },
        {
          "field": "location",
          "like": "[parameters('testObject').location]"
        }
      ]
    },
    "then": {
       "effect": "AuditIfNotExists",
        "details": {
          "type": "Microsoft.Insights/diagnosticSettings",
          "existenceCondition": {
            "allOf": [
            {
              "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.enabled",
              "equals": "true"
            },
            {
              "field": "Microsoft.Insights/diagnosticSettings/logs[*].retentionPolicy.days",
              "equals": "[parameters('requiredRetentionDays')]"
            }
          ]
        }
      }
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
    },
    "requiredRetentionDays": {
        "type": "Integer",
        "defaultValue": 365,
        "allowedValues": [
          0,
          30,
          90,
          180,
          365
        ],
        "metadata": {
          "displayName": "Required retention (days)",
          "description": "The required diagnostic logs retention in days"
      }
    },
    "testObject": {
      "type": "Object",
      "defaultValue": {
        "location": "westeurope"
      },
      "metadata": {
        "description": "test",
        "displayName": "test"
      },
      "schema": {
        "description": "test",
        "type": "object",
        "properties": {
          "location": {
            "description": "test",
            "type": "string"
          }
        },
        "additionalProperties": false
      }
    }
  }
PARAMETERS
}