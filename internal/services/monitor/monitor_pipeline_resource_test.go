// Copyright IBM Corp. 2014, 2025
// SPDX-License-Identifier: MPL-2.0

package monitor_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/pointer"
	"github.com/hashicorp/go-azure-sdk/resource-manager/monitor/2026-04-01/pipelinegroups"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
)

type MonitorPipelineResource struct{}

func TestAccMonitorPipeline_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline", "test")
	r := MonitorPipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorPipeline_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline", "test")
	r := MonitorPipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config:      r.requiresImport(data),
			ExpectError: acceptance.RequiresImportError("azurerm_monitor_pipeline"),
		},
	})
}

func TestAccMonitorPipeline_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline", "test")
	r := MonitorPipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestAccMonitorPipeline_update(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_monitor_pipeline", "test")
	r := MonitorPipelineResource{}

	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.update(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r MonitorPipelineResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := pipelinegroups.ParsePipelineGroupID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.Monitor.PipelineGroupsClient.Get(ctx, *id)
	if err != nil {
		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return pointer.To(resp.Model != nil), nil
}

// template provisions a self-contained Arc-enabled Kubernetes cluster (an AKS cluster
// connected via the Azure Arc agent), the cert-manager and pipeline controller cluster
// extensions it needs, a Custom Location targeting it, and the supporting Log
// Analytics/Data Collection resources the exporter in this test sends data to.
func (r MonitorPipelineResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

# Fixed multi-tenant app ID for the Custom Locations RP; its per-tenant service principal
# object ID is required to enable Custom Locations on the Arc agent.
# https://learn.microsoft.com/en-us/azure/azure-arc/kubernetes/custom-locations#enable-custom-locations-on-your-cluster
data "azuread_service_principal" "custom_locations" {
  client_id = "bc313c14-388c-4e7d-a58e-70017303ee3b"
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-monitor-pipeline-%[1]d"
  location = "%[2]s"
}

# Arc agent key pair, generated entirely in Terraform. tls_private_key.private_key_pem is
# already PKCS#1 ("-----BEGIN RSA PRIVATE KEY-----"), which is what the Arc agent requires.
resource "tls_private_key" "arc_agent" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

locals {
  # tls only exposes the SPKI public key (public_key_pem). Azure Arc needs the
  # PKCS#1 (RSAPublicKey) form. For an RSA key >= 2048 bits the SPKI DER has a
  # fixed 24-byte prefix before the embedded PKCS#1; 24 bytes == 32 base64
  # chars, so stripping the first 32 chars of the base64 body yields PKCS#1.
  arc_agent_public_key_spki_b64 = replace(replace(replace(tls_private_key.arc_agent.public_key_pem, "-----BEGIN PUBLIC KEY-----\n", ""), "\n-----END PUBLIC KEY-----\n", ""), "\n", "")
  arc_agent_public_key_pkcs1    = substr(local.arc_agent_public_key_spki_b64, 32, -1)
}

resource "azurerm_kubernetes_cluster" "test" {
  name                = "acctest-aks-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  node_resource_group = "acctest-aks-nodes-%[1]d"
  location            = azurerm_resource_group.test.location
  dns_prefix          = "acctest-aks-%[1]d"

  node_provisioning_profile {
    mode = "Manual"
  }

  # cert-manager + the pipeline controller + diagnostics extensions together need more
  # headroom than a single small node provides.
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
  content         = azurerm_kubernetes_cluster.test.kube_config_raw
  filename        = %[3]q
  file_permission = "0600"
}

resource "azurerm_arc_kubernetes_cluster" "test" {
  name                         = "acctest-akcc-%[1]d"
  resource_group_name          = azurerm_resource_group.test.name
  location                     = azurerm_resource_group.test.location
  agent_public_key_certificate = local.arc_agent_public_key_pkcs1

  identity {
    type = "SystemAssigned"
  }

  provisioner "local-exec" {
    command = "python3 testdata/install_agent.py --subscriptionId ${data.azurerm_client_config.current.subscription_id} --resourceGroupName ${azurerm_resource_group.test.name} --clusterName acctest-akcc-%[1]d --location ${azurerm_resource_group.test.location} --tenantId ${data.azurerm_client_config.current.tenant_id} --privatePemEnvironmentVariable ARC_PRIVATE_PEM --kubeConfig '%[3]s' --customLocationsOid ${data.azuread_service_principal.custom_locations.object_id}"

    environment = {
      KUBECONFIG      = %[3]q
      ARC_PRIVATE_PEM = tls_private_key.arc_agent.private_key_pem
    }
  }

  depends_on = [
    local_sensitive_file.kubeconfig,
  ]
}

resource "azurerm_arc_kubernetes_cluster_extension" "cert_manager" {
  name           = "acctest-cm-%[1]d"
  cluster_id     = azurerm_arc_kubernetes_cluster.test.id
  extension_type = "microsoft.certmanagement"
  release_train  = "preview"

  identity {
    type = "SystemAssigned"
  }
}

# Keep this name short: the operator chart derives a Kubernetes Service name as
# "<name>-pipeline-operator-metrics-service", which must fit the 63-char DNS label limit.
resource "azurerm_arc_kubernetes_cluster_extension" "pipeline" {
  name              = "acctpipe%[1]d"
  cluster_id        = azurerm_arc_kubernetes_cluster.test.id
  extension_type    = "microsoft.monitor.pipelinecontroller"
  release_train     = "Preview"
  release_namespace = "monitoring"

  identity {
    type = "SystemAssigned"
  }

  depends_on = [
    azurerm_arc_kubernetes_cluster_extension.cert_manager,
  ]
}

resource "azurerm_extended_location_custom_location" "test" {
  name                = "acctest-cl-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  namespace           = "monitoring"
  cluster_extension_ids = [
    azurerm_arc_kubernetes_cluster_extension.pipeline.id,
  ]
  host_resource_id = azurerm_arc_kubernetes_cluster.test.id
}

# Azure Files backing for the pipeline's persistent volume, and self-signed TLS material for
# the mutualTls tlsConfiguration exercised by the complete test.
resource "azurerm_storage_account" "test" {
  name                     = "acct%[1]d"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "LRS"
}

resource "azurerm_storage_share" "test" {
  name               = "pipeline"
  storage_account_id = azurerm_storage_account.test.id
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
    common_name = "acctest-monitor-pipeline-%[1]d-client-ca"
  }
}

resource "tls_private_key" "pipeline_tls" {
  algorithm = "RSA"
  rsa_bits  = 2048
}

resource "tls_self_signed_cert" "pipeline_tls" {
  private_key_pem       = tls_private_key.pipeline_tls.private_key_pem
  validity_period_hours = 8760
  allowed_uses          = ["key_encipherment", "digital_signature", "server_auth", "client_auth"]

  subject {
    common_name = "acctest-monitor-pipeline-%[1]d-pipeline-tls"
  }
}

# Works around the cert-manager extension's default ClusterIssuers referencing "-current"
# secrets that nothing else creates, which otherwise leaves the pipeline group stuck at
# provisioningState "Creating" indefinitely. Also creates the mTLS secrets and the Azure
# Files-backed PersistentVolume the complete test's tls_configuration and
# persistent_volume_name reference.
resource "terraform_data" "cluster_prereqs" {
  provisioner "local-exec" {
    command = "export PATH=\"$HOME/go/bin:$PATH\"; bash testdata/configure_pipeline_cluster_prereqs.sh && kubectl label nodes --all gpu-enabled=true high-memory=true node-type=dedicated --overwrite && kubectl create configmap server-tls-cert-updated --namespace \"$PIPELINE_NAMESPACE\" --from-literal=tls-updated.crt=\"$PIPELINE_TLS_CERT\" --dry-run=client -o yaml | kubectl apply -f - && kubectl create secret generic server-tls-key-updated --namespace \"$PIPELINE_NAMESPACE\" --from-literal=tls-updated.key=\"$PIPELINE_TLS_KEY\" --dry-run=client -o yaml | kubectl apply -f -"

    environment = {
      KUBECONFIG          = %[3]q
      PIPELINE_NAMESPACE  = "monitoring"
      CLIENT_CA_CERT      = tls_self_signed_cert.client_ca.cert_pem
      PIPELINE_TLS_CERT   = tls_self_signed_cert.pipeline_tls.cert_pem
      PIPELINE_TLS_KEY    = tls_private_key.pipeline_tls.private_key_pem
      RESOURCE_GROUP      = azurerm_resource_group.test.name
      STORAGE_ACCOUNT     = azurerm_storage_account.test.name
      STORAGE_ACCOUNT_KEY = azurerm_storage_account.test.primary_access_key
      STORAGE_SHARE       = azurerm_storage_share.test.name
    }
  }

  depends_on = [
    azurerm_extended_location_custom_location.test,
  ]
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "acctest-law-%[1]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azurerm_log_analytics_workspace_table_custom_log" "test" {
  name         = "PipelineTestData%[1]d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
  column {
    name = "Message"
    type = "string"
  }
  column {
    name = "MessageUpdated"
    type = "string"
  }
  column {
    name = "ResourceColumnUpd"
    type = "string"
  }
  column {
    name = "ResourceColumnUpdV2"
    type = "string"
  }
  column {
    name = "ScopeColumnUpd"
    type = "string"
  }
  column {
    name = "ScopeColumnUpdV2"
    type = "string"
  }
}

resource "azurerm_log_analytics_workspace_table_custom_log" "updated" {
  name         = "PipelineUpdatedData%[1]d_CL"
  workspace_id = azurerm_log_analytics_workspace.test.id

  column {
    name = "TimeGenerated"
    type = "dateTime"
  }
  column {
    name = "MessageUpdated"
    type = "string"
  }
  column {
    name = "ResourceColumnUpdV2"
    type = "string"
  }
  column {
    name = "ScopeColumnUpdV2"
    type = "string"
  }
}

resource "azurerm_monitor_data_collection_endpoint" "test" {
  name                = "acctest-dce-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_endpoint" "updated" {
  name                = "acctest-dce-updated-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm_monitor_data_collection_rule" "test" {
  name                = "acctest-dcr-%[1]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
      name                  = "test-destination-log"
    }
  }

  data_flow {
    streams      = ["Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"]
    destinations = ["test-destination-log"]
  }

  stream_declaration {
    stream_name = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"
    column {
      name = "TimeGenerated"
      type = "datetime"
    }
    column {
      name = "Message"
      type = "string"
    }
    column {
      name = "MessageUpdated"
      type = "string"
    }
    column {
      name = "ResourceColumnUpd"
      type = "string"
    }
    column {
      name = "ResourceColumnUpdV2"
      type = "string"
    }
    column {
      name = "ScopeColumnUpd"
      type = "string"
    }
    column {
      name = "ScopeColumnUpdV2"
      type = "string"
    }
  }
}

resource "azurerm_monitor_data_collection_rule" "updated" {
  name                        = "acctest-dcr-updated-%[1]d"
  resource_group_name         = azurerm_resource_group.test.name
  location                    = azurerm_resource_group.test.location
  data_collection_endpoint_id = azurerm_monitor_data_collection_endpoint.updated.id

  destinations {
    log_analytics {
      workspace_resource_id = azurerm_log_analytics_workspace.test.id
      name                  = "updated-destination-log"
    }
  }

  data_flow {
    streams      = ["Custom-${azurerm_log_analytics_workspace_table_custom_log.updated.name}"]
    destinations = ["updated-destination-log"]
  }

  stream_declaration {
    stream_name = "Custom-${azurerm_log_analytics_workspace_table_custom_log.updated.name}"
    column {
      name = "TimeGenerated"
      type = "datetime"
    }
    column {
      name = "MessageUpdated"
      type = "string"
    }
    column {
      name = "ResourceColumnUpdV2"
      type = "string"
    }
    column {
      name = "ScopeColumnUpdV2"
      type = "string"
    }
  }
}
`, data.RandomInteger, data.Locations.Primary, filepath.Join(os.TempDir(), fmt.Sprintf("acctest-monitor-pipeline-kubeconfig-%d", data.RandomInteger)))
}

func (r MonitorPipelineResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  workspace_logs_exporter_name = "acctest-exporter"
  syslog_receiver_name         = "acctest-receiver"
}

resource "azurerm_monitor_pipeline" "test" {
  name                = "acctest-pg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_extended_location_custom_location.test.id

  exporter {
    name = local.workspace_logs_exporter_name

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"

        schema {
          record_map {
            from = "body"
            to   = "Message"
          }
          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }
    }
  }

  receiver {
    name = local.syslog_receiver_name
    type = "Syslog"

    syslog {
      endpoint = "0.0.0.0:514"
    }
  }

  service {
    pipeline {
      name      = "acctest-pipeline"
      exporters = [local.workspace_logs_exporter_name]
      receivers = [local.syslog_receiver_name]
    }
  }

  depends_on = [
    terraform_data.cluster_prereqs,
  ]
}
`, r.template(data), data.RandomInteger)
}

func (r MonitorPipelineResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_monitor_pipeline" "import" {
  name                = azurerm_monitor_pipeline.test.name
  resource_group_name = azurerm_monitor_pipeline.test.resource_group_name
  location            = azurerm_monitor_pipeline.test.location
  custom_location_id  = azurerm_monitor_pipeline.test.custom_location_id

  exporter {
    name = azurerm_monitor_pipeline.test.exporter.0.name

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.data_collection_endpoint_url
        data_collection_rule_immutable_id = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.data_collection_rule_immutable_id
        stream                            = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.stream

        schema {
          record_map {
            from = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.schema.0.record_map.0.from
            to   = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.schema.0.record_map.0.to
          }
          record_map {
            from = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.schema.0.record_map.1.from
            to   = azurerm_monitor_pipeline.test.exporter.0.azure_monitor_workspace_logs.0.api.0.schema.0.record_map.1.to
          }
        }
      }
    }
  }

  receiver {
    name = azurerm_monitor_pipeline.test.receiver.0.name
    type = azurerm_monitor_pipeline.test.receiver.0.type

    syslog {
      endpoint = azurerm_monitor_pipeline.test.receiver.0.syslog.0.endpoint
    }
  }

  service {
    pipeline {
      name      = azurerm_monitor_pipeline.test.service.0.pipeline.0.name
      exporters = azurerm_monitor_pipeline.test.service.0.pipeline.0.exporters
      receivers = azurerm_monitor_pipeline.test.service.0.pipeline.0.receivers
    }
  }
}
`, r.basic(data))
}

func (r MonitorPipelineResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  workspace_logs_exporter_name   = "acctest-exporter"
  batch_processor_name           = "acctest-batch-processor"
  transform_processor_name       = "acctest-transform-processor"
  cef_processor_name             = "acctest-cef-processor"
  syslog_processor_name          = "acctest-syslog-processor"
  syslog_receiver_name           = "acctest-receiver"
  otlp_receiver_name             = "acctest-otlp-receiver"
  tls_configuration_name         = "acctest-tls-config"
  mtls_configuration_name        = "acctest-mtls-config"
  server_only_configuration_name = "acctest-server-only-tls-config"
}

resource "azurerm_monitor_pipeline" "test" {
  name                = "acctest-pg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_extended_location_custom_location.test.id
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
    name = local.workspace_logs_exporter_name

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"

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
    name = local.batch_processor_name
    type = "Batch"

    batch {
      batch_size              = 8192
      timeout_in_milliseconds = 300000
    }
  }

  processor {
    name = local.transform_processor_name
    type = "TransformLanguage"

    transform_statement = "source | extend FooColumn = 'bar'"
  }

  processor {
    name = local.cef_processor_name
    type = "MicrosoftCommonSecurityLog"
  }

  processor {
    name = local.syslog_processor_name
    type = "MicrosoftSyslog"
  }

  receiver {
    name                   = local.syslog_receiver_name
    type                   = "Syslog"
    tls_configuration_name = local.tls_configuration_name

    syslog {
      allow_skip_priority_header = true
      endpoint                   = "0.0.0.0:514"
      allowed_formats            = ["syslogRfc5424", "syslogRfc3164"]
    }
  }

  receiver {
    name                   = local.otlp_receiver_name
    type                   = "OTLP"
    tls_configuration_name = local.mtls_configuration_name

    otlp {
      endpoint = "0.0.0.0:4317"
    }
  }

  tls_configuration {
    name = local.tls_configuration_name
    mode = "disabled"
  }

  tls_configuration {
    name = local.mtls_configuration_name
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
    name = local.server_only_configuration_name
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
    persistent_volume_name = "acctest-pipeline-pv"

    pipeline {
      name       = "acctest-pipeline"
      exporters  = [local.workspace_logs_exporter_name]
      receivers  = [local.syslog_receiver_name]
      processors = [local.batch_processor_name, local.cef_processor_name, local.syslog_processor_name]
    }

    pipeline {
      name       = "acctest-otlp-pipeline"
      exporters  = [local.workspace_logs_exporter_name]
      receivers  = [local.otlp_receiver_name]
      processors = [local.transform_processor_name]
    }
  }

  tags = {
    environment = "test"
  }

  depends_on = [
    terraform_data.cluster_prereqs,
  ]
}
`, r.template(data), data.RandomInteger)
}

// update reuses the same template() dependencies as complete() and only changes values on
// the azurerm_monitor_pipeline resource itself, so the update test exercises an
// in-place update of the pipeline group without recreating the Arc cluster/extensions.
func (r MonitorPipelineResource) update(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

locals {
  workspace_logs_exporter_name           = "acctest-exporter"
  workspace_logs_exporter_secondary_name = "acctest-exporter-2"
  batch_processor_name                   = "acctest-batch-processor"
  transform_processor_name               = "acctest-transform-processor"
  cef_processor_name                     = "acctest-cef-processor"
  syslog_processor_name                  = "acctest-syslog-processor"
  syslog_receiver_name                   = "acctest-receiver"
  syslog_receiver_secondary_name         = "acctest-receiver-2"
  otlp_receiver_name                     = "acctest-otlp-receiver"
  tls_configuration_name                 = "acctest-tls-config"
  mtls_configuration_name                = "acctest-mtls-config"
  server_only_configuration_name         = "acctest-server-only-tls-config"
}

resource "azurerm_monitor_pipeline" "test" {
  name                = "acctest-pg-%[2]d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  custom_location_id  = azurerm_extended_location_custom_location.test.id
  replicas            = 2

  exporter {
    name = local.workspace_logs_exporter_name

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.updated.logs_ingestion_endpoint
        data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.updated.immutable_id
        stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.updated.name}"

        schema {
          record_map {
            from = "body"
            to   = "MessageUpdated"
          }
          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
          resource_map {
            from = "ResourceColumn"
            to   = "ResourceColumnUpdV2"
          }
          scope_map {
            from = "ScopeColumn"
            to   = "ScopeColumnUpdV2"
          }
        }
      }

      persistence {
        maximum_storage_usage_in_gb = 200
        retention_period_in_minutes = 20
      }
    }
  }

  exporter {
    name = local.workspace_logs_exporter_secondary_name

    azure_monitor_workspace_logs {
      api {
        data_collection_endpoint_url      = azurerm_monitor_data_collection_endpoint.test.logs_ingestion_endpoint
        data_collection_rule_immutable_id = azurerm_monitor_data_collection_rule.test.immutable_id
        stream                            = "Custom-${azurerm_log_analytics_workspace_table_custom_log.test.name}"

        schema {
          record_map {
            from = "body"
            to   = "Message"
          }
          record_map {
            from = "time_unix_nano"
            to   = "TimeGenerated"
          }
        }
      }
    }
  }

  processor {
    name = local.batch_processor_name
    type = "Batch"

    batch {
      batch_size              = 4096
      timeout_in_milliseconds = 150000
    }
  }

  processor {
    name = local.transform_processor_name
    type = "TransformLanguage"

    transform_statement = "source | extend FooColumn = 'baz'"
  }

  processor {
    name = local.cef_processor_name
    type = "MicrosoftCommonSecurityLog"
  }

  processor {
    name = local.syslog_processor_name
    type = "MicrosoftSyslog"
  }

  receiver {
    name = local.syslog_receiver_name
    type = "Syslog"

    syslog {
      allow_skip_priority_header = false
      endpoint                   = "0.0.0.0:1514"
      allowed_formats            = ["syslogRfc3164"]
      transport_protocol         = "udp"
    }
  }

  receiver {
    name                   = local.otlp_receiver_name
    type                   = "OTLP"
    tls_configuration_name = local.mtls_configuration_name

    otlp {
      endpoint = "0.0.0.0:4318"
    }
  }

  receiver {
    name = local.syslog_receiver_secondary_name
    type = "Syslog"

    syslog {
      endpoint        = "0.0.0.0:6514"
      allowed_formats = ["all"]
    }
  }

  tls_configuration {
    name = local.mtls_configuration_name
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
    name = local.server_only_configuration_name
    mode = "mutualTls"

    client_certificate_authority {
      location     = "client-ca-bundle"
      sub_location = "ca.crt"
      type         = "kubernetesSecret"
    }

    tls_certificate {
      certificate {
        location     = "server-tls-cert-updated"
        sub_location = "tls-updated.crt"
        type         = "kubernetesConfigMap"
      }
      private_key {
        location     = "server-tls-key-updated"
        sub_location = "tls-updated.key"
      }
    }
  }

  tls_configuration {
    name = local.tls_configuration_name
    mode = "disabled"
  }

  service {
    persistent_volume_name = "acctest-pipeline-pv"

    pipeline {
      name       = "acctest-pipeline"
      exporters  = [local.workspace_logs_exporter_name, local.workspace_logs_exporter_secondary_name]
      receivers  = [local.syslog_receiver_name, local.syslog_receiver_secondary_name]
      processors = [local.batch_processor_name, local.syslog_processor_name]
    }

    pipeline {
      name       = "acctest-otlp-pipeline"
      exporters  = [local.workspace_logs_exporter_name]
      receivers  = [local.otlp_receiver_name]
      processors = [local.transform_processor_name]
    }

    pipeline {
      name       = "acctest-secondary-pipeline"
      exporters  = [local.workspace_logs_exporter_secondary_name]
      receivers  = [local.syslog_receiver_secondary_name]
      processors = [local.batch_processor_name]
    }
  }

  tags = {
    environment = "staging"
    updated     = "true"
  }

  depends_on = [
    terraform_data.cluster_prereqs,
  ]
}
`, r.template(data), data.RandomInteger)
}
