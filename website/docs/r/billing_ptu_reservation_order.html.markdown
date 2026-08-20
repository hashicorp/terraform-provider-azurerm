---
subcategory: "Billing"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_billing_ptu_reservation_order"
description: |-
  Manages an Azure OpenAI PTU Reservation Order.
---

# azurerm_billing_ptu_reservation_order

Manages an Azure OpenAI PTU (Provisioned Throughput Unit) Reservation Order.

PTU reservations are billing commitments at the Azure subscription level that grant access to a fixed amount of provisioned throughput capacity for Azure OpenAI models. A single reservation order applies to all matching deployments within the specified scope.

~> **Note:** PTU reservations with `billing_plan = "Upfront"` are non-refundable after the return window (typically 30 days). Use `lifecycle { prevent_destroy = true }` to guard against accidental deletion.

~> **Note:** All arguments force a new resource to be created because Azure does not support in-place modification of reservation order properties.

## Example Usage

```hcl
data "azurerm_subscription" "current" {}

resource "azurerm_billing_ptu_reservation_order" "example" {
  name              = "example-ptu-reservation"
  location          = "eastus"
  capacity          = 100
  billing_scope_id  = data.azurerm_subscription.current.id
  sku_name          = "DataZoneProvisionedManaged"
  term              = "P1Y"
  billing_plan      = "Monthly"
  applied_scope_type = "Shared"
  renew             = false
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The display name for the PTU reservation order. Changing this forces a new resource to be created.

* `location` - (Required) The Azure region where the PTU capacity is reserved (e.g. `eastus`). Changing this forces a new resource to be created.

* `capacity` - (Required) The number of PTUs (Provisioned Throughput Units) to reserve. Must be at least `1`. Changing this forces a new resource to be created.

* `billing_scope_id` - (Required) The billing scope for the reservation, for example `/subscriptions/00000000-0000-0000-0000-000000000000`. Changing this forces a new resource to be created.

---

* `sku_name` - (Optional) The SKU name for the PTU reservation. Defaults to `DataZoneProvisionedManaged`. Changing this forces a new resource to be created.

* `term` - (Optional) The reservation term expressed as an ISO 8601 duration. Possible values are `P1M`, `P1Y`, and `P3Y`. Defaults to `P1Y`. Changing this forces a new resource to be created.

* `billing_plan` - (Optional) The billing plan for the reservation. Possible values are `Upfront` and `Monthly`. Defaults to `Upfront`. Changing this forces a new resource to be created.

* `applied_scope_type` - (Optional) The scope type to which the reservation benefit is applied. Possible values are `Shared`, `Single`, and `ManagementGroup`. Defaults to `Shared`. Changing this forces a new resource to be created.

* `renew` - (Optional) Whether the reservation automatically renews at the end of the term. Defaults to `true`. Changing this forces a new resource to be created.

## Attributes Reference

In addition to the Arguments listed above, the following Attributes are exported:

* `id` - The ARM resource ID of the reservation order, in the format `/providers/Microsoft.Capacity/reservationOrders/<uuid>`.

* `order_id` - The UUID of the reservation order in Azure.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the PTU Reservation Order.
* `read` - (Defaults to 5 minutes) Used when retrieving the PTU Reservation Order.
* `update` - (Defaults to 30 minutes) Used when updating the PTU Reservation Order.
* `delete` - (Defaults to 30 minutes) Used when deleting the PTU Reservation Order.

## Import

PTU Reservation Orders can be imported using the `id`, e.g.

```shell
terraform import azurerm_billing_ptu_reservation_order.example /providers/Microsoft.Capacity/reservationOrders/00000000-0000-0000-0000-000000000000
```
