load balancer terraform module
===========

A terraform module to provide load balancers in Azure.

Module Input Variables
----------------------
#### Required
- `location` - Name of the Azure datacenter location https://azure.microsoft.com/regions
- `prefix` - name prefix for the lb and public IP address names
- `number_of_endpoints`  - Number of load balancer endpoints to create - default 2
- `frontend_name` - Specifies the name of the frontend ip configuration - default `myPublicIP`
- `public_ip_address_allocation` - Defines how an IP address is assigned. Options are Static or Dynamic. - default static


#### Optional

- `lb_probe_unhealthy_threshold` - Number of times the load balancer health probe has an unsuccessful attempt before considering the endpoint unhealthy. - default 2
- `lb_probe_interval` - Interval in seconds the load balancer health probe rule does a check - default 5

- `remote_port` - desired port(s) for remote connections. Defaults below values.  Set to blank to disable.

```hcl
variable "remote_port" {
  default = { 
      ssh = ["Tcp", "22"]
    }
}
```

- `lb_port` - desired port(s) for load balancer rules. Defaults below values. Set to blank to disable.

```hcl
variable "lb_port" {
  default = {
      http = ["80", "Tcp", "80"]
      https = ["443", "Tcp", "443"]
    }
}
```

```hcl
variable "tags" {
  type = "map"
  default = {
    tag1 = ""
    tag2 = ""
  }
}
```

Usage
-----

```hcl
module "mylb" {
  source   = "../examples/tf-azurerm-loadbalancer"
  prefix   = "mylbpre"
  location = "west us"
  number_of_endpoints = 2
}
```

Outputs
=======

- `azurerm_resource_group_name`
- `number_of_nodes`

...

Authors
=======

* [David Tesar](https://github.com/dtzar)

License
=======

[MIT](LICENSE)