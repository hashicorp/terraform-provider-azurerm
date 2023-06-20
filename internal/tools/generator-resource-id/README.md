## Generator: Resource ID

> **Note:** Resource ID Formatters, Parsers and Structs are now available in [`hashicorp/go-azure-sdk`](https://github.com/hashicorp/go-azure-sdk) and [`commonids` within `hashicorp/go-azure-helpers`](https://github.com/hashicorp/go-azure-helpers/tree/main/resourcemanager/commonids), where available those Resource ID's can be used directly - at this point this package is mostly for legacy resources which haven't yet been migrated to `hashicorp/go-azure-sdk`.

Each Service Definition contains one or more Resource ID's - this tool allows the generation of:

* Resource ID Formatters
* Resource ID Parsers
* Resource ID Structs

This is run via go:generate whenever the provider is compiled - at this time this doesn't wipe an existing "parse" folder so it's possible to mix and match if necessary.

## Example Usage

```
go run main.go -path=./ -name=MyResourceType -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AnalysisServices/servers/Server1
```

## Arguments

* `help` - Show help?

* `id` - An example of the Azure Resource ID for this Resource.

* `name` - The name of this Resource Type, without the Service Name. For example `AnalysisServicesServer` becomes `Server`.

* `path` - The Relative Path to the Service Package.

* `rewrite` - should an `insensitive` parser also be generated to allow for these ID's being rewritten?
