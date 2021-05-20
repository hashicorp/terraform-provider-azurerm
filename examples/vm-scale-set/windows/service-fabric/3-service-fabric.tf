resource "azurerm_service_fabric_cluster" "main" {
  name                = "${var.prefix}-servicefabric"
  resource_group_name = azurerm_resource_group.main.name
  location            = azurerm_resource_group.main.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://${azurerm_public_ip.main.fqdn}:19080"
  add_on_features = ["DnsService"]

  node_type {
    name                 = "primary"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 2020
    http_endpoint_port   = 80
  }

  node_type {
    name                 = local.vmss_name
    instance_count       = local.vmss_count
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080

    application_ports {
      start_port = 20000
      end_port   = 30000
    }

    ephemeral_ports {
      start_port = 49152
      end_port   = 65534
    }
  }

  fabric_settings {
    name = "Security"

    parameters = {
      "ClusterProtectionLevel" = "EncryptAndSign"
    }
  }

  certificate {
    thumbprint           = azurerm_key_vault_certificate.primary.thumbprint
    thumbprint_secondary = azurerm_key_vault_certificate.secondary.thumbprint
    x509_store_name      = "My"
  }
}
