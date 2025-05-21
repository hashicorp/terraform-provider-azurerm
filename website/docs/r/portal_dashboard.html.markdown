---
subcategory: "Portal"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_portal_dashboard"
description: |-
  Manages a shared dashboard in the Azure Portal.
---

# azurerm_portal_dashboard

Manages a shared dashboard in the Azure Portal.

## Example Usage

```hcl
variable "md_content" {
  description = "Content for the MD tile"
  default     = "# Hello all :)"
}

variable "video_link" {
  description = "Link to a video"
  default     = "https://www.youtube.com/watch?v=......"
}

data "azurerm_subscription" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "mygroup"
  location = "West Europe"
}

resource "azurerm_portal_dashboard" "my-board" {
  name                = "my-cool-dashboard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  tags = {
    source = "terraform"
  }
  dashboard_properties = <<DASH
{
   "lenses": {
        "0": {
            "order": 0,
            "parts": {
                "0": {
                    "position": {
                        "x": 0,
                        "y": 0,
                        "rowSpan": 2,
                        "colSpan": 3
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/MarkdownPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "content": "${var.md_content}",
                                    "subtitle": "",
                                    "title": ""
                                }
                            }
                        }
                    }
                },               
                "1": {
                    "position": {
                        "x": 5,
                        "y": 0,
                        "rowSpan": 4,
                        "colSpan": 6
                    },
                    "metadata": {
                        "inputs": [],
                        "type": "Extension/HubsExtension/PartType/VideoPart",
                        "settings": {
                            "content": {
                                "settings": {
                                    "title": "Important Information",
                                    "subtitle": "",
                                    "src": "${var.video_link}",
                                    "autoplay": true
                                }
                            }
                        }
                    }
                },
                "2": {
                    "position": {
                        "x": 0,
                        "y": 4,
                        "rowSpan": 4,
                        "colSpan": 6
                    },
                    "metadata": {
                        "inputs": [
                            {
                                "name": "ComponentId",
                                "value": "/subscriptions/${data.azurerm_subscription.current.subscription_id}/resourceGroups/myRG/providers/microsoft.insights/components/myWebApp"
                            }
                        ],
                        "type": "Extension/AppInsightsExtension/PartType/AppMapGalPt",
                        "settings": {},
                        "asset": {
                            "idInputName": "ComponentId",
                            "type": "ApplicationInsights"
                        }
                    }
                }              
            }
        }
    },
    "metadata": {
        "model": {
            "timeRange": {
                "value": {
                    "relative": {
                        "duration": 24,
                        "timeUnit": 1
                    }
                },
                "type": "MsPortalFx.Composition.Configuration.ValueTypes.TimeRange"
            },
            "filterLocale": {
                "value": "en-us"
            },
            "filters": {
                "value": {
                    "MsPortalFx_TimeRange": {
                        "model": {
                            "format": "utc",
                            "granularity": "auto",
                            "relative": "24h"
                        },
                        "displayCache": {
                            "name": "UTC Time",
                            "value": "Past 24 hours"
                        },
                        "filteredPartIds": [
                            "StartboardPart-UnboundPart-ae44fef5-76b8-46b0-86f0-2b3f47bad1c7"
                        ]
                    }
                }
            }
        }
    }
}
DASH
}

```

It is recommended to follow the steps outlined
[here](https://docs.microsoft.com/azure/azure-portal/azure-portal-dashboards-create-programmatically#fetch-the-json-representation-of-the-dashboard) to create a Dashboard in the Portal and extract the relevant JSON to use in this resource. From the extracted JSON, the contents of the `properties: {}` object can used. Variables can be injected as needed - see above example.

### Using a `template_file` data source or the `templatefile` function

Since the contents of the dashboard JSON can be quite lengthy, use a template file to improve readability:

`dash.tpl`:

```JSON
{
    "lenses": {
         "0": {
             "order": 0,
             "parts": {
                 "0": {
                     "position": {
                         "x": 0,
                         "y": 0,
                         "rowSpan": 2,
                         "colSpan": 3
                     },
                     "metadata": {
                         "inputs": [],
                         "type": "Extension/HubsExtension/PartType/MarkdownPart",
                         "settings": {
                             "content": {
                                 "settings": {
                                     "content": "${md_content}", // <-- note the 'var.' is dropped
                                     "subtitle": "",
                                     "title": ""
                                 }
                             }
                         }
                     }
                 },  
                 ...
                 ...
```

This is then referenced in the `.tf` file by using a [`template_file`](https://www.terraform.io/docs/providers/template/d/file.html) data source (terraform 0.11 or earlier), or the [`templatefile`](https://www.terraform.io/docs/configuration/functions/templatefile.html) function (terraform 0.12+).

`main.tf` (terraform 0.11 or earlier):

```hcl
data "template_file" "dash-template" {
  template = "${file("${path.module}/dash.tpl")}"
  vars = {
    md_content = "Variable content here!"
    video_link = "https://www.youtube.com/watch?v=......"
    sub_id     = data.azurerm_subscription.current.subscription_id
  }
}

#...

resource "azurerm_dashboard" "my-board" {
  name                = "my-cool-dashboard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  tags = {
    source = "terraform"
  }
  dashboard_properties = data.template_file.dash-template.rendered
}

```

`main.tf` (terraform 0.12+)

```hcl
resource "azurerm_dashboard" "my-board" {
  name                = "my-cool-dashboard"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  tags = {
    source = "terraform"
  }
  dashboard_properties = templatefile("dash.tpl",
    {
      md_content = "Variable content here!",
      video_link = "https://www.youtube.com/watch?v=......",
      sub_id     = data.azurerm_subscription.current.subscription_id
  })
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Specifies the name of the Shared Dashboard. Changing this forces a new resource to be created.

-> **Note:** You can specify a tag with the key `hidden-title` to set a more user-friendly title for this Dashboard.

* `resource_group_name` - (Required) The name of the resource group in which to create the dashboard. Changing this forces a new resource to be created.

* `location` - (Required) Specifies the supported Azure location where the resource exists. Changing this forces a new resource to be created.

* `dashboard_properties` - (Required) JSON data representing dashboard body. See above for details on how to obtain this from the Portal.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Dashboard.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/language/resources/syntax#operation-timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the Dashboard.
* `read` - (Defaults to 5 minutes) Used when retrieving the Dashboard.
* `update` - (Defaults to 30 minutes) Used when updating the Dashboard.
* `delete` - (Defaults to 30 minutes) Used when deleting the Dashboard.

## Import

Dashboards can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_portal_dashboard.my-board /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Portal/dashboards/00000000-0000-0000-0000-000000000000
```

Note the URI in the above sample can be found using the Resource Explorer tool in the Azure Portal.
