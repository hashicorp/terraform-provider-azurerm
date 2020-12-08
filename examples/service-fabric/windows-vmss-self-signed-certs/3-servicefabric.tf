
resource "azurerm_service_fabric_cluster" "example" {
  name                = "${var.prefix}servicefabric"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  reliability_level   = "Bronze"
  upgrade_mode        = "Automatic"
  vm_image            = "Windows"
  management_endpoint = "https://${var.prefix}servicefabric.${var.location}.cloudapp.azure.com:19080"

  node_type {
    name                 = "Windows"
    instance_count       = 3
    is_primary           = true
    client_endpoint_port = 19000
    http_endpoint_port   = 19080
  }


  reverse_proxy_certificate {
    thumbprint      = azurerm_key_vault_certificate.example.thumbprint
    x509_store_name = "My"
  }

  certificate {
    thumbprint      = azurerm_key_vault_certificate.example.thumbprint
    x509_store_name = "My"
  }

  client_certificate_thumbprint {
    thumbprint = azurerm_key_vault_certificate.example.thumbprint
    is_admin   = true
  }
}

resource "azurerm_windows_virtual_machine_scale_set" "example" {
  name                 = "${var.prefix}examplesf"
  computer_name_prefix = var.prefix
  resource_group_name  = azurerm_resource_group.example.name
  location             = azurerm_resource_group.example.location
  sku                  = "Standard_D1_v2"
  instances            = azurerm_service_fabric_cluster.example.node_type[0].instance_count
  admin_password       = "P@55w0rd1234!"
  admin_username       = "adminuser"
  overprovision        = false
  upgrade_mode         = "Automatic"


  source_image_reference {
    publisher = "MicrosoftWindowsServer"
    offer     = "WindowsServer"
    sku       = "2016-Datacenter-with-Containers"
    version   = "latest"
  }

  os_disk {
    storage_account_type = "Standard_LRS"
    caching              = "ReadWrite"
  }

  network_interface {
    name    = "example"
    primary = true

    ip_configuration {
      name      = "internal"
      primary   = true
      subnet_id = azurerm_subnet.example.id
      load_balancer_backend_address_pool_ids = [
        azurerm_lb_backend_address_pool.example.id
      ]
      load_balancer_inbound_nat_rules_ids = [
        azurerm_lb_nat_pool.example.id
      ]
    }
  }

  secret {
    certificate {
      store = "My"
      url   = azurerm_key_vault_certificate.example.secret_id
    }
    key_vault_id = azurerm_key_vault.example.id
  }

  extension {
    name                       = "${var.prefix}ServiceFabricNode"
    publisher                  = "Microsoft.Azure.ServiceFabric"
    type                       = "ServiceFabricNode"
    type_handler_version       = "1.1"
    auto_upgrade_minor_version = false

    settings = jsonencode({
      "clusterEndpoint"    = azurerm_service_fabric_cluster.example.cluster_endpoint
      "nodeTypeRef"        = azurerm_service_fabric_cluster.example.node_type[0].name
      "durabilityLevel"    = "bronze"
      "nicPrefixOverride"  = azurerm_subnet.example.address_prefixes[0]
      "enableParallelJobs" = true
      "certificate" = {
        "commonNames" = [
          "${var.prefix}servicefabric.${var.location}.cloudapp.azure.com",
        ]
        "x509StoreName" = "My"
      }
    })

    protected_settings = jsonencode({
      "StorageAccountKey1" = azurerm_storage_account.example.primary_access_key
      "StorageAccountKey2" = azurerm_storage_account.example.secondary_access_key
    })
  }

}
