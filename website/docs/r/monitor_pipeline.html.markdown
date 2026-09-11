---
subcategory: "Monitor"
layout: "azurerm"
page_title: "Azure Resource Manager: azurerm_monitor_pipeline"
description: |-
  Manages an Azure Monitor pipeline.
---

# azurerm_monitor_pipeline

Manages an Azure Monitor pipeline.

-> **Note:** An Azure Monitor pipeline runs on an Arc-enabled Kubernetes cluster to receive, process, and forward telemetry (such as Syslog, CEF, and OpenTelemetry logs) to Azure Monitor. More information about prerequisites can be found in the [Azure documentation](https://learn.microsoft.com/azure/azure-monitor/data-collection/pipeline-configure#prerequisites).

## Example Usage

```hcl
resource "azurerm_resource_group" "example" {
  name     = "example-resource-group"
  location = "West Europe"
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "example-log-analytics-workspace"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_table_custom_log" "example" {
  name         = "ExampleTable_CL"
  workspace_id = azurerm_log_analytics_workspace.example.id

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
  column {
    name = "Message"
    type = "string"
  }
  column {
    name = "ResourceColumnUpd"
    type = "string"
  }
  column {
    name = "ScopeColumnUpd"
    type = "string"
  }
}

resource "azurerm_monitor_data_collection_endpoint" "example" {
  name                = "example-monitor-data-collection-endpoint"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_monitor_data_collection_rule" "example" {
  name                        = "example-monitor-data-collection-rule"
  resource_group_name         = azurerm_resource_group.example.name
  location                    = azurerm_resource_group.example.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.example.id

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.example.id
      name                  = "example-destination-log"
    }
  }

  data_flow {
    streams      = ["Custom-${azurerm_log_analytics_workspace_table_custom_log.example.name}"]
    destinations = ["example-destination-log"]
  }

  stream_declaration {
    stream_name = "Custom-${azurerm_log_analytics_workspace_table_custom_log.example.name}"
    column {
      name = "TimeGenerated"
      type = "datetime"
    }
    column {
      name = "Message"
      type = "string"
    }
    column {
      name = "ResourceColumnUpd"
      type = "string"
    }
    column {
      name = "ScopeColumnUpd"
      type = "string"
    }
  }
}

data "azurerm_extended_location_custom_location" "example" {
  name                = "existing-extended-location-custom-location"
  resource_group_name = "existing-resource-group"
}

resource "azurerm_monitor_pipeline" "example" {
  name                   = "example-monitor-pipeline"
  resource_group_name    = azurerm_resource_group.example.name
  location               = azurerm_resource_group.example.location
  custom_location_id     = data.azurerm_extended_location_custom_location.example.id
  persistent_volume_name = "example-monitor-pipeline-pv"
  replicas               = 1

  execution_placement_constraint {
    capability = "gpu-enabled"
    operator   = "Exists"
  }

  execution_placement_constraint {
    capability = "node-type"
    operator   = "In"
    values     = ["high-cpu", "dedicated"]
  }

  azure_monitor_workspace_log_exporter {
    name = "example-exporter"

    api {
      data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.example.logs_ingestion_endpoint
      data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.example.immutable_id
      stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.example.name}"

      schema {
        record_map {
          from = "body"
          to   = "Message"
        }
        record_map {
          from = "time_unix_nano"
          to   = "TimeGenerated"
        }
        resource_map {
          from = "ResourceColumn"
          to   = "ResourceColumnUpd"
        }
        scope_map {
          from = "ScopeColumn"
          to   = "ScopeColumnUpd"
        }
      }
    }

    persistence_maximum_storage_usage_in_gb = 100
    persistence_retention_period_in_minutes = 10
  }

  batch_processor {
    name = "example-batch-processor"
  }

  transform_language_processor {
    name                = "example-transform-processor"
    transform_statement = "source | extend FooColumn = 'bar'"
  }

  microsoft_common_security_log_processor {
    name = "example-cef-processor"
  }

  microsoft_syslog_processor {
    name = "example-syslog-processor"
  }

  syslog_receiver {
    name                       = "example-syslog-receiver"
    allow_skip_priority_header = true
    endpoint                   = "0.0.0.0:514"
    allowed_formats            = ["syslogRfc5424", "syslogRfc3164"]
    tls_configuration_name     = "example-disabled-tls"
  }

  otlp_receiver {
    name                   = "example-otlp-receiver"
    tls_configuration_name = "example-mutual-tls"
    endpoint               = "0.0.0.0:4317"
  }

  tls_configuration {
    name = "example-disabled-tls"
    mode = "disabled"
  }

  tls_configuration {
    name = "example-mutual-tls"
    mode = "mutualTls"

    client_certificate_authority {
      location     = "client-ca-bundle"
      sub_location = "ca.crt"
      type         = "kubernetesSecret"
    }

    tls_certificate {
      certificate {
        location     = "pipeline-tls-cert"
        sub_location = "tls.crt"
        type         = "kubernetesSecret"
      }
      private_key {
        location     = "pipeline-tls-cert"
        sub_location = "tls.key"
      }
    }
  }

  tls_configuration {
    name = "example-server-only-tls"
    mode = "serverOnly"

    tls_certificate {
      certificate {
        location     = "pipeline-tls-cert"
        sub_location = "tls.crt"
        type         = "kubernetesSecret"
      }
      private_key {
        location     = "pipeline-tls-cert"
        sub_location = "tls.key"
      }
    }
  }

  log_pipeline {
    name       = "example-syslog-pipeline"
    exporters  = ["example-exporter"]
    receivers  = ["example-syslog-receiver"]
    processors = ["example-batch-processor", "example-cef-processor", "example-syslog-processor"]
  }

  log_pipeline {
    name       = "example-otlp-pipeline"
    exporters  = ["example-exporter"]
    receivers  = ["example-otlp-receiver"]
    processors = ["example-transform-processor"]
  }

  tags = {
    environment = "example"
  }
}
```

-> **Note:** A complete example that provisions an Arc-enabled Kubernetes cluster and its Azure Monitor pipeline prerequisites can be found in [the `examples/azure-monitoring/monitor-pipeline` directory within the GitHub Repository](https://github.com/hashicorp/terraform-provider-azurerm/tree/main/examples/azure-monitoring/monitor-pipeline).

## Arguments Reference

The following arguments are supported:

* `name` - (Required) The name which should be used for this Azure Monitor pipeline. Changing this forces a new resource to be created.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `resource_group_name` - (Required) The name of the Resource Group where the Azure Monitor pipeline should exist. Changing this forces a new resource to be created.

* `location` - (Required) The Azure Region where the Azure Monitor pipeline should exist. Changing this forces a new resource to be created.

* `custom_location_id` - (Required) The ID of the Custom Location where the Azure Monitor pipeline should exist. Changing this forces a new resource to be created.

* `log_pipeline` - (Required) One or more `log_pipeline` blocks as defined below.

* `azure_monitor_workspace_log_exporter` - (Optional) One or more `azure_monitor_workspace_log_exporter` blocks as defined below.

* `batch_processor` - (Optional) One or more `batch_processor` blocks as defined below.

* `execution_placement_constraint` - (Optional) One or more `execution_placement_constraint` blocks as defined below.

* `microsoft_common_security_log_processor` - (Optional) One or more `microsoft_common_security_log_processor` blocks as defined below.

* `microsoft_syslog_processor` - (Optional) One or more `microsoft_syslog_processor` blocks as defined below.

* `otlp_receiver` - (Optional) One or more `otlp_receiver` blocks as defined below.

* `persistent_volume_name` - (Optional) The name of the Kubernetes persistent volume mounted for durable storage. Changing this forces a new resource to be created.

* `replicas` - (Optional) The number of replicas of the Azure Monitor pipeline instance.

-> **Note:** `replicas` must be at least `1`.

* `syslog_receiver` - (Optional) One or more `syslog_receiver` blocks as defined below.

* `tls_configuration` - (Optional) One or more `tls_configuration` blocks as defined below.

* `transform_language_processor` - (Optional) One or more `transform_language_processor` blocks as defined below.

* `tags` - (Optional) A mapping of tags which should be assigned to the Azure Monitor pipeline.

---

An `api` block supports the following:

* `data_collection_endpoint_url` - (Required) The logs ingestion URL of the Data Collection Endpoint.

* `data_collection_rule_immutable_id` - (Required) The immutable ID of the Data Collection Rule to send data to.

-> **Note:** `data_collection_rule_immutable_id` must use the form `dcr-` followed by 32 hexadecimal characters.

* `schema` - (Required) A `schema` block as defined below.

* `stream` - (Required) The name of the stream in the Data Collection Rule that the exported data is sent to.

---

An `azure_monitor_workspace_log_exporter` block supports the following:

* `name` - (Required) The name which should be used for this exporter.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `api` - (Required) An `api` block as defined above.

* `persistence_maximum_storage_usage_in_gb` - (Optional) The maximum local storage the exporter is allowed to use, in gigabytes.

* `persistence_retention_period_in_minutes` - (Optional) The retention period for persisted data that has not yet been exported, in minutes. Possible values range between `1` and `2880`.

~> **Note:** `persistent_volume_name` must be set when `persistence_maximum_storage_usage_in_gb` or `persistence_retention_period_in_minutes` is configured.

---

A `batch_processor` block supports the following:

* `name` - (Required) The name which should be used for this processor.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `batch_size` - (Optional) The size of the batch. Possible values range between `10` and `100000`. Defaults to `8192`.

* `timeout_in_milliseconds` - (Optional) The batch timeout, in milliseconds. Possible values range between `10` and `300000`. Defaults to `300000`.

---

An `execution_placement_constraint` block supports the following:

* `capability` - (Required) The capability or attribute key used to match compute unit properties.

* `operator` - (Required) The match operator to use for this constraint. Possible values are `DoesNotExist`, `Exists`, `In`, and `NotIn`.

* `values` - (Optional) A list of values to match against.

~> **Note:** `values` is required when `operator` is `In` or `NotIn`, and must not be set when `operator` is `Exists` or `DoesNotExist`.

---

A `certificate` block supports the following:

* `location` - (Required) The location of the certificate source.

* `sub_location` - (Required) The sub-location within the certificate source.

* `type` - (Required) The type of certificate source. Possible values are `kubernetesConfigMap` and `kubernetesSecret`.

---

A `client_certificate_authority` block supports the following:

* `location` - (Required) The location of the client CA certificate source.

* `sub_location` - (Required) The sub-location within the client CA certificate source.

* `type` - (Required) The type of client CA certificate source. Possible values are `kubernetesConfigMap` and `kubernetesSecret`.

---

A `private_key` block supports the following:

* `location` - (Required) The location of the private key source.

* `sub_location` - (Required) The sub-location within the private key source.

---

A `log_pipeline` block supports the following:

* `name` - (Required) The name which should be used for this pipeline.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `exporters` - (Required) A list of exporter names referenced by this pipeline.

* `receivers` - (Required) A list of receiver names referenced by this pipeline.

* `processors` - (Optional) A list of processor names referenced by this pipeline.

---

A `microsoft_common_security_log_processor` block supports the following:

* `name` - (Required) The name which should be used for this processor.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

---

A `microsoft_syslog_processor` block supports the following:

* `name` - (Required) The name which should be used for this processor.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

---

An `otlp_receiver` block supports the following:

* `name` - (Required) The name which should be used for this receiver.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `endpoint` - (Required) The endpoint the OTLP receiver listens on.

-> **Note:** `endpoint` must use the format `<host>:<port>`, for example `0.0.0.0:4317`, with a port between `1` and `65535`.

* `tls_configuration_name` - (Optional) The name of the `tls_configuration` block to secure this receiver with. When not specified, the default TLS configuration is used.

~> **Note:** `tls_configuration_name` must reference the `name` of a `tls_configuration` block defined on this resource.

---

A `record_map` block supports the following:

* `from` - (Required) The source field name to map from.

* `to` - (Required) The destination field name to map to.

~> **Note:** `record_map` must contain at least one entry with `to` set to `TimeGenerated`.

---

A `resource_map` block supports the following:

* `from` - (Required) The source field name to map from.

* `to` - (Required) The destination field name to map to.

---

A `schema` block supports the following:

* `record_map` - (Required) One or more `record_map` blocks as defined above.

* `resource_map` - (Optional) One or more `resource_map` blocks as defined above.

* `scope_map` - (Optional) One or more `scope_map` blocks as defined below.

---

A `scope_map` block supports the following:

* `from` - (Required) The source field name to map from.

* `to` - (Required) The destination field name to map to.

---

A `syslog_receiver` block supports the following:

* `name` - (Required) The name which should be used for this receiver.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `endpoint` - (Required) The endpoint the Syslog receiver listens on.

-> **Note:** `endpoint` must use the format `<host>:<port>`, for example `0.0.0.0:514`, with a port between `1` and `65535`.

* `allow_skip_priority_header` - (Optional) Whether the receiver allows parsing messages without the PRI header. Defaults to `false`.

* `allowed_formats` - (Optional) A list of allowed message formats for Syslog and CEF ingestion. Possible values are `all`, `cefRfc3164`, `cefRfc5424`, `rawCef`, `syslogRfc3164`, and `syslogRfc5424`. Defaults to `["all"]`.

~> **Note:** `all` cannot be combined with other values in `allowed_formats`. When `allow_skip_priority_header` is `true`, `allowed_formats` must include `all`, `syslogRfc3164`, or `cefRfc3164`.

* `tls_configuration_name` - (Optional) The name of the `tls_configuration` block to secure this receiver with. When not specified, the default TLS configuration is used.

~> **Note:** `tls_configuration_name` must reference the `name` of a `tls_configuration` block defined on this resource. It is not supported when `transport_protocol` is `udp`.

* `transport_protocol` - (Optional) The transport protocol used by the receiver. Possible values are `tcp` and `udp`. Defaults to `tcp`.

---

A `tls_configuration` block supports the following:

* `name` - (Required) The name which should be used for this TLS configuration.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `client_certificate_authority` - (Optional) A `client_certificate_authority` block as defined above.

* `mode` - (Optional) The TLS security mode for receivers using this configuration. Possible values are `disabled`, `mutualTls`, and `serverOnly`. Defaults to `mutualTls`.

* `tls_certificate` - (Optional) A `tls_certificate` block as defined below.

~> **Note:** `client_certificate_authority` and `tls_certificate` must not be set when `mode` is `disabled`, and `client_certificate_authority` must not be set when `mode` is `serverOnly`.

---

A `tls_certificate` block supports the following:

* `certificate` - (Required) A `certificate` block as defined above.

* `private_key` - (Required) A `private_key` block as defined above.

---

A `transform_language_processor` block supports the following:

* `name` - (Required) The name which should be used for this processor.

-> **Note:** `name` must be between 4 and 33 characters, contain only letters, numbers, and hyphens, and must not start or end with a hyphen.

* `transform_statement` - (Required) The transform statement to execute over the data passing through the processor.

-> **Note:** `transform_statement` must be between 1 and 10000 characters.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the Azure Monitor pipeline.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - (Defaults to 1 hour) Used when creating the Azure Monitor pipeline.
* `read` - (Defaults to 5 minutes) Used when retrieving the Azure Monitor pipeline.
* `update` - (Defaults to 1 hour) Used when updating the Azure Monitor pipeline.
* `delete` - (Defaults to 1 hour) Used when deleting the Azure Monitor pipeline.

## Import

An Azure Monitor pipeline can be imported using the `resource id`, e.g.

```shell
terraform import azurerm_monitor_pipeline.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroup1/providers/Microsoft.Monitor/pipelineGroups/pipelineGroup1
```

## API Providers
<!-- This section is generated, changes will be overwritten -->
This resource uses the following Azure API Providers:

* `Microsoft.Monitor` - 2026-04-01
