---
subcategory: "Database"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_mssql_execute_job"
description: |-
  Executes an Elastic Job
---

# Action: azurerm_mssql_execute_job

~> **Note:** `azurerm_mssql_execute_job` is in beta. Its interface and behaviour may change as the feature evolves, and breaking changes are possible. It is offered as a technical preview without compatibility guarantees until Terraform 1.14 is generally available.

Executes an Elastic Job

## Example Usage

```terraform
resource "azurerm_mssql_job" "example" {
  # ... Job configuration
}

resource "azurerm_mssql_job_step" "example" {
  # ... Job Step configuration
}

resource "terraform_data" "example" {
  input = azurerm_mssql_job_step.example.id

  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azurerm_mssql_execute_job.example]
    }
  }
}

action "azurerm_mssql_execute_job" "example" {
  config {
    job_id = azurerm_mssql_job.test.id
  }
}

```

## Argument Reference

This action supports the following arguments:

* `job_id` - (Required) The ID of the job to execute.

---

* `timeout` - (Optional) Timeout duration for the action to complete. Defaults to `15m`.

* `wait_for_completion` - (Optional) Whether to poll the job execution for completion. Defaults to `false`.
