## State Migration Generator

This application generates the state migration for a resource, mainly to be used for [resource state migration](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/state-migration). The `UpgradeFunc` body is the only concern to migrate the schema version of a resource with this tool.

## Example Usage

```
$ go run main.go <resource_type>
```

E.g.

```
$ go run main.go azurerm_resource_group
```

## Arguments

* `resource_type`: The resource type to generate the schema. 
