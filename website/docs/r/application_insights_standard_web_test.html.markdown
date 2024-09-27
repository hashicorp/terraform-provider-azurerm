---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_standard_web_test"
description: |-
  Manages a Application Insights Standard WebTest.
---

# azurerm_application_insights_standard_web_test

Manages a Application Insights Standard WebTest.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "rg-example"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "example"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_standard_web_test" "example" {
  name                    = "example-test"
  resource_group_name     = azurerm_resource_group.example.name
  location                = "West Europe"
  application_insights_id = azurerm_application_insights.example.id
  geo_locations           = ["example"]

  request {
    url = "http://www.example.com"
  }
}
```

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Application Insights Standard WebTest. Changing this forces a new Application Insights Standard WebTest to be created.

* `resource_group_name` - (Required) The name of the Resource Group where the Application Insights Standard WebTest should exist. Changing this forces a new Application Insights Standard WebTest to be created.

* `location` - (Required) The Azure Region where the Application Insights Standard WebTest should exist. Changing this forces a new Application Insights Standard WebTest to be created. It needs to correlate with location of the parent resource (azurerm_application_insights)

* `application_insights_id` - (Required) The ID of the Application Insights instance on which the WebTest operates. Changing this forces a new Application Insights Standard WebTest to be created.

* `geo_locations` - (Required) Specifies a list of where to physically run the tests from to give global coverage for accessibility of your application.

~> **Note:** [Valid options for geo locations are described here](https://docs.microsoft.com/azure/azure-monitor/app/monitor-web-app-availability#location-population-tags)

* `request` - (Required) A `request` block as defined below.


---

* `description` - (Optional) Purpose/user defined descriptive test for this WebTest.

* `enabled` - (Optional) Should the WebTest be enabled?

* `frequency` - (Optional) Interval in seconds between test runs for this WebTest. Valid options are `300`, `600` and `900`. Defaults to `300`.

* `retry_enabled` - (Optional) Should the retry on WebTest failure be enabled?

* `tags` - (Optional) A mapping of tags which should be assigned to the Application Insights Standard WebTest.

* `timeout` - (Optional) Seconds until this WebTest will timeout and fail. Default is `30`.

* `validation_rules` - (Optional) A `validation_rules` block as defined below.

---

A `content` block supports the following:

* `content_match` - (Required) A string value containing the content to match on.

* `ignore_case` - (Optional) Ignore the casing in the `content_match` value.

* `pass_if_text_found` - (Optional) If the content of `content_match` is found, pass the test. If set to `false`, the WebTest is failing if the content of `content_match` is found.

---

A `header` block supports the following:

* `name` - (Required) The name which should be used for a header in the request.

* `value` - (Required) The value which should be used for a header in the request.

---

A `request` block supports the following:

* `url` - (Required) The WebTest request URL.

* `body` - (Optional) The WebTest request body.

* `follow_redirects_enabled` - (Optional) Should the following of redirects be enabled? Defaults to `true`.

* `header` - (Optional) One or more `header` blocks as defined above.

* `http_verb` - (Optional) Which HTTP verb to use for the call. Options are 'GET', 'POST', 'PUT', 'PATCH', 'DELETE', 'HEAD', and 'OPTIONS'. Defaults to `GET`.

* `parse_dependent_requests_enabled` - (Optional) Should the parsing of dependend requests be enabled? Defaults to `true`.

---

A `validation_rules` block supports the following:

* `content` - (Optional) A `content` block as defined above.

* `expected_status_code` - (Optional) The expected status code of the response. Default is '200', '0' means 'response code < 400'

* `ssl_cert_remaining_lifetime` - (Optional) The number of days of SSL certificate validity remaining for the checked endpoint. If the certificate has a shorter remaining lifetime left, the test will fail. This number should be between 1 and 365.

* `ssl_check_enabled` - (Optional) Should the SSL check be enabled?

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Application Insights Standard WebTest.

* `synthetic_monitor_id` - Unique ID of this WebTest. This is typically the same value as the Name field.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Standard WebTest.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Standard WebTest.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Standard WebTest.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Standard WebTest.

## Import

Application Insights Standard WebTests can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_standard_web_test.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Insights/webTests/appinsightswebtest
```
