---
layout: "azurerm"
page_title: "Azure Resource Manager: Continuous Validation with Terraform Cloud"
description: |-
Azure Resource Manager: Continuous Validation with Terraform Cloud

---

# Continuous Validation with Terraform Cloud

The Continuous Validation feature in Terraform Cloud (TFC) allows users to make assertions about their infrastructure between applied runs. This helps users to identify issues at the time they first appear and avoid situations where a change is only identified during a future terraform plan/apply or once it causes a user-facing problem.

Checks can be added to Terraform configuration in Terraform Cloud (TFC) using check blocks. Check blocks contain assertions that are defined with a custom condition expression and an error message. When the condition expression evaluates to false, Terraform will show a warning message that includes the user-defined error message.

Custom conditions can be created using data from Terraform providers’ resources and data sources. Data can also be combined from multiple sources; for example, you can use checks to monitor expirable resources by comparing a resource’s expiration date attribute to the current time returned by Terraform’s built-in time functions.

Below, this guide shows examples of how data returned from the AzureRM provider can be used to define checks in your Terraform configuration. For more information about continuous validation visit the [Workspace Health](https://developer.hashicorp.com/terraform/cloud-docs/workspaces/health#continuous-validation) page in the Terraform Cloud documentation.

## Example - Check if a VM's is not running (`azurerm_linux_virtual_machine`, `azurerm_windows_virtual_machine`)

VM instances provisioned can pass through several different power states as past of the VM instance lifecycle. Once a VM is provisioned it could experience an error, or a user could suspend or stop that VM, without that change being detected until the next Terraform plan is generated. Continuous validation can be used to assert the state of a VM and detect if there are any unexpected status changes that occur out-of-band.

The example below shows how a check block can be used to assert that a VM is in the running state. You can force the check to fail in this example by provisioning the VM and manually stopping it, and then triggering a health check in TFC. The check will fail and report that the VM is not running.

```hcl
data "azurerm_virtual_machine" "example" {
  name                = azurerm_linux_virtual_machine.example.name
  resource_group_name = azurerm_resource_group.example.name
}

check "check_vm_state" {
  assert {
    condition = data.azurerm_virtual_machine.example.power_state == "running"
    error_message = format("Virtual Machine (%s) should be in a 'running' status, instead state is '%s'",
      data.azurerm_virtual_machine.example.id,
      data.azurerm_virtual_machine.example.power_state
    )
  }
}
```

The full example can be found in the provider's [examples](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/tfc-checks/vm-power-state) folder. 

## Example - Check if a Container App certificate will expire within a certain timeframe (`azurerm_app_service_certificate`)

Azure App Service Certificates (and other resources) can be provisioned using a user supplied certificate. In the example below we show how to check that a certificate should be valid for the next 30 days (see `local.month_in_hour_duration`).

```hcl
locals {
  month_in_hour_duration            = "${24 * 30}h"
  month_and_2min_in_second_duration = "${(60 * 60 * 24 * 30) + (60 * 2)}s"
}

data "azurerm_app_service_certificate" "example" {
  name                = azurerm_app_service_certificate.example.name
  resource_group_name = azurerm_app_service_certificate.example.resource_group_name
}

check "check_certificate_state" {
  assert {
    condition = timecmp(plantimestamp(), timeadd(
      data.azurerm_app_service_certificate.example.expiration_date,
    "-${local.month_in_hour_duration}")) < 0
    error_message = format("App Service Certificate (%s) is valid for at least 30 days, but is due to expire on `%s`.",
      data.azurerm_app_service_certificate.example.id,
      data.azurerm_app_service_certificate.example.expiration_date
    )
  }
}

```

The full example can be found in the provider's [examples](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/tfc-checks/app-service-certificate-expiry) folder. 

## Example - Check if an App Service Function or Web App has exceeded its usage limit (`azurerm_linux_function_app`, `azurerm_linux_web_app`)

App Service Function and Web Apps can exceed their usage limit. The example below shows how a check block can be used to assert that a Function or Web App has not exceeded its usage limit.

```hcl
data "azurerm_virtual_machine" "example" {
  name                = azurerm_linux_virtual_machine.example.name
  resource_group_name = azurerm_resource_group.example.name
}

check "check_vm_state" {
  assert {
    condition = data.azurerm_virtual_machine.example.power_state == "running"
    error_message = format("Virtual Machine (%s) should be in a 'running' status, instead state is '%s'",
      data.azurerm_virtual_machine.example.id,
      data.azurerm_virtual_machine.example.power_state
    )
  }
}
```

The full example can be found in the provider's [examples](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/tfc-checks/app-service-app-usage) folder. 
