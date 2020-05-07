provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "${var.prefix}-resources"
  location = "${var.location}"
}

resource "azurerm_container_group" "example" {
  name                = "${var.prefix}-continst"
  location            = "${azurerm_resource_group.example.location}"
  resource_group_name = "${azurerm_resource_group.example.name}"
  ip_address_type     = "public"
  os_type             = "linux"

  image_registry_credential {
    server   = "hub.docker.com"
    username = "yourusername1"
    password = "yourpassword"
  }

  image_registry_credential {
    server   = "2hub.docker.com"
    username = "2yourusername1"
    password = "2yourpassword"
  }

  container {
    name   = "hw"
    image  = "microsoft/aci-helloworld:latest"
    cpu    = "0.5"
    memory = "1.5"
    port   = "80"
  }

  container {
    name   = "sidecar"
    image  = "microsoft/aci-tutorial-sidecar"
    cpu    = "0.5"
    memory = "1.5"
  }

  tags = {
    environment = "testing"
  }
}
