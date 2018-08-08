---
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_scheduler_job"
sidebar_current: "docs-azurerm-resource-scheduler-job-x"
description: |-
  Manages a Scheduler Job.
---

# azurerm_scheduler_job

Manages a Scheduler Job.

## Example Usage

Complete examples of how to use the `azurerm_scheduler_job` resource can be found [in the `./examples/scheduler-job` folder within the Github Repository](https://github.com/terraform-providers/terraform-provider-azurerm/tree/master/examples/scheduler-job)


```hcl
resource "azurerm_resource_group" "example" {
  # ...
}

resource "azurerm_scheduler_job" "example" {
  name                = "example-job"
  resource_group_name = "${azurerm_resource_group.example.name}"
  job_collection_name = "${azurerm_scheduler_job_collection.example.name}"

  # re-enable it each run
  state = "enabled"

  action_web {
    # defaults to get
    url = "http://this.url.fails"
  }

  # default start time is now
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Scheduler Job. Changing this forces a new resource to be created.

* `resource_group_name` - (Required) The name of the resource group in which to create the Scheduler Job. Changing this forces a new resource to be created.

* `job_collection_name` - (Required) Specifies the name of the Scheduler Job Collection in which the Job should exist. Changing this forces a new resource to be created.

* `action_web` - (Optional) A `action_web` block defining the job action as described below. Note this is identical to an `error_action_web` block.

~> **NOTE** At least one of `error_action_web` or `action_storage_queue` needs to be set.

* `action_storage_queue` - (Optional) A `action_storage_queue` block defining a storage queue job action as described below. Note this is identical to an `error_action_storage_queue` block.

* `error_action_web` - (Optional) A `error_action_web` block defining the action to take on an error as described below. Note this is identical to an `action_web` block.

* `error_action_storage_queue` - (Optional) A `error_action_storage_queue` block defining the a web action to take on an error as described below. Note this is identical to an `action_storage_queue` block.

* `retry` - (Optional) A `retry` block defining how to retry as described below.

* `recurrence` - (Optional) A `recurrence` block defining a job occurrence schedule.

* `start_time` - (Optional) The time the first instance of the job is to start running at.

* `state` - (Optional) The sets or gets the current state of the job. Can be set to either `Enabled` or `Completed`


`web_action` & `error_web_action` block supports the following:

* `url` - (Required) Specifies the URL of the web request. Must be HTTPS for authenticated requests.
* `method` - (Optional) Specifies the method of the request. Defaults to `Get` and must be one of `Get`, `Put`, `Post`, `Delete`.
* `body` - (Optional) Specifies the request body.
* `headers` - (Optional) A map specifying the headers sent with the request.
* `authentication_basic` - (Optional) An `authentication_active_directory` block which defines the Active Directory oauth configuration to use.
* `authentication_certificate` - (Optional) An `authentication_certificate` block which defines the client certificate information to be use.
* `authentication_active_directory` - (Optional) An `authentication_active_directory` block which defines the OAUTH Active Directory information to use.


`authentication_basic` block supports the following:

* `username` - (Required) Specifies the username to use.
* `password` - (Required) Specifies the password to use.


`authentication_certificate` block supports the following:

* `pfx` - (Required) Specifies the pfx certificate in base-64 format.
* `password` - (Required) Specifies the certificate password.


`authentication_active_directory` block supports the following:

* `client_id` - (Required) Specifies the client ID to use.
* `tenant_id` - (Required) Specifies the tenant ID to use.
* `client_secret` - (Required) Specifies the secret to use.
* `audience` - (Optional) Specifies the audience.

`action_storage_queue` & `error_action_storage_queue` block supports the following:

* `storage_account_name` - (Required) Specifies the the storage account name.
* `storage_queue_name` - (Required) Specifies the the storage account queue.
* `sas_token` - (Required) Specifies a SAS token/key to authenticate with.
* `message` - (Required) The message to send into the queue.

`retry` block supports the following:

* `interval` - (Required) Specifies the duration between retries.
* `count` - (Required) Specifies the number of times a retry should be attempted.

`recurrence` block supports the following:

* `frequency` - (Required) Specifies the frequency of recurrence. Must be one of `Minute`, `Hour`, `Day`, `Week`, `Month`.
* `interval` - (Optional) Specifies the interval between executions. Defaults to `1`.
* `count` - (Optional) Specifies the maximum number of times that the job should run.
* `end_time` - (Optional) Specifies the time at which the job will cease running. Must be less then 500 days into the future.
* `minutes` - (Optional) Specifies the minutes of the hour that the job should execute at. Must be between `0` and `59`
* `hours` - (Optional) Specifies the hours of the day that the job should execute at. Must be between `0` and `23`
* `week_days` - (Optional) Specifies the days of the week that the job should execute on. Must be one of `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`. Only applies when `Week` is used for frequency.
* `month_days` - (Optional) Specifies the days of the month that the job should execute on. Must be non zero and between `-1` and `31`. Only applies when `Month` is used for frequency.
* `monthly_occurrences` - (Optional) Specifies specific monthly occurrences like "last sunday of the month" with `monthly_occurrences` blocks. Only applies when `Month` is used for frequency.

`monthly_occurrences` block supports the following:

* `day` - (Optional) Specifies the day of the week that the job should execute on. Must be one of  one of `Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`.
* `occurrence` - (Optional) Specifies the week the job should run on. For example  `1` for the first week, `-1` for the last week of the month. Must be between `-5` and `5`.


## Attributes Reference

The following attributes are exported:

* `id` - The Scheduler Job ID.

`authentication_certificate` block exports the following:

* `thumbprint` - (Computed) The certificate thumbprint.
* `expiration` - (Computed)  The certificate expiration date.
* `subject_name` - (Computed) The certificate's certificate subject name.

## Import

Scheduler Job can be imported using a `resource id`, e.g.

```shell
terraform import azurerm_scheduler_job.job1 /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Scheduler/jobCollections/jobCollection1/jobs/job1
```
