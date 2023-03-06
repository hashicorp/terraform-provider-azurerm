---
subcategory: "Sentinel"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_sentinel_metadata"
description: |-
  Manages a Sentinel Metadata.
---

# azurerm_sentinel_metadata

Manages a Sentinel Metadata.

## Example Usage

```hcl
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-resources"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "pergb2018"
}

resource "azurerm_log_analytics_solution" "example" {
  solution_name         = "SecurityInsights"
  location              = azurerm_resource_group.example.location
  resource_group_name   = azurerm_resource_group.example.name
  workspace_resource_id = azurerm_log_analytics_workspace.example.id
  workspace_name        = azurerm_log_analytics_workspace.example.name

  plan {
    publisher = "Microsoft"
    product   = "OMSGallery/SecurityInsights"
  }
}

resource "azurerm_sentinel_alert_rule_nrt" "example" {
  name                       = "example"
  log_analytics_workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  display_name               = "example"
  severity                   = "High"
  query                      = <<QUERY
AzureActivity |
  where OperationName == "Create or Update Virtual Machine" or OperationName =="Create Deployment" |
  where ActivityStatus == "Succeeded" |
  make-series dcount(ResourceId) default=0 on EventSubmissionTimestamp in range(ago(7d), now(), 1d) by Caller
QUERY
}

resource "azurerm_sentinel_metadata" "example" {
  name         = "exampl"
  workspace_id = azurerm_log_analytics_solution.example.workspace_resource_id
  content_id   = azurerm_sentinel_alert_rule_nrt.example.name
  kind         = "AnalyticsRule"
  parent_id    = azurerm_sentinel_alert_rule_nrt.example.id
}
```

## Arguments Reference

The following arguments are supported:

* `content_id` - (Required) The ID of the content. Used to identify dependencies and content from solutions or community.

* `kind` - (Required) The kind of content the metadata is for. Possible values are `AnalyticsRule`, `AnalyticsRuleTemplate`, `AutomationRule`, `AzureFunction`, `DataConnector`, `DataType`, `HuntingQuery`, `InvestigationQuery`, `LogicAppsCustomConnector`, `Parser`, `Playbook`, `PlaybookTemplate`, `Solution`, `Watchlist`, `WatchlistTemplate`, `Workbook` and `WorkbookTemplate`.

* `name` - (Required) The name which should be used for this Sentinel Metadata. Changing this forces a new Sentinel Metadata to be created.

* `parent_id` - (Required) The ID of the parent resource ID of the content item, which the metadata belongs to.

* `workspace_id` - (Required) The ID of the Log Analytics Workspace. Changing this forces a new Sentinel Metadata to be created.

---

* `author` - (Optional) An `author` blocks as defined below.

* `category` - (Optional) A `category` block as defined below.

* `content_schema_version` - (Optional) Schema version of the content. Can be used to distinguish between flow based on the schema version.

* `custom_version` - (Optional) The Custom version of the content.

* `dependency` - (Optional) A JSON formatted `dependency` block as defined below. Dependency for the content item, what other content items it requires to work.

* `first_publish_date` - (Optional) The first publish date of solution content item.

* `icon_id` - (Optional) The ID of the icon, this id can be fetched from the solution template.

* `last_publish_date` - (Optional) The last publish date of solution content item.

* `preview_image` - (Optional) Specifies a list of preview image file names. These will be taken from solution artifacts.

* `preview_image_dark` - (Optional) Specifies a list of preview image file names used for dark theme. These will be taken from solution artifacts.

* `providers` - (Optional) Specifies a list of providers for the solution content item.

* `source` - (Optional) A `source` block as defined below.

* `support` - (Optional) A `support` block as defined below.

* `threat_analysis_tactics` - (Optional) Specifies a list of tactics the resource covers.

* `threat_analysis_techniques` - (Optional) Specifies a list of techniques the resource covers.

* `version` - (Optional) Version of the content.

---

A `author` block supports the following:

* `name` - (Optional) The name of the author, company or person.

* `email` - (Optional) The email address of the author contact.

* `link` - (Optional) The link for author/vendor page.

---

A `category` block supports the following:

* `domains` - (Optional) Specifies a list of domains for the solution content item.

* `verticals` - (Optional) Specifies a list of industry verticals for the solution content item.

---

A `dependency` block supports the following:

* `contentId` - (Optional) ID of the content item that is depended on.

* `kind` - (Optional) Type of the content item that is depended on.

* `version` - (Optional) Version of the content item that is depended on.

* `operator` - (Optional) Operator used for list of dependencies in `criteria` array.

* `criteria` - (Optional) Specifies a list of `dependency` which must be fulfilled, according to the `operator`.

---

A `source` block supports the following:

* `name` - (Optional) The name of the content source, repo name, solution name, Log Analytics Workspace name, etc.

* `kind` - (Required) The kind of the content source. Possible values are `LocalWorkspace`, `Communtity`, `Solution` and `SourceRepository`.

* `id` - (Optional) The id of the content source, the solution ID, Log Analytics Workspace name etc.

---

A `support` block supports the following:

* `tier` - (Required) The type of support for content item. Possible values are `Microsoft`, `Partner` and `Community`.

* `email` - (Optional) The email address of the support contact.

* `link` - (Optional) The link for support help.

* `name` - (Optional) The name of the support contact.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported: 

* `id` - The ID of the Sentinel Metadata.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Sentinel Metadata.
* `read` - (Defaults to 5 minutes) Used when retrieving the Sentinel Metadata.
* `update` - (Defaults to 30 minutes) Used when updating the Sentinel Metadata.
* `delete` - (Defaults to 30 minutes) Used when deleting the Sentinel Metadata.

## Import

Sentinel Metadata can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_sentinel_metadata.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourcegroup1/providers/Microsoft.OperationalInsights/workspaces/workspace1/providers/Microsoft.SecurityInsights/metadata/metadata1
```
