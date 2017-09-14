load balancer terraform module
===========

A terraform module to provide load balancers in Azure.

Module Input Variables
----------------------
#### Required
- `location` - Name of the Azure datacenter location https://azure.microsoft.com/regions
- `resource_group_name` - The name of the resource group where the load balancer resources will be placed. - default `azure_lb-rg`
- `prefix` - name prefix for the lb and public IP address names - default `azure_lb`
- `number_of_endpoints`  - Number of load balancer endpoints to create - default 2
- `frontend_name` - Specifies the name of the frontend ip configuration - default `myPublicIP`
- `public_ip_address_allocation` - Defines how an IP address is assigned. Options are Static or Dynamic. - default `static`

#### Optional

- `remote_port` - map variable type for the desired port(s) for remote connections to all hosts behind the load balancer. Defaults to an empty map which effectively disables nat rules for remote access to the backend pool.
- `lb_port` - map variable type for desired port(s) for load balancer health probes and rules. Defaults to an empty map which effectively disables health probes and rules.
- `lb_probe_unhealthy_threshold` - Number of times the load balancer health probe has an unsuccessful attempt before considering the endpoint unhealthy. - default 2
- `lb_probe_interval` - Interval in seconds the load balancer health probe rule does a check - default 5
- `tags` - map variable for tags to be placed on the resource group.  Defaults to the below value.

```hcl
variable "tags" {
  type = "map"
  default = {
    source = "terraform"
  }
}
```

Usage
-----

```hcl
module "mylb" {
  source   = "../examples/tf-azurerm-loadbalancer"
  location = "North Central US"
  "remote_port" {
    ssh = ["Tcp", "22"]
  }
  "lb_port" {
    http = ["80", "Tcp", "80"]
    https = ["443", "Tcp", "443"]
  }
  "tags" {
    cost-center = "12345"
    source     = "terraform"
  }
}
```

Outputs
=======
- `azurerm_resource_group_name` - name of the resource group provisioned
- `azurerm_resource_group_tags` - the tags provided for the resource group
- `number_of_nodes` - the number of load balancer nodes provisioned
- `azurerm_lb_id` - the id for the `azurerm_lb` resource
- `azurerm_lb_frontend_ip_configuration` - the `frontend_ip_configuration` for the `azurerm_lb` resource
- `azurerm_lb_probe_ids` - the ids for the `azurerm_lb_probe` resource(s)
- `azurerm_lb_nat_rule_ids` - the ids for the `azurerm_lb_nat_rule` resource(s)
- `azurerm_lb_public_ip_id` - the id for the `azurerm_lb_public_ip` resource
- `azurerm_lb_backend_address_pool_id` - the id for the `azurerm_lb_backend_address_pool` resource

...

Authors
=======

* [David Tesar](https://github.com/dtzar)

License
=======

[MIT](LICENSE)