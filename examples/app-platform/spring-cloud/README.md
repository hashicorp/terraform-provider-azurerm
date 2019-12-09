# Example: Azure Spring Cloud Service

This example provisions a Spring Cloud Service.

### Notes

* now Azure Spring Cloud only supports location: `East US`, `Southeast Asia`, `West Europe`, `West US 2`.

```hcl
resource "azurerm_spring_cloud" "example" {
  name                     = "example"
  resource_group           = "example"
  location                 = "example"

  tags = {
    env = "staging"
  }
}

resource "azurerm_spring_cloud_config_server" "example" {
  spring_cloud_id = "${azurerm_spring_cloud.example.id}"

  uri = "https://github.com/Azure-Samples/piggymetrics"
  label = "config"
}
```
