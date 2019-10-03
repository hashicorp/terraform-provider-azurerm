---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_notification_hub_authorization_rule"
sidebar_current: "docs-azurerm-resource-messaging-notification-hub-authorization-rule"
description: |-
  Manages an Authorization Rule associated with a Notification Hub within a Notification Hub Namespace.

---

# azurerm_notification_hub_authorization_rule

Manages an Authorization Rule associated with a Notification Hub within a Notification Hub Namespace.

## Example Usage

```hcl
resource "azurerm_resource_group" "test" {
  name     = "notificationhub-resources"
  location = "Australia East"
}

resource "azurerm_notification_hub_namespace" "test" {
  name                = "myappnamespace"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  namespace_type      = "NotificationHub"

  sku_name = "Free"
}

resource "azurerm_notification_hub" "test" {
  name                = "mynotificationhub"
  namespace_name      = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
}

resource "azurerm_notification_hub_authorization_rule" "test" {
  name                  = "management-auth-rule"
  notification_hub_name = "${azurerm_notification_hub.test.name}"
  namespace_name        = "${azurerm_notification_hub_namespace.test.name}"
  resource_group_name   = "${azurerm_resource_group.test.name}"
  manage                = true
  send                  = true
  listen                = true
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name to use for this Authorization Rule. Changing this forces a new resource to be created.

* `notification_hub_name` - (Required) The name of the Notification Hub for which the Authorization Rule should be created. Changing this forces a new resource to be created.

* `namespace_name` - (Required) The name of the Notification Hub Namespace in which the Notification Hub exists. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the Resource Group in which the Notification Hub Namespace exists. Changing this forces a new resource to be created.

* `manage` - (Optional) Does this Authorization Rule have Manage access to the Notification Hub? Defaults to `false`.

-> **NOTE:** If `manage` is set to `true` then both `send` and `listen` must also be set to `true`.

* `send` - (Optional) Does this Authorization Rule have Send access to the Notification Hub? Defaults to `false`.

* `listen` - (Optional) Does this Authorization Rule have Listen access to the Notification Hub? Defaults to `false`.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the Authorization Rule.

* `primary_access_key` - The Primary Access Key associated with this Authorization Rule.

* `secondary_access_key` - The Secondary Access Key associated with this Authorization Rule.

## Import

Notification Hub Authorization Rule can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_notification_hub_authorization_rule.rule1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.NotificationHubs/namespaces/{namespaceName}/notificationHubs/hub1/AuthorizationRules/rule1
```
