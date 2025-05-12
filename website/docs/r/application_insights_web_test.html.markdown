---
subcategory: "Application Insights"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_application_insights_web_test"
description: |-
  Manages an Application Insights WebTest.
---

# azurerm_application_insights_web_test

Manages an Application Insights WebTest.

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "tf-test"
  location = "West Europe"
}

resource "azurerm_application_insights" "example" {
  name                = "tf-test-appinsights"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  application_type    = "web"
}

resource "azurerm_application_insights_web_test" "example" {
  name                    = "tf-test-appinsights-webtest"
  location                = azurerm_application_insights.example.location
  resource_group_name     = azurerm_resource_group.example.name
  application_insights_id = azurerm_application_insights.example.id
  kind                    = "ping"
  frequency               = 300
  timeout                 = 60
  enabled                 = true
  geo_locations           = ["us-tx-sn1-azr", "us-il-ch1-azr"]

  configuration = <<XML
<WebTest Name="WebTest1" Id="ABD48585-0831-40CB-9069-682EA6BB3583" Enabled="True" CssProjectStructure="" CssIteration="" Timeout="0" WorkItemIds="" xmlns="http://microsoft.com/schemas/VisualStudio/TeamTest/2010" Description="" CredentialUserName="" CredentialPassword="" PreAuthenticate="True" Proxy="default" StopOnError="False" RecordedResultFile="" ResultsLocale="">
  <Items>
    <Request Method="GET" Guid="a5f10126-e4cd-570d-961c-cea43999a200" Version="1.1" Url="http://microsoft.com" ThinkTime="0" Timeout="300" ParseDependentRequests="True" FollowRedirects="True" RecordResult="True" Cache="False" ResponseTimeGoal="0" Encoding="utf-8" ExpectedHttpStatusCode="200" ExpectedResponseUrl="" ReportingName="" IgnoreHttpStatusCode="False" />
  </Items>
</WebTest>
XML

}

output "webtest_id" {
  value = azurerm_application_insights_web_test.example.id
}

output "webtests_synthetic_id" {
  value = azurerm_application_insights_web_test.example.synthetic_monitor_id
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Application Insights WebTest. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Application Insights WebTest. Changing this forces a new resource

* `application_insights_id` - (Required) The ID of the Application Insights component on which the WebTest operates. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created. It needs to correlate with location of parent resource (azurerm_application_insights).

* `kind` - (Required) The kind of web test that this web test watches. Choices are `ping` and `multistep`. Changing this forces a new resource to be created.

* `geo_locations` - (Required) A list of where to physically run the tests from to give global coverage for accessibility of your application.

~> **Note:** [Valid options for geo locations are described here](https://docs.microsoft.com/azure/azure-monitor/app/monitor-web-app-availability#location-population-tags)

* `configuration` - (Required) An XML configuration specification for a WebTest ([see here for more information](https://docs.microsoft.com/rest/api/application-insights/webtests/createorupdate/)).

* `frequency` - (Optional) Interval in seconds between test runs for this WebTest. Valid options are `300`, `600` and `900`. Defaults to `300`.

* `timeout` - (Optional) Seconds until this WebTest will timeout and fail. Default is `30`.

* `enabled` - (Optional) Is the test actively being monitored.

* `retry_enabled` - (Optional) Allow for retries should this WebTest fail.

* `description` - (Optional) Purpose/user defined descriptive test for this WebTest.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Application Insights Web Test.
* `read` - (Defaults to 5 minutes) Used when retrieving the Application Insights Web Test.
* `update` - (Defaults to 30 minutes) Used when updating the Application Insights Web Test.
* `delete` - (Defaults to 30 minutes) Used when deleting the Application Insights Web Test.

## Import

Application Insights Web Tests can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_application_insights_web_test.my_test /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Insights/webTests/my_test
```
