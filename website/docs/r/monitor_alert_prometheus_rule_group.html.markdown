---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_alert_prometheus_rule_group"
description: |-
  Manages an Alert Management Prometheus Rule Group.
---

# azurerm_monitor_alert_prometheus_rule_group

Manages an Alert Management Prometheus Rule Group.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_monitor_action_group" "example" {
  name                = "example-mag"
  resource_group_name = azurerm_resource_group.example.name
  short_name          = "testag"
}

resource "azurerm_monitor_workspace" "example" {
  name                = "example-amw"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "example-cluster"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  dns_prefix          = "example-aks"

  default_node_pool {
    name                   = "default"
    node_count             = 1
    vm_size                = "Standard_DS2_v2"
    enable_host_encryption = true
  }

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_monitor_alert_prometheus_rule_group" "example" {
  name                = "example-amprg"
  location            = "West Europe"
  resource_group_name = azurerm_resource_group.example.name
  cluster_name        = azurerm_kubernetes_cluster.example.name
  description         = "This is the description of the following rule group"
  rule_group_enabled  = false
  interval            = "PT1M"
  scopes              = [azurerm_monitor_workspace.example.id]
  rule {
    enabled    = false
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    record     = "job_type:billing_jobs_duration_seconds:99p5m"
    labels = {
      team = "prod"
    }
  }

  rule {
    alert      = "Billing_Processing_Very_Slow"
    enabled    = true
    expression = <<EOF
histogram_quantile(0.99, sum(rate(jobs_duration_seconds_bucket{service="billing-processing"}[5m])) by (job_type))
EOF
    for        = "PT5M"
    severity   = 2

    action {
      action_group_id = azurerm_monitor_action_group.example.id
    }

    alert_resolution {
      auto_resolved   = true
      time_to_resolve = "PT10M"
    }

    annotations = {
      annotationName = "annotationValue"
    }

    labels = {
      team = "prod"
    }
  }
  tags = {
    key = "value"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) Specifies the name which should be used for this Alert Management Prometheus Rule Group. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the Azure Region where the Alert Management Prometheus Rule Group should exist. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) Specifies the name of the Resource Group where the Alert Management Prometheus Rule Group should exist. Changing this forces a new resource to be created.

* `rule` - (Required) A `rule` block as defined below.

* `scopes` - (Required) Specifies the resource ID of the Azure Monitor Workspace.

* `cluster_name` - (Optional) Specifies the name of the Managed Kubernetes Cluster.

* `description` - (Optional) The description of the Alert Management Prometheus Rule Group.

* `rule_group_enabled` - (Optional) Is this Alert Management Prometheus Rule Group enabled? Possible values are `true` and `false`.

* `interval` - (Optional) Specifies the interval in which to run the Alert Management Prometheus Rule Group represented in ISO 8601 duration format. Possible values are between `PT1M` and `PT15M`.

* `tags` - (Optional) A mapping of tags to assign to the Alert Management Prometheus Rule Group.

---

A `rule` block supports the following:

* `action` - (Optional) An `action` block as defined below.

* `alert` - (Optional) Specifies the Alert rule name.

* `annotations` - (Optional) Specifies a set of informational labels that can be used to store longer additional information such as alert descriptions or runbook links.

* `enabled` - (Optional) Is this rule enabled? Possible values are `true` and `false`.

* `expression` - (Required) Specifies the Prometheus Query Language expression to evaluate. For more details see [this doc](https://prometheus.io/docs/prometheus/latest/querying/basics). Evaluate at the period given by `interval` and record the result as a new set of time series with the metric name given by `record`. 

* `for` - (Optional) Specifies the amount of time alert must be active before firing, represented in ISO 8601 duration format.

* `labels` - (Optional) Specifies the labels to add or overwrite before storing the result.

* `record` - (Optional) Specifies the recorded metrics name.

* `alert_resolution` - (Optional) An `alert_resolution` block as defined below.

* `severity` - (Optional) Specifies the severity of the alerts fired by the rule. Possible values are between 0 and 4.

---

An `action` block supports the following:

* `action_group_id` - (Required) Specifies the resource id of the monitor action group.

* `action_properties` - (Optional) Specifies the properties of an action group object.
 
-> **Note:** `action_properties` can only be configured for IcM Connector Action Groups for now. Other public features will be supported in the future.

---

An `alert_resolution` block supports the following:

* `auto_resolved` - (Optional) Is the alert auto-resolution? Possible values are `true` and `false`.

* `time_to_resolve` - (Optional) Specifies the alert auto-resolution interval, represented in ISO 8601 duration format.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Alert Management Prometheus Rule Group.


## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Alert Management Prometheus Rule Group.
* `read` - (Defaults to 5 minutes) Used when retrieving the Alert Management Prometheus Rule Group.
* `update` - (Defaults to 30 minutes) Used when updating the Alert Management Prometheus Rule Group.
* `delete` - (Defaults to 30 minutes) Used when deleting the Alert Management Prometheus Rule Group.

## Import

Alert Management Prometheus Rule Group can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_alert_prometheus_rule_group.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.AlertsManagement/prometheusRuleGroups/ruleGroup1
```
