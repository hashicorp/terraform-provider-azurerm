## Typed Model Generator

This application generates the Go structures for a typed resource by providing its schema.

## Example Usage

Given following schema of the `azurerm_resource_group` and ensure the `azurerm_resource_group` is registered as a resource to the provider:

```go
Schema: map[string]*pluginsdk.Schema{
    "name": commonschema.ResourceGroupName(),

    "location": commonschema.Location(),

    "tags": tags.Schema(),
},
```

To generate its corresponding typed SDK model, run the following command:

```go
$ go run main.go azurerm_resource_group
```

This generates:

```go
package main

type ResourceGroupModel struct {
        Location string            `tfschema:"location"`
        Name     string            `tfschema:"name"`
        Tags     map[string]string `tfschema:"tags"`
}
```
