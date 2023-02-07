## Schema Snapshot Generator

This application generates the schema snapshot for a resource, mainly to be used for [resource state migration](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/state-migration).

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
