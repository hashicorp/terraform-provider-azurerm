# Create basic network in Azure 

You can use this module to create a basic network in Azure. Below is an example of how you would use this module in your code:

```
module "internal-network" {
  source   = "./path/to/module"
  prefix   = "microservice"
  location = "west us"
}
```