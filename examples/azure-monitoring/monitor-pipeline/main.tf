# Copyright IBM Corp. 2014, 2025
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    azuread = {
      source  = "hashicorp/azuread"
      version = "~> 3.0"
    }
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 4.0"
    }
    local = {
      source  = "hashicorp/local"
      version = "~> 2.5"
    }
    random = {
      source  = "hashicorp/random"
      version = "~> 3.7"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.1"
    }
  }
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

data "azuread_service_principal" "custom_locations" {
  client_id = "bc313c14-388c-4e7d-a58e-70017303ee3b"
}

resource "random_string" "suffix" {
  length  = 8
  special = false
  upper   = false
}

locals {
  pipeline_namespace = "monitoring"
  storage_account    = "amp${random_string.suffix.result}"
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-monitor-pipeline"
  location = var.location
}

resource "tls_private_key" "arc_agent" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

locals {
  arc_agent_public_key_spki_b64 = replace(replace(replace(tls_private_key.arc_agent.public_key_pem, "-----BEGIN PUBLIC KEY-----\n", ""), "\n-----END PUBLIC KEY-----\n", ""), "\n", "")
  arc_agent_public_key_pkcs1    = substr(local.arc_agent_public_key_spki_b64, 32, -1)
}

resource "azurerm_kubernetes_cluster" "example" {
  name                = "${var.prefix}-monitor-pipeline-aks"
  resource_group_name = azurerm_resource_group.example.name
  node_resource_group = "${var.prefix}-monitor-pipeline-nodes"
  location            = azurerm_resource_group.example.location
  dns_prefix          = "${var.prefix}-monitor-pipeline"

  node_provisioning_profile {
    mode = "Manual"
  }

  default_node_pool {
    name       = "system"
    node_count = 1
    vm_size    = "Standard_D4s_v3"
  }

  identity {
    type = "SystemAssigned"
  }

  lifecycle {
    ignore_changes = [default_node_pool[0].upgrade_settings]
  }
}

resource "local_sensitive_file" "kubeconfig" {
  content         = azurerm_kubernetes_cluster.example.kube_config_raw
  filename        = "${path.root}/monitor-pipeline-kubeconfig"
  file_permission = "0600"
}

resource "azurerm_arc_kubernetes_cluster" "example" {
  name                         = "${var.prefix}-monitor-pipeline-arc"
  resource_group_name          = azurerm_resource_group.example.name
  location                     = azurerm_resource_group.example.location
  agent_public_key_certificate = local.arc_agent_public_key_pkcs1

  identity {
    type = "SystemAssigned"
  }

  provisioner "local-exec" {
    command = "python3 ${path.module}/testdata/install_agent.py --subscriptionId ${data.azurerm_client_config.current.subscription_id} --resourceGroupName ${azurerm_resource_group.example.name} --clusterName ${var.prefix}-monitor-pipeline-arc --location ${azurerm_resource_group.example.location} --tenantId ${data.azurerm_client_config.current.tenant_id} --privatePemEnvironmentVariable ARC_PRIVATE_PEM --kubeConfig '${local_sensitive_file.kubeconfig.filename}' --customLocationsOid ${data.azuread_service_principal.custom_locations.object_id}"

    environment = {
      ARC_PRIVATE_PEM = tls_private_key.arc_agent.private_key_pem
      KUBECONFIG      = local_sensitive_file.kubeconfig.filename
    }
  }

  depends_on = [local_sensitive_file.kubeconfig]
}

resource "azurerm_arc_kubernetes_cluster_extension" "cert_manager" {
  name           = "${var.prefix}-cert-manager"
  cluster_id     = azurerm_arc_kubernetes_cluster.example.id
  extension_type = "microsoft.certmanagement"
  release_train  = "preview"

  identity {
    type = "SystemAssigned"
  }
}

resource "azurerm_arc_kubernetes_cluster_extension" "pipeline" {
  name              = "amp-${random_string.suffix.result}"
  cluster_id        = azurerm_arc_kubernetes_cluster.example.id
  extension_type    = "microsoft.monitor.pipelinecontroller"
  release_train     = "Preview"
  release_namespace = local.pipeline_namespace

  identity {
    type = "SystemAssigned"
  }

  depends_on = [azurerm_arc_kubernetes_cluster_extension.cert_manager]
}

resource "azurerm_extended_location_custom_location" "example" {
  name                = "${var.prefix}-monitor-pipeline"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  namespace           = local.pipeline_namespace
  cluster_extension_ids = [
    azurerm_arc_kubernetes_cluster_extension.pipeline.id,
  ]
  host_resource_id = azurerm_arc_kubernetes_cluster.example.id
}

resource "azurerm_storage_account" "example" {
  name                     = local.storage_account
  resource_group_name      = azurerm_resource_group.example.name
  location                 = azurerm_resource_group.example.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "example" {
  name               = "pipeline"
  storage_account_id = azurerm_storage_account.example.id
  quota              = 5
}

resource "tls_private_key" "client_ca" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "client_ca" {
  private_key_pem       = tls_private_key.client_ca.private_key_pem
  validity_period_hours = 8760
  is_ca_certificate     = true
  allowed_uses          = ["cert_signing", "key_encipherment", "digital_signature", "server_auth", "client_auth"]

  subject {
    common_name = "${var.prefix}-monitor-pipeline-client-ca"
  }
}

resource "tls_private_key" "pipeline" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "pipeline" {
  private_key_pem       = tls_private_key.pipeline.private_key_pem
  validity_period_hours = 8760
  allowed_uses          = ["key_encipherment", "digital_signature", "server_auth", "client_auth"]

  subject {
    common_name = "${var.prefix}-monitor-pipeline-tls"
  }
}

resource "local_file" "pipeline_persistent_volume" {
  content = templatefile("${path.module}/testdata/pipeline_persistent_volume.yaml.tftpl", {
    pipeline_namespace     = local.pipeline_namespace
    persistent_volume_name = "example-monitor-pipeline-pv"
    resource_group         = azurerm_resource_group.example.name
    storage_account        = azurerm_storage_account.example.name
    storage_share          = azurerm_storage_share.example.name
  })
  filename = "${path.root}/pipeline-persistent-volume.yaml"
}

resource "terraform_data" "cluster_prerequisites" {
  provisioner "local-exec" {
    command = "bash ${path.module}/testdata/configure_pipeline_cluster_prereqs.sh && kubectl label nodes --all gpu-enabled=true high-memory=true node-type=dedicated --overwrite"

    environment = {
      CLIENT_CA_CERT                    = tls_self_signed_cert.client_ca.cert_pem
      KUBECONFIG                        = local_sensitive_file.kubeconfig.filename
      PIPELINE_NAMESPACE                = local.pipeline_namespace
      PIPELINE_PERSISTENT_VOLUME_CONFIG = local_file.pipeline_persistent_volume.filename
      PIPELINE_TLS_CERT                 = tls_self_signed_cert.pipeline.cert_pem
      PIPELINE_TLS_KEY                  = tls_private_key.pipeline.private_key_pem
      STORAGE_ACCOUNT                   = azurerm_storage_account.example.name
      STORAGE_ACCOUNT_KEY               = azurerm_storage_account.example.primary_access_key
    }
  }

  depends_on = [azurerm_extended_location_custom_location.example]
}

resource "azurerm_log_analytics_workspace" "example" {
  name                = "${var.prefix}-monitor-pipeline"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_table_custom_log" "example" {
  name         = "PipelineExample_CL"
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
  name                = "${var.prefix}-monitor-pipeline"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

resource "azurerm_monitor_data_collection_rule" "example" {
  name                        = "${var.prefix}-monitor-pipeline"
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

resource "azurerm_monitor_pipeline" "example" {
  name                = "${var.prefix}-monitor-pipeline"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  custom_location_id  = azurerm_extended_location_custom_location.example.id
  replicas            = 1

  execution_placement_constraint {
    capability = "gpu-enabled"
    operator   = "Exists"
  }

  execution_placement_constraint {
    capability = "node-type"
    operator   = "In"
    values     = ["high-cpu", "dedicated"]
  }

  exporter {
    name = "example-exporter"

    azure_monitor_workspace_logs {
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

      persistence {
        maximum_storage_usage_in_gb = 100
        retention_period_in_minutes = 10
      }
    }
  }

  processor {
    name = "example-batch-processor"
    type = "Batch"

    batch {
      batch_size              = 8192
      timeout_in_milliseconds = 300000
    }
  }

  processor {
    name                = "example-transform-processor"
    type                = "TransformLanguage"
    transform_statement = "source | extend FooColumn = 'bar'"
  }

  processor {
    name = "example-cef-processor"
    type = "MicrosoftCommonSecurityLog"
  }

  processor {
    name = "example-syslog-processor"
    type = "MicrosoftSyslog"
  }

  receiver {
    name                   = "example-syslog-receiver"
    type                   = "Syslog"
    tls_configuration_name = "example-disabled-tls"

    syslog {
      allow_skip_priority_header = true
      endpoint                   = "0.0.0.0:514"
      allowed_formats            = ["syslogRfc5424", "syslogRfc3164"]
    }
  }

  receiver {
    name                   = "example-otlp-receiver"
    type                   = "OTLP"
    tls_configuration_name = "example-mutual-tls"

    otlp {
      endpoint = "0.0.0.0:4317"
    }
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

  service {
    persistent_volume_name = "example-monitor-pipeline-pv"

    pipeline {
      name       = "example-syslog-pipeline"
      exporters  = ["example-exporter"]
      receivers  = ["example-syslog-receiver"]
      processors = ["example-batch-processor", "example-cef-processor", "example-syslog-processor"]
    }

    pipeline {
      name       = "example-otlp-pipeline"
      exporters  = ["example-exporter"]
      receivers  = ["example-otlp-receiver"]
      processors = ["example-transform-processor"]
    }
  }

  tags = {
    environment = "example"
  }

  depends_on = [terraform_data.cluster_prerequisites]
}