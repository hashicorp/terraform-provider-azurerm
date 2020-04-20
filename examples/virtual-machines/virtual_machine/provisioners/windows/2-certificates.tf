data "azurerm_client_config" "current" {}

resource "azurerm_key_vault" "example" {
  name                = "${var.prefix}keyvault"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  tenant_id           = "${data.azurerm_client_config.current.tenant_id}"

  enabled_for_deployment          = true
  enabled_for_template_deployment = true

  sku_name = "standard"

  access_policy {
    tenant_id = "${data.azurerm_client_config.current.tenant_id}"
    object_id = "${data.azurerm_client_config.current.object_id}"

    certificate_permissions = [
      "create",
      "delete",
      "get",
      "update",
    ]

    key_permissions    = []
    secret_permissions = []
  }
}

resource "azurerm_key_vault_certificate" "example" {
  name         = "${local.virtual_machine_name}-cert"
  key_vault_id = "${azurerm_key_vault.example.id}"

  certificate_policy {
    issuer_parameters {
      name = "Self"
    }

    key_properties {
      exportable = true
      key_size   = 2048
      key_type   = "RSA"
      reuse_key  = true
    }

    lifetime_action {
      action {
        action_type = "AutoRenew"
      }

      trigger {
        days_before_expiry = 30
      }
    }

    secret_properties {
      content_type = "application/x-pkcs12"
    }

    x509_certificate_properties {
      # Server Authentication = 1.3.6.1.5.5.7.3.1
      # Client Authentication = 1.3.6.1.5.5.7.3.2
      extended_key_usage = ["1.3.6.1.5.5.7.3.1"]

      key_usage = [
        "cRLSign",
        "dataEncipherment",
        "digitalSignature",
        "keyAgreement",
        "keyCertSign",
        "keyEncipherment",
      ]

      subject            = "CN=${local.virtual_machine_name}"
      validity_in_months = 12
    }
  }
}
