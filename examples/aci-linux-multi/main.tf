resource "azurerm_resource_group" "aci-rg" {
    name="aci-test"
    location="west us"
}

resource "azurerm_container_group" "aci-test" {
    name = "mynginx"
    location = "west us"
    resource_group_name =  "${azurerm_resource_group.aci-rg.name}"
    ip_address_type="public"
    os_type = "linux"

    container {
        name = "hw"
        image = "microsoft/aci-helloworld:latest"
        cpu ="0.5"
        memory =  "1.5"
        port = "80"
    }
    container {
        name = "sidecar"
        image = "microsoft/aci-tutorial-sidecar"
        cpu="0.5"
        memory="1.5"
    }
    tags {
        environment = "testing"
    }
}

