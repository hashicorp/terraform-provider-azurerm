## Generator: Resource ID

Each Service Definition contains one or more Resource ID's - this tool allows the generation of:

* Resource ID Formatters
* Resource ID Parsers
* Resource ID Structs

This is run via go:generate whenever the provider is compiled - at this time this doesn't wipe an existing "parse" folder so it's possible to mix and match if necessary.

## Example Usage

```
go run main.go -path=-path=./ -name=MyResourceType -id=/subscriptions/12345678-1234-9876-4563-123456789012/resourceGroups/resGroup1/providers/Microsoft.AnalysisServices/servers/Server1
```

## Arguments

* `help` - Show help?

* `id` - An example of the Azure Resource ID for this Resource.

* `name` - The name of this Resource Type, without the Service Name. For example `AnalysisServicesServer` becomes `Server`.

* `path` - The Relative Path to the Service Package.

* `rewrite` - should an `insensitive` parser also be generated to allow for these ID's being rewritten?
