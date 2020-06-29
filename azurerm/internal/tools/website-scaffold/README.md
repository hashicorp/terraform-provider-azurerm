## Website Scaffolder

This application scaffolds the documentation for a Data Source/Resource.

**Note:** the documentation generated from this application is intended to be a starting point, which when finished requires human review - rather than generating a finished product. 

## Example Usage

```
$ go run main.go -name azurerm_resource_group -brand-name "Resource Group" -type "resource" -resource-id "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/group1" -website-path ../../../../website/
```

## Arguments

* `-name` - (Required) The Name used for the Resource in Terraform e.g. `azurerm_resource_group`

* `-brand-name` - (Required) The Brand Name used for this Resource in Azure e.g. `Resource Group` or `App Service (Web Apps)`

* `-type` - (Required) The Type of Documentation to generate. Possible values are `data` (for a Data Source) or `resource` (for a Resource).

* `-resource-id` - (Required when scaffolding a Resource) An Azure Resource ID which can be used as a placeholder in the import documentation.

* `-website-path` - (Required) The path to the `./website` directory in the root of this repository.
